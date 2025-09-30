package supabase

import (
	"context"
	"strconv"

	"github.com/octokerbs/chronocode-go/internal/repository"
)

func (s *SupabaseService) GetRepository(ctx context.Context, id int64) (*repository.RepositoryRecord, bool, error) {
	var results []repository.RepositoryRecord
	err := s.client.DB.From("repositories").Select("*").Eq("id", strconv.FormatInt(id, 10)).Execute(&results)
	if err != nil {
		return nil, false, err
	}

	if len(results) > 0 {
		return &results[0], true, nil
	}

	return nil, false, nil
}
