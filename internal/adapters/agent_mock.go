package adapters

import (
	"context"

	"github.com/octokerbs/chronocode/internal/domain/subcommit"
)

type Agent struct {
}

func NewAgent() *Agent {
	return &Agent{}
}

func (a *Agent) AnalyzeCommitsIntoSubcommits(ctx context.Context, commitSHAs <-chan string, subcommits chan<- subcommit.Subcommit) {
	for range commitSHAs {
		subcommits <- subcommit.Subcommit{}
	}

}
