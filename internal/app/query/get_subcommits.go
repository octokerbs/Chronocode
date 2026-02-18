package query

import (
	"context"

	"github.com/octokerbs/chronocode/internal/domain/codehost"
	"github.com/octokerbs/chronocode/internal/domain/repo"
	"github.com/octokerbs/chronocode/internal/domain/subcommit"
)

type GetSubcommits struct {
	RepoID      int64
	AccessToken string
}

type GetSubcommitsHandler struct {
	repoRepository      repo.Repository
	subcommitRepository subcommit.Repository
	codeHostFactory     codehost.CodeHostFactory
}

func NewGetSubcommitsHandler(repoRepository repo.Repository, subcommitRepository subcommit.Repository, codeHostFactory codehost.CodeHostFactory) GetSubcommitsHandler {
	return GetSubcommitsHandler{repoRepository: repoRepository, subcommitRepository: subcommitRepository, codeHostFactory: codeHostFactory}
}

func (gs *GetSubcommitsHandler) Handle(ctx context.Context, cmd GetSubcommits) ([]subcommit.Subcommit, error) {
	foundRepo, err := gs.repoRepository.GetRepoByID(ctx, cmd.RepoID)
	if err != nil {
		return nil, err
	}

	codeHost, err := gs.codeHostFactory.Create(ctx, cmd.AccessToken)
	if err != nil {
		return nil, err
	}

	if err := codeHost.CanAccessRepo(ctx, foundRepo.URL()); err != nil {
		return nil, err
	}

	repoSubcommits, err := gs.subcommitRepository.GetSubcommits(ctx, foundRepo.ID())
	if err != nil {
		return nil, err
	}

	return repoSubcommits, nil
}
