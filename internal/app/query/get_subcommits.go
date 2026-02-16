package query

import (
	"context"

	"github.com/octokerbs/chronocode/internal/domain/repo"
	"github.com/octokerbs/chronocode/internal/domain/subcommit"
)

type GetSubcommits struct {
	RepoURL string
}

type GetSubcommitsHandler struct {
	repoRepository      repo.Repository
	subcommitRepository subcommit.Repository
}

func NewGetSubcommitsHandler(repoRepository repo.Repository, subcommitRepository subcommit.Repository) GetSubcommitsHandler {
	return GetSubcommitsHandler{repoRepository: repoRepository, subcommitRepository: subcommitRepository}
}

func (gs *GetSubcommitsHandler) Handle(ctx context.Context, cmd GetSubcommits) ([]subcommit.Subcommit, error) {
	foundRepo, err := gs.repoRepository.GetRepo(ctx, cmd.RepoURL)
	if err != nil {
		return nil, err
	}

	repoSubcommits, err := gs.subcommitRepository.GetSubcommits(ctx, foundRepo.ID())
	if err != nil {
		return nil, err
	}

	return repoSubcommits, nil
}
