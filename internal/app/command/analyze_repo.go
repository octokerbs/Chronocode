package command

import (
	"context"
	"errors"
	"sync"

	"github.com/octokerbs/chronocode/internal/domain/agent"
	"github.com/octokerbs/chronocode/internal/domain/codehost"
	"github.com/octokerbs/chronocode/internal/domain/repo"
	"github.com/octokerbs/chronocode/internal/domain/subcommit"
)

type AnalyzeRepo struct {
	RepoURL string
}
type AnalyzeRepoHandler struct {
	repoRepository      repo.Repository
	subcommitRepository subcommit.Repository
	agent               agent.Agent
	codeHost            codehost.CodeHost
}

func NewAnalyzeRepoHandler(repoRepository repo.Repository, subcommitRepository subcommit.Repository, agent agent.Agent, codeHost codehost.CodeHost) AnalyzeRepoHandler {
	return AnalyzeRepoHandler{repoRepository: repoRepository, subcommitRepository: subcommitRepository, agent: agent, codeHost: codeHost}
}

func (s *AnalyzeRepoHandler) Handle(ctx context.Context, cmd AnalyzeRepo) error {
	newRepo, err := s.repoRepository.GetRepo(ctx, cmd.RepoURL)
	if err != nil {
		if !errors.Is(err, repo.ErrRepositoryNotFound) {
			return err
		}

		newRepo, err = s.codeHost.CreateRepoFromURL(ctx, cmd.RepoURL)
		if err != nil {
			return err
		}
	}

	var wg sync.WaitGroup
	commitSHAs := make(chan string, 100)
	subcommits := make(chan subcommit.Subcommit, 100)

	wg.Add(3)

	go func() {
		defer wg.Done()
		s.codeHost.GetRepoCommitSHAsIntoChannel(ctx, newRepo, commitSHAs)
		close(commitSHAs)
	}()

	go func() {
		defer wg.Done()
		s.agent.AnalyzeCommitsIntoSubcommits(ctx, commitSHAs, subcommits)
		close(subcommits)
	}()

	go func() {
		defer wg.Done()
		s.subcommitRepository.StoreSubcommits(ctx, subcommits)
	}()

	wg.Wait()

	err = s.repoRepository.StoreRepo(ctx, newRepo)
	return err
}
