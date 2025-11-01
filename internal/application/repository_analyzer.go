package application

import (
	"context"
	"sync"

	"github.com/octokerbs/chronocode-backend/internal/domain"
)

type RepositoryAnalyzer struct {
	Agent           domain.Agent
	CodeHostFactory domain.CodeHostFactory
	Database        domain.Database

	resultsMutex    sync.Mutex
	analyzedCommits []*domain.Commit
}

func NewRepositoryAnalyzer(ctx context.Context, agent domain.Agent, codehostFactory domain.CodeHostFactory, database domain.Database) *RepositoryAnalyzer {
	return &RepositoryAnalyzer{
		Agent:           agent,
		CodeHostFactory: codehostFactory,
		Database:        database,
		analyzedCommits: []*domain.Commit{},
	}
}

func (ra *RepositoryAnalyzer) AnalyzeRepository(ctx context.Context, repoURL string, accessToken string) error {
	codeHost := ra.CodeHostFactory.Create(ctx, accessToken)

	commits := make(chan string)

	repo, err := ra.fetchOrCreateRepository(ctx, repoURL, codeHost)
	if err != nil {
		return err
	}

	// Clear slices from any previous run
	ra.analyzedCommits = nil

	var wg sync.WaitGroup
	for range 200 {
		wg.Add(1)
		go ra.commitAnalyzerWorker(ctx, repoURL, codeHost, commits, &wg)
	}

	lastAnalyzedCommitSHA := repo.LastAnalyzedCommit

	go func() {
		// When ProduceCommits is done, it will close the channel
		codeHost.ProduceCommitSHAs(ctx, repoURL, lastAnalyzedCommitSHA, commits)
	}()

	wg.Wait()

	if len(ra.analyzedCommits) > 0 {
		if err := ra.Database.StoreCommits(ctx, ra.analyzedCommits); err != nil {
			return err
		}
	}

	return nil
}

func (ra *RepositoryAnalyzer) fetchOrCreateRepository(ctx context.Context, repoURL string, codeHost domain.CodeHost) (*domain.Repository, error) {
	id, err := codeHost.FetchRepositoryID(ctx, repoURL)
	if err != nil {
		return nil, err
	}

	repo, ok, err := ra.Database.GetRepository(ctx, id)
	if err != nil {
		return nil, err
	}

	if ok {
		return repo, nil
	}

	repo, err = codeHost.FetchRepository(ctx, repoURL)
	if err != nil {
		return nil, err
	}

	ra.Database.StoreRepository(ctx, repo)

	return repo, nil
}

func (ra *RepositoryAnalyzer) commitAnalyzerWorker(ctx context.Context, repoURL string, codeHost domain.CodeHost, commits <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for commitSHA := range commits {
		diff, err := codeHost.FetchCommitDiff(ctx, repoURL, commitSHA)
		if err != nil {
			continue
		}

		analysis, err := ra.Agent.AnalyzeCommitDiff(ctx, diff)
		if err != nil {
			continue
		}

		commit, err := codeHost.FetchCommit(ctx, repoURL, commitSHA)
		if err != nil {
			continue
		}

		commit.ApplyAnalysis(&analysis)

		ra.resultsMutex.Lock()
		ra.analyzedCommits = append(ra.analyzedCommits, commit)
		ra.resultsMutex.Unlock()
	}
}
