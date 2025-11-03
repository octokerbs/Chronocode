package application

import (
	"context"
	"strconv"

	"github.com/octokerbs/chronocode-backend/internal/domain"
)

type TimelineService struct {
	db     domain.Database
	logger domain.Logger
}

func NewTimelineService(db domain.Database, logger domain.Logger) *TimelineService {
	return &TimelineService{
		db:     db,
		logger: logger,
	}
}

func (s *TimelineService) GetSubcommitsFromRepo(ctx context.Context, repoIDStr string) ([]*domain.Subcommit, error) {
	log := s.logger.With("repoID", repoIDStr)

	repoID, err := strconv.ParseInt(repoIDStr, 10, 64)
	if err != nil {
		log.Error("Invalid repoID format", err)
		return nil, domain.NewError(domain.ErrBadRequest, err)
	}

	subcommits, err := s.db.GetSubcommitsByRepoID(ctx, repoID)
	if err != nil {
		log.Error("Failed to get subcommits from database", err)
		return nil, err
	}

	log.Info("Successfully fetched subcommits for timeline", "count", len(subcommits))
	return subcommits, nil
}
