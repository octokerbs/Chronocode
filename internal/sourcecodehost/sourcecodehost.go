package sourcecodehost

import (
	"context"
)

type SourcecodeHostService interface {
	ProduceCommits(ctx context.Context, lastAnalyzedCommitSHA string, commits chan<- string, errors chan<- error)
	GetCommitDiff(ctx context.Context, commitSHA string) (string, error)

	RepositoryID() int64

	GetRepositoryData() map[string]interface{}
	GetCommitData(ctx context.Context, commitSHA string) map[string]interface{}
}
