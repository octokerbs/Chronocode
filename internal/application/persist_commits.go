package application

import (
	"context"
	"fmt"
	"sync"

	"github.com/octokerbs/chronocode-backend/internal/domain/database"
)

type PersistCommits struct {
	Database database.Database
}

func (pc *PersistCommits) HandleCommitAnalyzed(ctx context.Context, events <-chan CommitAnalyzed) error {
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go pc.persistWorker(ctx, events, &wg)
	}

	wg.Wait()
	return nil
}

func (pc *PersistCommits) persistWorker(ctx context.Context, events <-chan CommitAnalyzed, wg *sync.WaitGroup) {
	defer wg.Done()

	for event := range events {
		if err := pc.Database.StoreCommit(ctx, event.Commit); err != nil {
			fmt.Printf("Failed to store commit: %e", err)
		}
	}
}
