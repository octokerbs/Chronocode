package codehost

import (
	"context"
	"errors"

	"github.com/octokerbs/chronocode/internal/domain/repo"
)

var (
	ErrInvalidRepoURL  = errors.New("invalid repository URL")
	ErrAccessDenied    = errors.New("access denied to repository")
	ErrDiffFetchFailed = errors.New("failed to fetch commit diff")
)

type CodeHostFactory interface {
	Create(ctx context.Context, accessToken string) (CodeHost, error)
}

type CodeHost interface {
	CanAccessRepo(ctx context.Context, repoURL string) error
	CreateRepoFromURL(ctx context.Context, url string) (*repo.Repo, error)
	GetRepoCommitSHAsIntoChannel(ctx context.Context, repo *repo.Repo, commitSHAs chan<- string) error
	GetCommitDiff(ctx context.Context, repo *repo.Repo, commitSHA string) (string, error)
}
