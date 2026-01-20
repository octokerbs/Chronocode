package repo

import "context"

type RepoRepository interface {
	Get(ctx context.Context, id int64) (*Repository, error)
	Store(ctx context.Context, repo *Repository) error
}
