package query

import (
	"context"
	"strconv"

	"github.com/octokerbs/chronocode-backend/internal/domain/analysis"
	"github.com/octokerbs/chronocode-backend/internal/domain/database"
	"github.com/octokerbs/chronocode-backend/pkg/errors"
	"github.com/octokerbs/chronocode-backend/pkg/log"
)

type QuerierService struct {
	db     database.Database
	logger log.Logger
}

func NewQuerier(db database.Database, logger log.Logger) *QuerierService {
	return &QuerierService{
		db:     db,
		logger: logger,
	}
}

func (ts *QuerierService) GetSubcommitsFromRepo(ctx context.Context, repoIDStr string) ([]*analysis.Subcommit, error) {
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
