package domain

import (
	"context"
)

type Database interface {
	GetRepository(ctx context.Context, id int64) (*Repository, error)
	GetSubcommitsByRepoID(ctx context.Context, repoID int64) ([]*Subcommit, error)

	StoreRepository(ctx context.Context, repo *Repository) error
	StoreCommit(ctx context.Context, commit *Commit) error
}
