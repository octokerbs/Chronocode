package adapters

import (
	"context"

	"github.com/octokerbs/chronocode/internal/domain/agent"
)

type Agent struct{}

func NewAgent() *Agent {
	return &Agent{}
}

func (a *Agent) AnalyzeDiff(ctx context.Context, diff string) ([]agent.AnalysisResult, error) {
	if diff == FailingDiff {
		return nil, agent.ErrAnalysisFailed
	}

	return []agent.AnalysisResult{
		{Title: "title", Idea: "idea", Description: "description", Epic: "epic", ModificationType: "FEATURE", Files: []string{}},
	}, nil
}
