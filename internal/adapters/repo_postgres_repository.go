package adapters

import (
	"context"
	"database/sql"
	"errors"
	"time"

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
	const query = `SELECT id, name, url, last_analyzed_commit_sha, created_at FROM repository WHERE url = $1`

	var id int64
	var name, repoURL, lastSHA string
	var createdAt time.Time
	err := pg.db.QueryRowContext(ctx, query, url).Scan(&id, &name, &repoURL, &lastSHA, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repo.ErrRepositoryNotFound
		}
		return nil, err
	}

	return repo.NewRepo(id, name, repoURL, lastSHA, createdAt), nil
}

func (pg *PostgresRepoRepository) GetRepoByID(ctx context.Context, id int64) (*repo.Repo, error) {
	const query = `SELECT id, name, url, last_analyzed_commit_sha, created_at FROM repository WHERE id = $1`

	var repoID int64
	var name, repoURL, lastSHA string
	var createdAt time.Time
	err := pg.db.QueryRowContext(ctx, query, id).Scan(&repoID, &name, &repoURL, &lastSHA, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repo.ErrRepositoryNotFound
		}
		return nil, err
	}

	return repo.NewRepo(repoID, name, repoURL, lastSHA, createdAt), nil
}

func (pg *PostgresRepoRepository) ListRepos(ctx context.Context) ([]*repo.Repo, error) {
	const query = `SELECT id, name, url, last_analyzed_commit_sha, created_at FROM repository ORDER BY created_at DESC`

	rows, err := pg.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var repos []*repo.Repo
	for rows.Next() {
		var id int64
		var name, repoURL, lastSHA string
		var createdAt time.Time
		if err := rows.Scan(&id, &name, &repoURL, &lastSHA, &createdAt); err != nil {
			return nil, err
		}
		repos = append(repos, repo.NewRepo(id, name, repoURL, lastSHA, createdAt))
	}

	return repos, rows.Err()
}

func (pg *PostgresRepoRepository) StoreRepo(ctx context.Context, r *repo.Repo) error {
	const query = `
		INSERT INTO repository (id, name, url, last_analyzed_commit_sha, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			url = EXCLUDED.url,
			last_analyzed_commit_sha = EXCLUDED.last_analyzed_commit_sha`

	_, err := pg.db.ExecContext(ctx, query, r.ID(), r.Name(), r.URL(), r.LastAnalyzedCommitSHA(), r.CreatedAt())
	return err
}
