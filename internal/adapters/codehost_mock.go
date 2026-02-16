package adapters

import (
	"context"

	"github.com/octokerbs/chronocode/internal/domain/codehost"
	"github.com/octokerbs/chronocode/internal/domain/repo"
)

var (
	ValidRepoURL             = "https/validRepo"
	ValidRepoID        int64 = 123456789
	ValidRepoCommitSHA       = "CommitSHA-1"
	ValidEmptyRepoURL        = "https/emptyRepo"
	ValidEmptyRepoID   int64 = 9876543221
	InvalidRepoURL           = "https/invalidRepo"
)

type CodeHost struct {
}

func NewCodeHost() *CodeHost {
	return &CodeHost{}
}

func (c *CodeHost) CreateRepoFromURL(ctx context.Context, url string) (*repo.Repo, error) {
	if url == InvalidRepoURL {
		return nil, codehost.ErrInvalidRepoURL
	}

	newRepo := repo.NewRepo(1, "chronocode", ValidRepoURL, "FFFFFF")
	return newRepo, nil
}

func (c *CodeHost) GetRepoCommitSHAsIntoChannel(ctx context.Context, repo *repo.Repo, commitSHAs chan<- string) {
	if repo.URL() == ValidEmptyRepoURL {
		return
	}

	commitSHAs <- ValidRepoCommitSHA
}
