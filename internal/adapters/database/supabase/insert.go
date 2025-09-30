package supabase

import (
	"context"

	"github.com/octokerbs/chronocode-go/internal/repository"
)

func (s *SupabaseService) InsertRepository(ctx context.Context, repo *repository.RepositoryRecord) error {
	var results []repository.RepositoryRecord
	err := s.client.DB.From("repositories").Upsert(repo).Execute(&results)
	if err != nil {
		return err
	}
	return nil
}

func (s *SupabaseService) InsertCommit(ctx context.Context, commit *repository.CommitRecord) error {
	var results []repository.CommitRecord
	err := s.client.DB.From("commits").Upsert(commit).Execute(&results)
	if err != nil {
		return err
	}
	return nil
}

func (s *SupabaseService) InsertSubcommit(ctx context.Context, subcommit *repository.SubcommitRecord) error {
	var results []repository.SubcommitRecord
	err := s.client.DB.From("subcommits").Upsert(subcommit).Execute(&results)
	if err != nil {
		return err
	}
	return nil
}
