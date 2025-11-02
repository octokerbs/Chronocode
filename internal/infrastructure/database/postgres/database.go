package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/octokerbs/chronocode-backend/internal/domain"
)

var errEmptyConnectionString = errors.New("database connection string is empty")

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

func (d *Database) GetRepository(ctx context.Context, id int64) (*domain.Repository, error) {
	const query = `
		SELECT id, created_at, name, url, last_analyzed_commit
		FROM repository
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
			return nil, domain.ErrBadRequest
		}

		return nil, domain.NewError(domain.ErrInternalFailure, err)
	}

	return repo, nil
}

func (d *Database) StoreRepository(ctx context.Context, repo *domain.Repository) error {
	const upsertQuery = `
        INSERT INTO repository (id, created_at, name, url, last_analyzed_commit)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (id) DO UPDATE SET
            created_at = EXCLUDED.created_at,
            name = EXCLUDED.name,
            url = EXCLUDED.url,
            last_analyzed_commit = EXCLUDED.last_analyzed_commit`

	_, err := d.postgres.DB.ExecContext(ctx, upsertQuery,
		repo.ID,
		repo.CreatedAt,
		repo.Name,
		repo.URL,
		repo.LastAnalyzedCommit,
	)

	if err != nil {
		return domain.NewError(domain.ErrInternalFailure, err)
	}

	return nil
}

func (d *Database) StoreCommits(ctx context.Context, commits []*domain.Commit) error {
	if len(commits) == 0 {
		return nil
	}

	tx, err := d.postgres.DB.BeginTx(ctx, nil)
	if err != nil {
		return domain.NewError(domain.ErrInternalFailure, err)
	}
	defer tx.Rollback()

	columns := []string{
		"sha", "created_at", "author", "date", "message", "url",
		"author_email", "description", "author_url", "files", "repo_id",
	}

	stmt, err := tx.Prepare(pq.CopyIn("commit", columns...))
	if err != nil {
		return domain.NewError(domain.ErrInternalFailure, err)
	}

	for _, c := range commits {
		_, err = stmt.Exec(
			c.SHA, c.CreatedAt, c.Author, c.Date, c.Message, c.URL,
			c.AuthorEmail, c.Description, c.AuthorURL, pq.Array(c.Files), c.RepoID,
		)
		if err != nil {
			return domain.NewError(domain.ErrInternalFailure, err)
		}
	}

	if _, err = stmt.Exec(); err != nil {
		return domain.NewError(domain.ErrInternalFailure, err)
	}

	if err = stmt.Close(); err != nil {
		return domain.NewError(domain.ErrInternalFailure, err)
	}

	for _, commit := range commits {
		if err := d.StoreSubcommits(ctx, tx, commit.Subcommits); err != nil {
			return domain.NewError(domain.ErrInternalFailure, err)
		}
	}

	return tx.Commit()
}

func (d *Database) StoreSubcommits(ctx context.Context, tx *sql.Tx, subcommits []*domain.Subcommit) error {
	if len(subcommits) == 0 {
		return nil
	}

	columns := []string{
		"created_at", "title", "idea", "description",
		"commit_sha", "type", "epic", "files",
	}

	stmt, err := tx.Prepare(pq.CopyIn("subcommit", columns...))
	if err != nil {
		return domain.NewError(domain.ErrInternalFailure, err)
	}

	for _, sc := range subcommits {
		_, err = stmt.Exec(
			sc.CreatedAt, sc.Title, sc.Idea, sc.Description,
			sc.CommitSHA, sc.Type, sc.Epic, pq.Array(sc.Files),
		)
		if err != nil {
			return domain.NewError(domain.ErrInternalFailure, err)
		}
	}

	if _, err = stmt.Exec(); err != nil {
		return domain.NewError(domain.ErrInternalFailure, err)
	}

	if err = stmt.Close(); err != nil {
		return domain.NewError(domain.ErrInternalFailure, err)
	}

	return nil
}

type postgresClient struct {
	DB *sql.DB
}

func newPostgresClient(dsn string) (*postgresClient, error) {
	if dsn == "" {
		return nil, errEmptyConnectionString
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &postgresClient{DB: db}, nil
}
