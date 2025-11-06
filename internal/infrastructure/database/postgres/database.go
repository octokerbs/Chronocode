package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/octokerbs/chronocode-backend/internal/domain"
	"github.com/octokerbs/chronocode-backend/internal/domain/analysis"
)

type Database struct {
	postgres *postgresClient
}

func NewPostgresDatabase(connectionString string) (*Database, error) {
	client, err := newPostgresClient(connectionString)
	if err != nil {
		return nil, domain.NewError(domain.ErrInternalFailure, err)
	}

	return &Database{client}, nil
}

func (d *Database) GetRepository(ctx context.Context, id int64) (*analysis.Repository, error) {
	const query = `
		SELECT id, created_at, name, url, last_analyzed_commit
		FROM repository
		WHERE id = $1`

	repo := &analysis.Repository{}

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
			return nil, domain.ErrNotFound
		}

		return nil, domain.NewError(domain.ErrInternalFailure, err)
	}

	return repo, nil
}

func (d *Database) StoreRepository(ctx context.Context, repo *analysis.Repository) error {
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

func (d *Database) StoreCommit(ctx context.Context, commit *analysis.Commit) error {
	if commit == nil {
		return nil
	}

	tx, err := d.postgres.DB.BeginTx(ctx, nil)
	if err != nil {
		return domain.NewError(domain.ErrInternalFailure, err)
	}
	defer tx.Rollback()

	const sqlInsert = `
        INSERT INTO commit (
            sha, created_at, author, date, message, url,
            author_email, description, author_url, files, repo_id
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    `

	_, err = tx.ExecContext(ctx, sqlInsert,
		commit.SHA,
		commit.CreatedAt,
		commit.Author,
		commit.Date,
		commit.Message,
		commit.URL,
		commit.AuthorEmail,
		commit.Description,
		commit.AuthorURL,
		pq.Array(commit.Files),
		commit.RepoID,
	)
	if err != nil {
		return domain.NewError(domain.ErrInternalFailure, err)
	}

	if err := d.storeSubcommits(tx, commit.Subcommits); err != nil {
		return domain.NewError(domain.ErrInternalFailure, err)
	}

	if err = tx.Commit(); err != nil {
		return domain.NewError(domain.ErrInternalFailure, err)
	}

	return nil
}

func (d *Database) storeSubcommits(tx *sql.Tx, subcommits []*analysis.Subcommit) error {
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

func (d *Database) GetSubcommitsByRepoID(ctx context.Context, repoID int64) ([]*analysis.Subcommit, error) {
	// Esta consulta une subcommit con commit para filtrar por repo_id
	const query = `
		SELECT 
			sc.id, sc.created_at, sc.title, sc.idea, sc.description,
			sc.commit_sha, sc.type, sc.epic, sc.files
		FROM 
			subcommit sc
		JOIN 
			commit c ON sc.commit_sha = c.sha
		WHERE 
			c.repo_id = $1
		ORDER BY 
			sc.created_at DESC`

	rows, err := d.postgres.DB.QueryContext(ctx, query, repoID)
	if err != nil {
		return nil, domain.NewError(domain.ErrInternalFailure, err)
	}
	defer rows.Close()

	var subcommits []*analysis.Subcommit

	for rows.Next() {
		var sc analysis.Subcommit
		var files pq.StringArray

		err := rows.Scan(
			&sc.ID,
			&sc.CreatedAt,
			&sc.Title,
			&sc.Idea,
			&sc.Description,
			&sc.CommitSHA,
			&sc.Type,
			&sc.Epic,
			&files,
		)
		if err != nil {
			return nil, domain.NewError(domain.ErrInternalFailure, err)
		}

		sc.Files = files
		subcommits = append(subcommits, &sc)
	}

	if err = rows.Err(); err != nil {
		return nil, domain.NewError(domain.ErrInternalFailure, err)
	}

	return subcommits, nil
}

type postgresClient struct {
	DB *sql.DB
}

func newPostgresClient(dsn string) (*postgresClient, error) {
	if dsn == "" {
		return nil, errors.New("database connection string is empty")
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
