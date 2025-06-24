package supabase

import (
	"context"

	"github.com/chrono-code-hackathon/chronocode-go/internal/database"
)

func (s *SupabaseService) InsertRepository(ctx context.Context, repo *database.RepositoryRecord) error {
	var results []database.RepositoryRecord
	err := s.client.DB.From("repositories").Upsert(repo).Execute(&results)
	if err != nil {
		return err
	}
	return nil
}

func (s *SupabaseService) InsertCommit(ctx context.Context, commit *database.CommitRecord) error {
	var results []database.CommitRecord
	err := s.client.DB.From("commits").Upsert(commit).Execute(&results)
	if err != nil {
		return err
	}
	return nil
}

func (s *SupabaseService) InsertSubcommit(ctx context.Context, subcommit *database.SubcommitRecord) error {
	var results []database.SubcommitRecord
	err := s.client.DB.From("subcommits").Upsert(subcommit).Execute(&results)
	if err != nil {
		return err
	}
	return nil
}
