package command

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/octokerbs/chronocode/internal/domain/agent"
	"github.com/octokerbs/chronocode/internal/domain/codehost"
	"github.com/octokerbs/chronocode/internal/domain/repo"
	"github.com/octokerbs/chronocode/internal/domain/subcommit"
)

const maxConcurrentAnalyses = 10

type AnalyzeRepo struct {
	RepoURL     string
	AccessToken string
}

type AnalyzeRepoHandler struct {
	repoRepository      repo.Repository
	subcommitRepository subcommit.Repository
	agent               agent.Agent
	codeHostFactory     codehost.CodeHostFactory
}

func NewAnalyzeRepoHandler(repoRepository repo.Repository, subcommitRepository subcommit.Repository, agent agent.Agent, codeHostFactory codehost.CodeHostFactory) AnalyzeRepoHandler {
	return AnalyzeRepoHandler{repoRepository: repoRepository, subcommitRepository: subcommitRepository, agent: agent, codeHostFactory: codeHostFactory}
}

func (s *AnalyzeRepoHandler) Handle(ctx context.Context, cmd AnalyzeRepo) error {
	codeHost, err := s.codeHostFactory.Create(ctx, cmd.AccessToken)
	if err != nil {
		return err
	}

	if err := codeHost.CanAccessRepo(ctx, cmd.RepoURL); err != nil {
		return err
	}

	newRepo, err := s.repoRepository.GetRepo(ctx, cmd.RepoURL)
	if err != nil {
		if !errors.Is(err, repo.ErrRepositoryNotFound) {
			return err
		}

		newRepo, err = codeHost.CreateRepoFromURL(ctx, cmd.RepoURL)
		if err != nil {
			return err
		}
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	commitSHAs := make(chan string, 100)
	subcommits := make(chan subcommit.Subcommit, 100)

	var fetchErr, analysisErr, storageErr error
	wg.Add(3)

	go func() {
		defer wg.Done()
		defer close(commitSHAs)
		fetchErr = codeHost.GetRepoCommitSHAsIntoChannel(ctx, newRepo, commitSHAs)
		if fetchErr != nil {
			cancel()
		}
	}()

	go func() {
		defer wg.Done()
		defer close(subcommits)
		analysisErr = s.analyzeCommits(ctx, codeHost, newRepo, commitSHAs, subcommits)
	}()

	go func() {
		defer wg.Done()
		storageErr = s.subcommitRepository.StoreSubcommits(ctx, subcommits)
	}()

	wg.Wait()

	if fetchErr != nil {
		return fetchErr
	}
	if analysisErr != nil {
		return analysisErr
	}
	if storageErr != nil {
		return storageErr
	}

	return s.repoRepository.StoreRepo(ctx, newRepo)
}

func (s *AnalyzeRepoHandler) analyzeCommits(ctx context.Context, codeHost codehost.CodeHost, r *repo.Repo, commitSHAs <-chan string, subcommits chan<- subcommit.Subcommit) error {
	var wg sync.WaitGroup
	var firstErr error
	var errOnce sync.Once
	sem := make(chan struct{}, maxConcurrentAnalyses)

	for sha := range commitSHAs {
		if ctx.Err() != nil {
			break
		}

		sem <- struct{}{}

		wg.Add(1)
		go func(sha string) {
			defer func() { <-sem }()
			defer wg.Done()

			if ctx.Err() != nil {
				return
			}

			diff, err := codeHost.GetCommitDiff(ctx, r, sha)
			if err != nil {
				errOnce.Do(func() { firstErr = fmt.Errorf("%w: %v", codehost.ErrDiffFetchFailed, err) })
				return
			}

			results, err := s.agent.AnalyzeDiff(ctx, diff)
			if err != nil {
				errOnce.Do(func() { firstErr = fmt.Errorf("%w: %v", agent.ErrAnalysisFailed, err) })
				return
			}

			for _, result := range results {
				subcommits <- subcommit.NewSubcommit(result.Title, result.Description, result.ModificationType, sha, result.Files, r.ID())
			}
		}(sha)
	}

	wg.Wait()
	return firstErr
}
