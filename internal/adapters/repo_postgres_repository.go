package adapters

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
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

	slog.Debug("Querying repository by URL", "url", url)

	var id int64
	var name, repoURL, lastSHA string
	var createdAt time.Time
	err := pg.db.QueryRowContext(ctx, query, url).Scan(&id, &name, &repoURL, &lastSHA, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Debug("Repository not found by URL", "url", url)
			return nil, repo.ErrRepositoryNotFound
		}
		slog.Error("Database error querying repository by URL", "url", url, "error", err)
		return nil, err
	}

	slog.Debug("Repository found by URL", "repo_id", id, "name", name)
	return repo.NewRepo(id, name, repoURL, lastSHA, createdAt), nil
}

func (pg *PostgresRepoRepository) GetRepoByID(ctx context.Context, id int64) (*repo.Repo, error) {
	const query = `SELECT id, name, url, last_analyzed_commit_sha, created_at FROM repository WHERE id = $1`

	slog.Debug("Querying repository by ID", "repo_id", id)

	var repoID int64
	var name, repoURL, lastSHA string
	var createdAt time.Time
	err := pg.db.QueryRowContext(ctx, query, id).Scan(&repoID, &name, &repoURL, &lastSHA, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Debug("Repository not found by ID", "repo_id", id)
			return nil, repo.ErrRepositoryNotFound
		}
		slog.Error("Database error querying repository by ID", "repo_id", id, "error", err)
		return nil, err
	}

	slog.Debug("Repository found by ID", "repo_id", repoID, "name", name)
	return repo.NewRepo(repoID, name, repoURL, lastSHA, createdAt), nil
}

func (pg *PostgresRepoRepository) ListRepos(ctx context.Context) ([]*repo.Repo, error) {
	const query = `SELECT id, name, url, last_analyzed_commit_sha, created_at FROM repository ORDER BY created_at DESC`

	slog.Debug("Listing all repositories from database")

	rows, err := pg.db.QueryContext(ctx, query)
	if err != nil {
		slog.Error("Database error listing repositories", "error", err)
		return nil, err
	}
	defer rows.Close()

	var repos []*repo.Repo
	for rows.Next() {
		var id int64
		var name, repoURL, lastSHA string
		var createdAt time.Time
		if err := rows.Scan(&id, &name, &repoURL, &lastSHA, &createdAt); err != nil {
			slog.Error("Database error scanning repository row", "error", err)
			return nil, err
		}
		repos = append(repos, repo.NewRepo(id, name, repoURL, lastSHA, createdAt))
	}

	slog.Debug("Repositories listed from database", "count", len(repos))
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

	slog.Debug("Storing repository", "repo_id", r.ID(), "name", r.Name(), "url", r.URL(), "last_sha", r.LastAnalyzedCommitSHA())

	_, err := pg.db.ExecContext(ctx, query, r.ID(), r.Name(), r.URL(), r.LastAnalyzedCommitSHA(), r.CreatedAt())
	if err != nil {
		slog.Error("Database error storing repository", "repo_id", r.ID(), "error", err)
		return err
	}

	slog.Info("Repository stored", "repo_id", r.ID(), "name", r.Name())
	return nil
}
