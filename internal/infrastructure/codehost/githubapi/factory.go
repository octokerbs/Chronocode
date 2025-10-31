package githubapi

import (
	"context"

	"github.com/octokerbs/chronocode-backend/internal/domain"
)

type GitHubFactory struct {
}

func NewGitHubFactory() *GitHubFactory {
	return &GitHubFactory{}
}

func (f *GitHubFactory) Create(ctx context.Context, accessToken string) domain.CodeHost {
	return NewGithubCodeHost(ctx, accessToken)
}
