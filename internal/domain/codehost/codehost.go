package codehost

import (
	"context"
	"errors"

	"github.com/octokerbs/chronocode/internal/domain/repo"
)

var (
	ErrInvalidRepoURL = errors.New("invalid repository URL")
)

type CodeHost interface {
	CreateRepoFromURL(ctx context.Context, url string) (*repo.Repo, error)
}
