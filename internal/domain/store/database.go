package store

import (
	"context"

	"github.com/octokerbs/chronocode-backend/internal/domain/analysis"
)

type Database interface {
	GetRepository(ctx context.Context, id int64) (*analysis.Repository, error)
	GetSubcommitsByRepoID(ctx context.Context, repoID int64) ([]*analysis.Subcommit, error)

	StoreRepository(ctx context.Context, repo *analysis.Repository) error
	StoreCommit(ctx context.Context, commit *analysis.Commit) error
}
