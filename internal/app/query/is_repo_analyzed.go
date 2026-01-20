package query

import (
	"context"

	"github.com/octokerbs/chronocode-backend/common/decorator"
	"go.uber.org/zap"
)

type IsRepoAnalyzed struct {
}

type IsRepoAnalyzedHandler decorator.QueryHandler[IsRepoAnalyzed, bool]

type isRepoAnalyzedHandler struct {
}

func NewIsRepoAnalyzedHandler(logger *zap.Logger) IsRepoAnalyzedHandler {
	return decorator.ApplyQueryDecorators[IsRepoAnalyzed, bool](
		isRepoAnalyzedHandler{},
		logger,
	)
}

func (h isRepoAnalyzedHandler) Handle(ctx context.Context, query IsRepoAnalyzed) (bool, error) {
	return true, nil
}
