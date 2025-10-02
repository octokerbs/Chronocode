package domain

import (
	"context"
)

type CodeHostFactory interface {
	Create(ctx context.Context, accessToken string) (CodeHost, error)
}

type CodeHost interface {
	ProduceCommits(ctx context.Context, repoURL string, lastAnalyzedCommitSHA string, commits chan<- string, errors chan<- string)
	GetCommitDiff(ctx context.Context, repoURL string, commitSHA string) (string, error)

	RepositoryID(ctx context.Context, repoURL string) (int64, error)

	GetRepositoryData(ctx context.Context, repoURL string) (map[string]interface{}, error)
	GetCommitData(ctx context.Context, repoURL string, commitSHA string) (map[string]interface{}, error)
}
