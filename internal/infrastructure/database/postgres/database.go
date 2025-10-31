package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/octokerbs/chronocode-backend/internal/domain"
)

type Database struct {
	postgres *postgresClient
}

func NewPostgresDatabase(connectionString string) (*Database, error) {
	client, err := newPostgresClient(connectionString)
	if err != nil {
		return nil, err
	}

	return &Database{client}, nil
}

func (d *Database) GetRepository(ctx context.Context, id int64) (*domain.Repository, bool, error) {
	const query = `
		SELECT id, created_at, name, url, last_analyzed_commit
		FROM repositories
		WHERE id = $1`

	repo := &domain.Repository{}

	row := d.postgres.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&repo.ID,
		&repo.CreatedAt,
		&repo.Name,
		&repo.URL,
		&repo.LastAnalyzedCommit,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}

		return nil, false, err
	}

	return repo, true, nil
}

func (d *Database) StoreRepository(ctx context.Context, repo *domain.Repository) error {
	if repo.ID == 0 {
		const insertQuery = `
			INSERT INTO repositories (created_at, name, url, last_analyzed_commit)
			VALUES ($1, $2, $3, $4)
			RETURNING id`

		err := d.postgres.DB.QueryRowContext(ctx, insertQuery,
			repo.CreatedAt,
			repo.Name,
			repo.URL,
			repo.LastAnalyzedCommit,
		).Scan(&repo.ID)

		if err != nil {
			return err
		}
		return nil
	}

	const updateQuery = `
		UPDATE repositories
		SET created_at = $1,
		    name = $2,
		    url = $3,
		    last_analyzed_commit = $4
		WHERE id = $5`

	_, err := d.postgres.DB.ExecContext(ctx, updateQuery,
		repo.CreatedAt,
		repo.Name,
		repo.URL,
		repo.LastAnalyzedCommit,
		repo.ID,
	)

	if err != nil {
		return err
	}
	return nil
}

func (d *Database) StoreCommits(ctx context.Context, commits []*domain.Commit) error {
	if len(commits) == 0 {
		return nil
	}

	tx, err := d.postgres.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	columns := []string{
		"sha", "created_at", "author", "date", "message", "url",
		"author_email", "description", "author_url", "files", "repo_id",
	}

	stmt, err := tx.Prepare(pq.CopyIn("commits", columns...))
	if err != nil {
		return err
	}

	for _, c := range commits {
		_, err = stmt.Exec(
			c.SHA,
			c.CreatedAt,
			c.Author,
			c.Date,
			c.Message,
			c.URL,
			c.AuthorEmail,
			c.Description,
			c.AuthorURL,
			pq.Array(c.Files),
			c.RepoID,
		)
		if err != nil {
			return err
		}
	}

	if _, err = stmt.Exec(); err != nil {
		return err
	}

	if err = stmt.Close(); err != nil {
		return err
	}

	return tx.Commit()
}

func (d *Database) StoreSubcommits(ctx context.Context, subcommits []*domain.Subcommit) error {
	if len(subcommits) == 0 {
		return nil
	}

	tx, err := d.postgres.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	columns := []string{
		"created_at", "title", "idea", "description",
		"commit_sha", "type", "epic", "files",
	}

	stmt, err := tx.Prepare(pq.CopyIn("subcommits", columns...))
	if err != nil {
		return err
	}

	for _, sc := range subcommits {
		_, err = stmt.Exec(
			sc.CreatedAt,
			sc.Title,
			sc.Idea,
			sc.Description,
			sc.CommitSHA,
			sc.Type,
			sc.Epic,
			pq.Array(sc.Files),
		)
		if err != nil {
			return err
		}
	}

	if _, err = stmt.Exec(); err != nil {
		return err
	}

	if err = stmt.Close(); err != nil {
		return err
	}

	return tx.Commit()
}

type postgresClient struct {
	DB *sql.DB
}

func newPostgresClient(dsn string) (*postgresClient, error) {
	if dsn == "" {
		return nil, fmt.Errorf("postgres dsn is required")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &postgresClient{DB: db}, nil
}
