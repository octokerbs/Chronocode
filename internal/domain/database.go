package domain

import (
	"context"
	"errors"
)

var ErrRepositoryNotFound = errors.New("repository not found")

type Database interface {
	GetRepository(ctx context.Context, id int64) (*Repository, error)

	StoreRepository(ctx context.Context, repo *Repository) error
	StoreCommits(ctx context.Context, commits []*Commit) error
}
