package mocks

import (
	"context"

	"github.com/chrono-code-hackathon/chronocode-go/internal/database"
	"github.com/stretchr/testify/mock"
)

type DatabaseServiceMock struct {
	mock.Mock
}

func (d *DatabaseServiceMock) InsertRepository(ctx context.Context, repo *database.RepositoryRecord) error {
	return nil
}
func (d *DatabaseServiceMock) InsertCommit(ctx context.Context, commit *database.CommitRecord) error {

	return nil
}
func (d *DatabaseServiceMock) InsertSubcommit(ctx context.Context, subcommit *database.SubcommitRecord) error {

	return nil
}

func (d *DatabaseServiceMock) GetRepository(ctx context.Context, id int64) (*database.RepositoryRecord, bool, error) {
	args := d.Called()
	if args.Get(2) != nil {
		return args.Get(0).(*database.RepositoryRecord), args.Get(1).(bool), args.Get(2).(error)
	}
	return args.Get(0).(*database.RepositoryRecord), args.Get(1).(bool), nil
}
