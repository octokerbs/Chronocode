package domain

import (
	"context"
)

type Database interface {
	InsertRepository(ctx context.Context, repo *RepositoryRecord) error
	InsertCommit(ctx context.Context, commit *CommitRecord) error
	InsertSubcommit(ctx context.Context, subcommit *SubcommitRecord) error

	GetRepository(ctx context.Context, id int64) (*RepositoryRecord, bool, error)
}

type Record interface {
	InsertIntoDatabase(ctx context.Context, database Database) error
}
