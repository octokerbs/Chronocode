package application

import (
	"context"
	"sync"
	"time"

	"github.com/octokerbs/chronocode-backend/internal/domain/analysis"
	"github.com/octokerbs/chronocode-backend/internal/domain/codehost"
)

type CommitAnalyzed struct {
	Commit    *analysis.Commit
	Timestamp time.Time
}

type Analyzer struct {
	agent           analysis.Agent
	codeHostFactory codehost.CodeHostFactory
}

func NewAnalyzer(agent analysis.Agent, factory codehost.CodeHostFactory) *Analyzer {
	return &Analyzer{
		agent:           agent,
		codeHostFactory: factory,
	}
}

func (a *Analyzer) AnalyzeCommits(ctx context.Context, repo *analysis.Repository, events chan<- CommitAnalyzed, accessToken string) error {
	codeHost := a.codeHostFactory.Create(ctx, accessToken)

	commitSHAs := make(chan string)
	var wg sync.WaitGroup

	// Workers only analyze and emit events
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go a.analyzerWorker(ctx, repo.URL, codeHost, events, commitSHAs, &wg)
	}

	go func() {
		defer close(commitSHAs)
		codeHost.ProduceCommitSHAs(ctx, repo.URL, repo.LastAnalyzedCommit, commitSHAs)
	}()

	wg.Wait()
	return nil
}

func (a *Analyzer) analyzerWorker(ctx context.Context, repoURL string, codeHost codehost.CodeHost, events chan<- CommitAnalyzed, commitSHAs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for commitSHA := range commitSHAs {
		diff, err := codeHost.FetchCommitDiff(ctx, repoURL, commitSHA)
		if err != nil {
			continue
		}

		analysis, err := a.agent.AnalyzeCommitDiff(ctx, diff)
		if err != nil {
			continue
		}

		commit, err := codeHost.FetchCommit(ctx, repoURL, commitSHA)
		if err != nil {
			continue
		}

		commit.ApplyAnalysis(&analysis)

		select {
		case events <- CommitAnalyzed{Commit: commit, Timestamp: time.Now()}:
		case <-ctx.Done():
			return
		}
	}
}
