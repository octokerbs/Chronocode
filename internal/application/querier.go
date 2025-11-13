package application

import (
	"context"
	"strconv"

	"github.com/octokerbs/chronocode-backend/internal/domain/analysis"
	"github.com/octokerbs/chronocode-backend/internal/domain/database"
	"github.com/octokerbs/chronocode-backend/internal/errors"
)

type Querier struct {
	db database.Database
}

func NewQuerier(db database.Database) *Querier {
	return &Querier{
		db: db,
	}
}

func (ts *Querier) GetSubcommitsFromRepo(ctx context.Context, repoIDStr string) ([]*analysis.Subcommit, error) {
	repoID, err := strconv.ParseInt(repoIDStr, 10, 64)
	if err != nil {
		return nil, errors.NewError(errors.ErrBadRequest, err)
	}

	subcommits, err := ts.db.GetSubcommitsByRepoID(ctx, repoID)
	if err != nil {
		return nil, err
	}

	return subcommits, nil
}
