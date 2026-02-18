package query

import (
	"context"

	"github.com/octokerbs/chronocode/internal/domain/codehost"
	"github.com/octokerbs/chronocode/internal/domain/repo"
	"github.com/octokerbs/chronocode/internal/domain/subcommit"
)

type GetSubcommits struct {
	RepoURL     string
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
	codeHost, err := gs.codeHostFactory.Create(ctx, cmd.AccessToken)
	if err != nil {
		return nil, err
	}

	if err := codeHost.CanAccessRepo(ctx, cmd.RepoURL); err != nil {
		return nil, err
	}

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
