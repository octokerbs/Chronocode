package domain

import (
	"context"
)

type Database interface {
	GetRepository(ctx context.Context, id int64) (*Repository, error)

	StoreRepository(ctx context.Context, repo *Repository) error
	StoreCommit(ctx context.Context, commit *Commit) error
}
