package codehost

import (
	"context"
	"errors"
	"time"

	"github.com/octokerbs/chronocode/internal/domain/repo"
)

var (
	ErrInvalidRepoURL  = errors.New("invalid repository URL")
	ErrAccessDenied    = errors.New("access denied to repository")
	ErrDiffFetchFailed = errors.New("failed to fetch commit diff")
)

type CommitReference struct {
	SHA         string
	CommittedAt time.Time
}

type CodeHostFactory interface {
	Create(ctx context.Context, accessToken string) (CodeHost, error)
}

type CodeHost interface {
	CanAccessRepo(ctx context.Context, repoURL string) error
	CreateRepoFromURL(ctx context.Context, url string) (*repo.Repo, error)
	// GetRepoCommitSHAsIntoChannel sends non-merge commits newest-first, stopping
	// at repo.LastAnalyzedCommitSHA() (exclusive). Returns the head SHA (first commit
	// sent) or "" if no commits were sent.
	GetRepoCommitSHAsIntoChannel(ctx context.Context, repo *repo.Repo, commits chan<- CommitReference) (headSHA string, err error)
	GetCommitDiff(ctx context.Context, repo *repo.Repo, commitSHA string) (string, error)
}
