package application

import (
	"context"
	"strconv"

	"github.com/octokerbs/chronocode-backend/internal/domain/analysis"
	"github.com/octokerbs/chronocode-backend/internal/domain/database"
	"github.com/octokerbs/chronocode-backend/internal/errors"
	"github.com/octokerbs/chronocode-backend/internal/log"
)

type Querier struct {
	db     database.Database
	logger log.Logger
}

func NewQuerier(db database.Database, logger log.Logger) *Querier {
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
		return nil, errors.NewError(errors.ErrBadRequest, err)
	}

	subcommits, err := ts.db.GetSubcommitsByRepoID(ctx, repoID)
	if err != nil {
		log.Error("Failed to get subcommits from database", err)
		return nil, err
	}

	log.Info("Successfully fetched subcommits for timeline", "count", len(subcommits))
	return subcommits, nil
}
