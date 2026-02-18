package query

import (
	"context"

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
	return h.repoRepository.ListRepos(ctx)
}
