package command

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/octokerbs/chronocode/internal/domain/agent"
	"github.com/octokerbs/chronocode/internal/domain/analysis"
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
	locker              analysis.Locker
}

func NewAnalyzeRepoHandler(repoRepository repo.Repository, subcommitRepository subcommit.Repository, agent agent.Agent, codeHostFactory codehost.CodeHostFactory, locker analysis.Locker) AnalyzeRepoHandler {
	return AnalyzeRepoHandler{repoRepository: repoRepository, subcommitRepository: subcommitRepository, agent: agent, codeHostFactory: codeHostFactory, locker: locker}
}

func (s *AnalyzeRepoHandler) Handle(ctx context.Context, cmd AnalyzeRepo) (int64, error) {
	codeHost, err := s.codeHostFactory.Create(ctx, cmd.AccessToken)
	if err != nil {
		return 0, err
	}

	if err := codeHost.CanAccessRepo(ctx, cmd.RepoURL); err != nil {
		return 0, err
	}

	release, err := s.locker.Acquire(ctx, cmd.RepoURL)
	if err != nil {
		return 0, err
	}
	defer release()

	newRepo, err := s.repoRepository.GetRepo(ctx, cmd.RepoURL)
	if err != nil {
		if !errors.Is(err, repo.ErrRepositoryNotFound) {
			return 0, err
		}

		newRepo, err = codeHost.CreateRepoFromURL(ctx, cmd.RepoURL)
		if err != nil {
			return 0, err
		}
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	commitRefs := make(chan codehost.CommitReference, 100)
	subcommits := make(chan subcommit.Subcommit, 100)

	var fetchErr, analysisErr, storageErr error
	var headSHA string
	wg.Add(3)

	go func() {
		defer wg.Done()
		defer close(commitRefs)
		headSHA, fetchErr = codeHost.GetRepoCommitSHAsIntoChannel(ctx, newRepo, commitRefs)
		if fetchErr != nil {
			cancel()
		}
	}()

	go func() {
		defer wg.Done()
		defer close(subcommits)
		analysisErr = s.analyzeCommits(ctx, codeHost, newRepo, commitRefs, subcommits)
	}()

	go func() {
		defer wg.Done()
		storageErr = s.subcommitRepository.StoreSubcommits(ctx, subcommits)
	}()

	wg.Wait()

	if fetchErr != nil {
		return 0, fetchErr
	}

	if analysisErr == nil && storageErr == nil && headSHA != "" {
		newRepo.SetLastAnalyzedCommitSHA(headSHA)
	}

	if err := s.repoRepository.StoreRepo(ctx, newRepo); err != nil {
		return 0, err
	}

	return newRepo.ID(), errors.Join(analysisErr, storageErr)
}

func (s *AnalyzeRepoHandler) analyzeCommits(ctx context.Context, codeHost codehost.CodeHost, r *repo.Repo, commitRefs <-chan codehost.CommitReference, subcommits chan<- subcommit.Subcommit) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error
	sem := make(chan struct{}, maxConcurrentAnalyses)

	for ref := range commitRefs {
		if ctx.Err() != nil {
			break
		}

		sem <- struct{}{}

		wg.Add(1)
		go func(ref codehost.CommitReference) {
			defer func() { <-sem }()
			defer wg.Done()

			if ctx.Err() != nil {
				return
			}

			alreadyAnalyzed, err := s.subcommitRepository.HasSubcommitsForCommit(ctx, r.ID(), ref.SHA)
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
				return
			}
			if alreadyAnalyzed {
				return
			}

			diff, err := codeHost.GetCommitDiff(ctx, r, ref.SHA)
			if err != nil {
				mu.Lock()
				errs = append(errs, fmt.Errorf("%w: %v", codehost.ErrDiffFetchFailed, err))
				mu.Unlock()
				return
			}

			results, err := s.agent.AnalyzeDiff(ctx, diff)
			if err != nil {
				mu.Lock()
				errs = append(errs, fmt.Errorf("%w: %v", agent.ErrAnalysisFailed, err))
				mu.Unlock()
				return
			}

			for _, result := range results {
				subcommits <- subcommit.NewSubcommit(result.Title, result.Idea, result.Description, result.Epic, result.ModificationType, ref.SHA, result.Files, r.ID(), ref.CommittedAt)
			}
		}(ref)
	}

	wg.Wait()
	return errors.Join(errs...)
}
