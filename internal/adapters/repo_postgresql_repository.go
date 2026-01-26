package adapters

import (
	"context"
	"database/sql"
	"errors"

	"github.com/octokerbs/chronocode-backend/internal/domain/repository"
)

type RepoPostgresRepository struct {
	db *sql.DB
}

func NewRepoPostgresRepository(db *sql.DB) *RepoPostgresRepository {
	return &RepoPostgresRepository{db: db}
}

func (r *RepoPostgresRepository) Get(ctx context.Context, id int64) (*repository.Repo, error) {
	const query = `
		SELECT id, created_at, name, url, last_analyzed_commit
		FROM repository
		WHERE id = $1`

	repository := &repository.Repo{}

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&repository.ID,
		&repository.CreatedAt,
		&repository.Name,
		&repository.URL,
		&repository.LastAnalyzedCommit,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return repository, nil
}

func (r *RepoPostgresRepository) Store(ctx context.Context, repository *repository.Repo) error {
	const upsertQuery = `
		INSERT INTO repository (id, created_at, name, url, last_analyzed_commit)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			created_at = EXCLUDED.created_at,
			name = EXCLUDED.name,
			url = EXCLUDED.url,
			last_analyzed_commit = EXCLUDED.last_analyzed_commit`

	_, err := r.db.ExecContext(ctx, upsertQuery,
		repository.ID,
		repository.CreatedAt,
		repository.Name,
		repository.URL,
		repository.LastAnalyzedCommit,
	)

	return err
}
