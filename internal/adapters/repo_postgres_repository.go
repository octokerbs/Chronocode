package adapters

import (
	"context"
	"database/sql"
	"errors"

	"github.com/octokerbs/chronocode/internal/domain/repo"
)

type PostgresRepoRepository struct {
	db *sql.DB
}

func NewPostgresRepoRepository(db *sql.DB) (*PostgresRepoRepository, error) {
	if db == nil {
		return nil, errors.New("missing postgres client")
	}

	return &PostgresRepoRepository{db: db}, nil
}

func (pg *PostgresRepoRepository) GetRepo(ctx context.Context, url string) (*repo.Repo, error) {
	const query = `SELECT id, name, url, last_analyzed_commit_sha FROM repository WHERE url = $1`

	var id int64
	var name, repoURL, lastSHA string
	err := pg.db.QueryRowContext(ctx, query, url).Scan(&id, &name, &repoURL, &lastSHA)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repo.ErrRepositoryNotFound
		}
		return nil, err
	}

	return repo.NewRepo(id, name, repoURL, lastSHA), nil
}

func (pg *PostgresRepoRepository) StoreRepo(ctx context.Context, r *repo.Repo) error {
	const query = `
		INSERT INTO repository (id, name, url, last_analyzed_commit_sha)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			url = EXCLUDED.url,
			last_analyzed_commit_sha = EXCLUDED.last_analyzed_commit_sha`

	_, err := pg.db.ExecContext(ctx, query, r.ID(), r.Name(), r.URL(), r.LastAnalyzedCommitSHA())
	return err
}
