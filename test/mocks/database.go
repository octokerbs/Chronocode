package mocks

import (
	"context"

	"github.com/octokerbs/chronocode-go/internal/repository"
	"github.com/stretchr/testify/mock"
)

type DatabaseServiceMock struct {
	mock.Mock
}

func (d *DatabaseServiceMock) InsertRepository(ctx context.Context, repo *repository.RepositoryRecord) error {
	return nil
}
func (d *DatabaseServiceMock) InsertCommit(ctx context.Context, commit *repository.CommitRecord) error {

	return nil
}
func (d *DatabaseServiceMock) InsertSubcommit(ctx context.Context, subcommit *repository.SubcommitRecord) error {

	return nil
}

func (d *DatabaseServiceMock) GetRepository(ctx context.Context, id int64) (*repository.RepositoryRecord, bool, error) {
	args := d.Called()
	if args.Get(2) != nil {
		return args.Get(0).(*repository.RepositoryRecord), args.Get(1).(bool), args.Get(2).(error)
	}
	return args.Get(0).(*repository.RepositoryRecord), args.Get(1).(bool), nil
}
