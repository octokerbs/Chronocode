package githubapi

import (
	"context"

	"github.com/octokerbs/chronocode-backend/internal/domain/codehost"
)

type GitHubFactory struct {
}

func NewGitHubFactory() *GitHubFactory {
	return &GitHubFactory{}
}

func (f *GitHubFactory) Create(ctx context.Context, accessToken string) codehost.CodeHost {
	return NewGithubCodeHost(ctx, accessToken)
}
