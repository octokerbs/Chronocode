package query

import (
	"context"
	"log/slog"

	"github.com/octokerbs/chronocode/internal/domain/codehost"
	"github.com/octokerbs/chronocode/internal/domain/repo"
	"github.com/octokerbs/chronocode/internal/domain/subcommit"
)

type GetSubcommits struct {
	RepoID      int64
	AccessToken string
}

type GetSubcommitsResult struct {
	Subcommits []subcommit.Subcommit
	RepoURL    string
}

type GetSubcommitsHandler struct {
	repoRepository      repo.Repository
	subcommitRepository subcommit.Repository
	codeHostFactory     codehost.CodeHostFactory
}

func NewGetSubcommitsHandler(repoRepository repo.Repository, subcommitRepository subcommit.Repository, codeHostFactory codehost.CodeHostFactory) GetSubcommitsHandler {
	return GetSubcommitsHandler{repoRepository: repoRepository, subcommitRepository: subcommitRepository, codeHostFactory: codeHostFactory}
}

func (gs *GetSubcommitsHandler) Handle(ctx context.Context, cmd GetSubcommits) (GetSubcommitsResult, error) {
	slog.Info("GetSubcommits query received", "repo_id", cmd.RepoID)

	foundRepo, err := gs.repoRepository.GetRepoByID(ctx, cmd.RepoID)
	if err != nil {
		slog.Warn("Repository not found by ID", "repo_id", cmd.RepoID, "error", err)
		return GetSubcommitsResult{}, err
	}

	codeHost, err := gs.codeHostFactory.Create(ctx, cmd.AccessToken)
	if err != nil {
		slog.Error("Failed to create code host client for access check", "repo_id", cmd.RepoID, "error", err)
		return GetSubcommitsResult{}, err
	}

	if err := codeHost.CanAccessRepo(ctx, foundRepo.URL()); err != nil {
		slog.Warn("Access denied to repository", "repo_id", cmd.RepoID, "repo_url", foundRepo.URL(), "error", err)
		return GetSubcommitsResult{}, err
	}

	repoSubcommits, err := gs.subcommitRepository.GetSubcommits(ctx, foundRepo.ID())
	if err != nil {
		slog.Error("Failed to fetch subcommits from database", "repo_id", foundRepo.ID(), "error", err)
		return GetSubcommitsResult{}, err
	}

	slog.Info("GetSubcommits query completed", "repo_id", foundRepo.ID(), "count", len(repoSubcommits))
	return GetSubcommitsResult{
		Subcommits: repoSubcommits,
		RepoURL:    foundRepo.URL(),
	}, nil
}
