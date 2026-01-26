package repository

import "context"

type Repository interface {
	Get(ctx context.Context, id int64) (*Repo, error)
	Store(ctx context.Context, repo *Repo) error
}
