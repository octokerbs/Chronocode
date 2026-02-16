package agent

import (
	"context"

	"github.com/octokerbs/chronocode/internal/domain/subcommit"
)

type Agent interface {
	AnalyzeCommitsIntoSubcommits(ctx context.Context, commitSHAs <-chan string, subcommits chan<- subcommit.Subcommit)
}
