package subcommit

import "context"

type Repository interface {
	GetSubcommits(ctx context.Context, repoID int64) ([]Subcommit, error)
	HasSubcommitsForCommit(ctx context.Context, repoID int64, commitSHA string) (bool, error)
	StoreSubcommits(ctx context.Context, subcommits <-chan Subcommit) error
}
