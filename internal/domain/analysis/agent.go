package analysis

import (
	"context"
)

type Agent interface {
	AnalyzeCommitDiff(ctx context.Context, diff string) (CommitAnalysis, error)
}
