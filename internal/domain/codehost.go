package domain

import (
	"context"
)

type CodeHost interface {
	FetchRepository(ctx context.Context, repoURL string) (*Repository, error)
	FetchRepositoryID(ctx context.Context, repoURL string) (int64, error)
	FetchCommit(ctx context.Context, repoURL string, commitSHA string) (*Commit, error)
	FetchCommitDiff(ctx context.Context, repoURL string, commitSHA string) (string, error)

	ProduceCommitSHAs(ctx context.Context, repoURL string, lastAnalyzedCommitSHA string, commits chan<- string)
}
