package query

import (
	"context"
	"strconv"

	"github.com/octokerbs/chronocode-backend/internal/domain"
	"github.com/octokerbs/chronocode-backend/internal/domain/analysis"
	"github.com/octokerbs/chronocode-backend/internal/domain/store"
)

type Querier struct {
	db     store.Database
	logger domain.Logger
}

func NewQuerier(db store.Database, logger domain.Logger) *Querier {
	return &Querier{
		db:     db,
		logger: logger,
	}
}

func (ts *Querier) GetSubcommitsFromRepo(ctx context.Context, repoIDStr string) ([]*analysis.Subcommit, error) {
	log := ts.logger.With("repoID", repoIDStr)

	repoID, err := strconv.ParseInt(repoIDStr, 10, 64)
	if err != nil {
		log.Error("Invalid repoID format", err)
		return nil, domain.NewError(domain.ErrBadRequest, err)
	}

	subcommits, err := ts.db.GetSubcommitsByRepoID(ctx, repoID)
	if err != nil {
		log.Error("Failed to get subcommits from database", err)
		return nil, err
	}

	log.Info("Successfully fetched subcommits for timeline", "count", len(subcommits))
	return subcommits, nil
}
