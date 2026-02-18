package adapters

import (
	"context"

	"github.com/octokerbs/chronocode/internal/domain/repo"
)

type RepoRepository struct {
	repos map[string]repo.Repo
}

func NewRepoRepository() *RepoRepository {
	return &RepoRepository{map[string]repo.Repo{}}
}

func (r *RepoRepository) GetRepo(ctx context.Context, url string) (*repo.Repo, error) {
	obtainedRepo, ok := r.repos[url]
	if !ok {
		return &repo.Repo{}, repo.ErrRepositoryNotFound
	}

	return &obtainedRepo, nil
}

func (r *RepoRepository) GetRepoByID(ctx context.Context, id int64) (*repo.Repo, error) {
	for _, rp := range r.repos {
		if rp.ID() == id {
			return &rp, nil
		}
	}
	return &repo.Repo{}, repo.ErrRepositoryNotFound
}

func (r *RepoRepository) ListRepos(ctx context.Context) ([]*repo.Repo, error) {
	var result []*repo.Repo
	for _, rp := range r.repos {
		rp := rp
		result = append(result, &rp)
	}
	return result, nil
}

func (r *RepoRepository) StoreRepo(ctx context.Context, aRepo *repo.Repo) error {
	r.repos[aRepo.URL()] = *aRepo
	return nil
}
