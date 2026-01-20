package adapters

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/octokerbs/chronocode-backend/internal/domain/commit"
	"github.com/octokerbs/chronocode-backend/internal/domain/subcommit"
)

type CommitPostgresRepository struct {
	db *sql.DB
}

func NewCommitPostgresRepository(db *sql.DB) *CommitPostgresRepository {
	return &CommitPostgresRepository{db: db}
}

func (r *CommitPostgresRepository) Store(ctx context.Context, c *commit.Commit) error {
	if c == nil {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	const sqlInsert = `
		INSERT INTO commit (
			sha, created_at, author, date, message, url,
			author_email, description, author_url, files, repo_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err = tx.ExecContext(ctx, sqlInsert,
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

	if err := r.storeSubcommits(tx, c.Subcommits); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *CommitPostgresRepository) storeSubcommits(tx *sql.Tx, subcommits []*subcommit.Subcommit) error {
	if len(subcommits) == 0 {
		return nil
	}

	columns := []string{
		"created_at", "title", "idea", "description",
		"commit_sha", "type", "epic", "files",
	}

	stmt, err := tx.Prepare(pq.CopyIn("subcommit", columns...))
	if err != nil {
		return err
	}

	for _, sc := range subcommits {
		_, err = stmt.Exec(
			sc.CreatedAt, sc.Title, sc.Idea, sc.Description,
			sc.CommitSHA, sc.Type, sc.Epic, pq.Array(sc.Files),
		)
		if err != nil {
			return err
		}
	}

	if _, err = stmt.Exec(); err != nil {
		return err
	}

	return stmt.Close()
}
