package domain

import (
	"context"
)

type Database interface {
	GetRepository(ctx context.Context, id int64) (*Repository, bool, error)

	StoreRepository(ctx context.Context, repo *Repository) error
	StoreCommits(ctx context.Context, commits []*Commit) error
	StoreSubcommits(ctx context.Context, subcommits []*Subcommit) error
}
