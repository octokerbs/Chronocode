package command

import (
	"context"
	"errors"

	"github.com/octokerbs/chronocode/internal/domain/codehost"
	"github.com/octokerbs/chronocode/internal/domain/repo"
)

type AnalyzeRepo struct {
	RepoURL string
}
type AnalyzeRepoHandler struct {
	repoRepository repo.Repository
	codeHost       codehost.CodeHost
}

func NewAnalyzeRepoHandler(repoRepository repo.Repository, codeHost codehost.CodeHost) AnalyzeRepoHandler {
	return AnalyzeRepoHandler{repoRepository: repoRepository, codeHost: codeHost}
}

func (s *AnalyzeRepoHandler) Handle(ctx context.Context, cmd AnalyzeRepo) error {
	r, err := s.repoRepository.GetRepo(ctx, cmd.RepoURL)
	if err != nil {
		if !errors.Is(err, repo.ErrRepositoryNotFound) {
			return err
		}

		r, err = s.codeHost.CreateRepoFromURL(ctx, cmd.RepoURL)
		if err != nil {
			return err
		}
	}

	err = s.repoRepository.StoreRepo(ctx, r)
	return err
}
