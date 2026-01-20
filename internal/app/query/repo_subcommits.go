package query

import (
	"context"

	"github.com/octokerbs/chronocode-backend/common/decorator"
	"go.uber.org/zap"
)

type RepoSubcommits struct {
}

type RepoSubcommitsHandler decorator.QueryHandler[RepoSubcommits, bool]

type repoSubcommitsHandler struct {
}

func NewRepoSubcommitsHandler(logger *zap.Logger) RepoSubcommitsHandler {
	return decorator.ApplyQueryDecorators[RepoSubcommits, bool](
		repoSubcommitsHandler{},
		logger,
	)
}

func (h repoSubcommitsHandler) Handle(ctx context.Context, query RepoSubcommits) (bool, error) {
	return true, nil
}
