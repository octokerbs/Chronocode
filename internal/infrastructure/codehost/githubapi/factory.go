package githubapi

import (
	"context"

	"github.com/octokerbs/chronocode-backend/internal/domain/codehost"
)

type GitHubCodeHostFactory struct {
}

func NewGitHubCodeHostFactory() *GitHubCodeHostFactory {
	return &GitHubCodeHostFactory{}
}

func (f *GitHubCodeHostFactory) Create(ctx context.Context, accessToken string) codehost.CodeHost {
	return NewGithubCodeHost(ctx, accessToken)
}
