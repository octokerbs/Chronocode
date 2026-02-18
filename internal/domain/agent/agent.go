package agent

import (
	"context"
	"errors"
)

var (
	ErrAnalysisFailed = errors.New("agent analysis failed")
)

type AnalysisResult struct {
	Title            string
	Idea             string
	Description      string
	Epic             string
	ModificationType string
	Files            []string
}

type Agent interface {
	AnalyzeDiff(ctx context.Context, diff string) ([]AnalysisResult, error)
}
