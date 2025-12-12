package domain

import (
	"context"

	"golang.org/x/oauth2"
)

type CodeHost interface {
	FetchRepository(ctx context.Context, repoURL string) (*Repository, error)
	FetchRepositoryID(ctx context.Context, repoURL string) (int64, error)
	FetchCommit(ctx context.Context, repoURL string, commitSHA string) (*Commit, error)
	FetchCommitDiff(ctx context.Context, repoURL string, commitSHA string) (string, error)

	ProduceCommitSHAs(ctx context.Context, repoURL string, lastAnalyzedCommitSHA string, commitSHAs chan<- string) (string, error)
}

type Auth interface {
	GetAuthURL(state string) string
	ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error)
}

type CodeHostFactory interface {
	Create(ctx context.Context, accessToken string) CodeHost
}
