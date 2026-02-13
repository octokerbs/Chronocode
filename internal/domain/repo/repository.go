package repo

import (
	"context"
	"errors"
)

var (
	ErrRepositoryNotFound = errors.New("repository not found")
)

type Repository interface {
	GetRepo(ctx context.Context, url string) (*Repo, error)
	StoreRepo(ctx context.Context, aRepo *Repo) error
}
