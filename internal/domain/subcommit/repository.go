package subcommit

import "context"

type Repository interface {
	GetByRepoID(ctx context.Context, repoID int64) ([]*Subcommit, error)
}
