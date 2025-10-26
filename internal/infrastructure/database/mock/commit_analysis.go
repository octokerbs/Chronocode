package mock

import (
	"context"

	"github.com/octokerbs/chronocode-backend/internal/application"
	"github.com/octokerbs/chronocode-backend/internal/domain"
	"github.com/stretchr/testify/mock"
)

type PostgresMock struct {
	mock.Mock
}

func (pm *PostgresMock) InsertRepositoryRecord(ctx context.Context, repo *domain.Repository) error {
	return nil
}

func (pm *PostgresMock) InsertCommitRecord(ctx context.Context, commit *domain.Commit) error {
	return nil
}

func (pm *PostgresMock) InsertSubcommitRecord(ctx context.Context, subcommit *domain.Subcommit) error {

	return nil
}

func (pm *PostgresMock) GetRepository(ctx context.Context, id int64) (*domain.Repository, bool, error) {
	return nil, false, nil
}

func (pm *PostgresMock) ProcessRecords(ctx context.Context, records <-chan application.DatabaseRecord, errors chan<- string) {

}
