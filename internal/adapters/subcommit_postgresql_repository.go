package adapters

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/octokerbs/chronocode-backend/internal/domain/subcommit"
)

type SubcommitPostgresRepository struct {
	db *sql.DB
}

func NewSubcommitPostgresRepository(db *sql.DB) *SubcommitPostgresRepository {
	return &SubcommitPostgresRepository{db: db}
}

func (r *SubcommitPostgresRepository) GetByRepoID(ctx context.Context, repoID int64) ([]*subcommit.Subcommit, error) {
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

	rows, err := r.db.QueryContext(ctx, query, repoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subcommits []*subcommit.Subcommit

	for rows.Next() {
		var sc subcommit.Subcommit
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
			return nil, err
		}

		sc.Files = files
		subcommits = append(subcommits, &sc)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return subcommits, nil
}
