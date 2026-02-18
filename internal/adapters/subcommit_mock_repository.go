package adapters

import (
	"context"

	"github.com/octokerbs/chronocode/internal/domain/subcommit"
)

type SubcommitRepository struct {
	subcommits []subcommit.Subcommit
}

func NewSubcommitRepository() *SubcommitRepository {
	return &SubcommitRepository{}
}

func (s *SubcommitRepository) GetSubcommits(ctx context.Context, repoID int64) ([]subcommit.Subcommit, error) {
	repoSubcommits := []subcommit.Subcommit{}
	for _, sc := range s.subcommits {
		if sc.RepoID() == repoID {
			repoSubcommits = append(repoSubcommits, sc)
		}
	}

	return repoSubcommits, nil
}

func (s *SubcommitRepository) HasSubcommitsForCommit(ctx context.Context, repoID int64, commitSHA string) (bool, error) {
	for _, sc := range s.subcommits {
		if sc.RepoID() == repoID && sc.CommitSHA() == commitSHA {
			return true, nil
		}
	}
	return false, nil
}

func (s *SubcommitRepository) StoreSubcommits(ctx context.Context, subcommits <-chan subcommit.Subcommit) error {
	for sc := range subcommits {
		s.subcommits = append(s.subcommits, sc)
	}
	return nil
}
