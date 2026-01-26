package command

import (
	"context"
	"sync"

	"github.com/octokerbs/chronocode-backend/common/decorator"
	"github.com/octokerbs/chronocode-backend/internal/domain/repository"
	"go.uber.org/zap"
)

type AnalyzeRepo struct {
	repoURL     string
	accessToken string
}

type AnalyzeRepoHandler decorator.CommandHandler[AnalyzeRepo]

type analyzeRepoHandler struct {
	agent           Agent
	codeHostFactory CodeHostFactory
	repoRepository  repository.Repository
}

func NewAnalyzeRepoHandler(repoRepository repository.Repository, agent Agent, codeHostFactory CodeHostFactory, logger *zap.Logger) AnalyzeRepoHandler {
	return decorator.ApplyCommandDecorators[AnalyzeRepo](
		analyzeRepoHandler{},
		logger,
	)
}

func (h analyzeRepoHandler) Handle(ctx context.Context, cmd AnalyzeRepo) error {
	repo, err := h.prepareRepository(ctx, cmd.repoURL, cmd.accessToken)
	if err != nil {
		return err
	}

	codeHost := h.codeHostFactory.Create(ctx, cmd.accessToken)

	rawCommitSHAs := make(chan string, 100)
	analyzedCommits := make(chan repository.Commit, 100)
	var wg sync.WaitGroup

	// Workers only analyze and emit events
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go h.analyzerWorker(ctx, cmd.repoURL, codeHost, rawCommitSHAs, analyzedCommits, &wg)
	}

	go func() {
		defer close(rawCommitSHAs)
		codeHost.ProduceCommitSHAs(ctx, cmd.repoURL, repo.LastAnalyzedCommit, rawCommitSHAs)
	}()

	wg.Wait()
	return nil
}

type Agent interface {
	AnalyzeCommitDiff(ctx context.Context, diff string) (repository.CommitAnalysis, error)
}

type CodeHostFactory interface {
	Create(ctx context.Context, accessToken string) CodeHost
}

type CodeHost interface {
	FetchRepository(ctx context.Context, repoURL string) (*repository.Repo, error)
	FetchRepositoryID(ctx context.Context, repoURL string) (int64, error)
	FetchCommit(ctx context.Context, repoURL string, commitSHA string) (*repository.Commit, error)
	FetchCommitDiff(ctx context.Context, repoURL string, commitSHA string) (string, error)

	ProduceCommitSHAs(ctx context.Context, repoURL string, lastAnalyzedCommitSHA string, commitSHAs chan<- string) (string, error)
}

func (h analyzeRepoHandler) prepareRepository(ctx context.Context, repoURL, accessToken string) (*repository.Repo, error) {
	codeHost := h.codeHostFactory.Create(ctx, accessToken)

	fetchedRepo, err := codeHost.FetchRepository(ctx, repoURL)
	if err != nil {
		return nil, err
	}

	repo, err := h.repoRepository.Get(ctx, fetchedRepo.ID)
	if err != nil {
		err := h.repoRepository.Store(ctx, fetchedRepo)
		if err != nil {
			return nil, err
		}
		return fetchedRepo, nil
	}

	return repo, err
}

func (h analyzeRepoHandler) analyzerWorker(ctx context.Context, repoURL string, codeHost CodeHost, rawCommitSHAs <-chan string, analyzedCommits chan<- repository.Commit, wg *sync.WaitGroup) {
	defer wg.Done()

	for commitSHA := range rawCommitSHAs {
		diff, err := codeHost.FetchCommitDiff(ctx, repoURL, commitSHA)
		if err != nil {
			continue
		}

		analysis, err := h.agent.AnalyzeCommitDiff(ctx, diff)
		if err != nil {
			continue
		}

		commit, err := codeHost.FetchCommit(ctx, repoURL, commitSHA)
		if err != nil {
			continue
		}

		commit.ApplyAnalysis(&analysis)

		select {
		case analyzedCommits <- *commit:
		case <-ctx.Done():
			return
		}
	}
}
