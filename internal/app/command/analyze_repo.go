package command

import (
	"context"

	"github.com/octokerbs/chronocode-backend/common/decorator"
	"go.uber.org/zap"
)

type AnalyzeRepo struct {
}

type AnalyzeRepoHandler decorator.CommandHandler[AnalyzeRepo]

type analyzeRepoHandler struct {
}

func NewAnalyzeRepoHandler(logger *zap.Logger) AnalyzeRepoHandler {
	return decorator.ApplyCommandDecorators[AnalyzeRepo](
		analyzeRepoHandler{},
		logger,
	)
}

func (h analyzeRepoHandler) Handle(ctx context.Context, cmd AnalyzeRepo) error {
	return nil
}
