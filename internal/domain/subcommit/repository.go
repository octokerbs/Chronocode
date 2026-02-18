package subcommit

import "context"

type Repository interface {
	GetSubcommits(ctx context.Context, repoID int64) ([]Subcommit, error)
	StoreSubcommits(ctx context.Context, subcommits <-chan Subcommit) error
}
