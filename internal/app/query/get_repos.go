package query

import (
	"context"
	"log/slog"

	"github.com/octokerbs/chronocode/internal/domain/repo"
)

type GetRepos struct{}

type GetReposHandler struct {
	repoRepository repo.Repository
}

func NewGetReposHandler(repoRepository repo.Repository) GetReposHandler {
	return GetReposHandler{repoRepository: repoRepository}
}

func (h *GetReposHandler) Handle(ctx context.Context, _ GetRepos) ([]*repo.Repo, error) {
	slog.Info("GetRepos query received")

	repos, err := h.repoRepository.ListRepos(ctx)
	if err != nil {
		slog.Error("Failed to list repositories", "error", err)
		return nil, err
	}

	slog.Info("GetRepos query completed", "count", len(repos))
	return repos, nil
}
