package mocks

import (
	"context"

	"github.com/octokerbs/chronocode/internal/domain/codehost"
	"github.com/octokerbs/chronocode/internal/domain/repo"
)

var (
	ValidRepoURL   = "https/validRepo"
	InvalidRepoURL = "https/invalidRepo"
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
