package domain

import (
	"context"
)

type Agent interface {
	AnalyzeCommitDiff(ctx context.Context, diff string) (CommitAnalysis, error)
}
