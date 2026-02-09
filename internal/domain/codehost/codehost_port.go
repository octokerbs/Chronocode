package codehost

import (
	"context"

	"github.com/octokerbs/chronocode-backend/internal/domain/analysis"
	"github.com/octokerbs/chronocode-backend/internal/domain/user"
)

type CodeHost interface {
	FetchRepository(ctx context.Context, repoURL string) (*analysis.Repository, error)
	FetchRepositoryID(ctx context.Context, repoURL string) (int64, error)
	FetchCommit(ctx context.Context, repoURL string, commitSHA string) (*analysis.Commit, error)
	FetchCommitDiff(ctx context.Context, repoURL string, commitSHA string) (string, error)

	ProduceCommitSHAs(ctx context.Context, repoURL string, lastAnalyzedCommitSHA string, commitSHAs chan<- string) (string, error)

	FetchAuthenticatedUser(ctx context.Context) (*user.GitHubProfile, error)
	SearchUserRepositories(ctx context.Context, query string) ([]*analysis.Repository, error)
}
