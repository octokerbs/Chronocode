package adapters

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"github.com/octokerbs/chronocode/internal/domain/subcommit"
)

type PostgresSubcommitRepository struct {
	db *sql.DB
}

func NewPostgresSubcommitRepository(db *sql.DB) (*PostgresSubcommitRepository, error) {
	if db == nil {
		return nil, errors.New("missing postgres client")
	}

	return &PostgresSubcommitRepository{db: db}, nil
}

func (pg *PostgresSubcommitRepository) GetSubcommits(ctx context.Context, repoID int64) ([]subcommit.Subcommit, error) {
	const query = `
		SELECT id, title, idea, description, epic, modification_type, commit_sha, files, repo_id, committed_at
		FROM subcommit
		WHERE repo_id = $1
		ORDER BY committed_at DESC`

	rows, err := pg.db.QueryContext(ctx, query, repoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subcommits []subcommit.Subcommit
	for rows.Next() {
		var id, rID int64
		var title, idea, desc, epic, modType, sha string
		var files pq.StringArray
		var committedAt time.Time

		if err := rows.Scan(&id, &title, &idea, &desc, &epic, &modType, &sha, &files, &rID, &committedAt); err != nil {
			return nil, err
		}

		subcommits = append(subcommits, subcommit.NewSubcommitFromDB(id, title, idea, desc, epic, modType, sha, []string(files), rID, committedAt))
	}

	return subcommits, rows.Err()
}

func (pg *PostgresSubcommitRepository) HasSubcommitsForCommit(ctx context.Context, repoID int64, commitSHA string) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM subcommit WHERE repo_id = $1 AND commit_sha = $2)`

	var exists bool
	err := pg.db.QueryRowContext(ctx, query, repoID, commitSHA).Scan(&exists)
	return exists, err
}

func (pg *PostgresSubcommitRepository) StoreSubcommits(ctx context.Context, subcommits <-chan subcommit.Subcommit) error {
	const query = `
		INSERT INTO subcommit (title, idea, description, epic, modification_type, commit_sha, files, repo_id, committed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	for sc := range subcommits {
		_, err := pg.db.ExecContext(ctx, query,
			sc.Title(), sc.Idea(), sc.Description(), sc.Epic(), sc.ModificationType(), sc.CommitSHA(),
			pq.Array(sc.Files()), sc.RepoID(), sc.CommittedAt())
		if err != nil {
			return err
		}
	}

	return nil
}
