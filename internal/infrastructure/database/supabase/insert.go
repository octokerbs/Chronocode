package supabase

import (
	"context"

	"github.com/octokerbs/chronocode-backend/internal/domain"
)

func (s *SupabaseService) InsertRepository(ctx context.Context, repo *domain.RepositoryRecord) error {
	var results []domain.RepositoryRecord
	err := s.client.DB.From("repositories").Upsert(repo).Execute(&results)
	if err != nil {
		return err
	}
	return nil
}

func (s *SupabaseService) InsertCommit(ctx context.Context, commit *domain.CommitRecord) error {
	var results []domain.CommitRecord
	err := s.client.DB.From("commits").Upsert(commit).Execute(&results)
	if err != nil {
		return err
	}
	return nil
}

func (s *SupabaseService) InsertSubcommit(ctx context.Context, subcommit *domain.SubcommitRecord) error {
	var results []domain.SubcommitRecord
	err := s.client.DB.From("subcommits").Upsert(subcommit).Execute(&results)
	if err != nil {
		return err
	}
	return nil
}
