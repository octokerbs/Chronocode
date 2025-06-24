package mocks

import (
	"context"
	"time"

	"github.com/chrono-code-hackathon/chronocode-go/internal/database"
)

type DatabaseMock struct {
	InsertionDuration time.Duration
}

func (s *DatabaseMock) InsertRepository(ctx context.Context, repo *database.RepositoryRecord) error {
	time.Sleep(s.InsertionDuration)
	return nil
}

func (s *DatabaseMock) InsertCommit(ctx context.Context, commit *database.CommitRecord) error {
	time.Sleep(s.InsertionDuration)
	return nil
}

func (s *DatabaseMock) InsertSubcommit(ctx context.Context, subcommit *database.SubcommitRecord) error {
	time.Sleep(s.InsertionDuration)
	return nil
}

func (s *DatabaseMock) GetRepository(ctx context.Context, id int64) (*database.RepositoryRecord, bool, error) {
	time.Sleep(s.InsertionDuration)
	today := time.Now()

	return &database.RepositoryRecord{
		ID:                 514170585,
		CreatedAt:          &today,
		Name:               "TusLibros",
		URL:                "https://github.com/octokerbs/TusLibros",
		LastAnalyzedCommit: "f41d775b09231a4d377a8133cae7d62dd1fa2fae",
	}, true, nil
}
