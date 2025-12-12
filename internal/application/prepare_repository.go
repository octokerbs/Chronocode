package application

import (
	"context"
	"errors"

	"github.com/octokerbs/chronocode-backend/internal/domain"
	"github.com/octokerbs/chronocode-backend/internal/domain/codehost"
	pkg_errors "github.com/octokerbs/chronocode-backend/internal/errors"
)

type PrepareRepository struct {
	CodeHostFactory codehost.CodeHostFactory
	Database        domain.Database
}

func (pr *PrepareRepository) Execute(ctx context.Context, repoURL, accessToken string) (*domain.Repository, error) {
	codeHost := pr.CodeHostFactory.Create(ctx, accessToken)

	fetchedRepo, err := codeHost.FetchRepository(ctx, repoURL)
	if err != nil {
		return nil, err
	}

	repo, err := pr.Database.GetRepository(ctx, fetchedRepo.ID)
	if errors.Is(err, pkg_errors.ErrNotFound) {
		// Create new repository
		if err := pr.Database.StoreRepository(ctx, fetchedRepo); err != nil {
			return nil, err
		}
		return fetchedRepo, nil
	}

	return repo, err
}
