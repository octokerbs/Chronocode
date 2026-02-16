package adapters

import (
	"context"

	"github.com/octokerbs/chronocode/internal/domain/subcommit"
)

type SubcommitRepository struct {
	subcommits map[int64]subcommit.Subcommit
}

func NewSubcommitRepository() *SubcommitRepository {
	return &SubcommitRepository{}
}

func (s *SubcommitRepository) GetSubcommits(ctx context.Context, repoID int64) ([]subcommit.Subcommit, error) {
	repoSubcommits := []subcommit.Subcommit{}
	for k, v := range s.subcommits {
		if k == repoID {
			repoSubcommits = append(repoSubcommits, v)
		}
	}

	return repoSubcommits, nil
}

func (s *SubcommitRepository) StoreSubcommits(ctx context.Context, subcommits <-chan subcommit.Subcommit) {
	for sc := range subcommits {
		s.subcommits[sc.RepoID()] = sc
	}

}

func (s *SubcommitRepository) StoreSubcommit(ctx context.Context, subcommit subcommit.Subcommit) {
}
