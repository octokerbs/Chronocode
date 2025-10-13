package postgres

import (
	"context"
	"database/sql"

	"github.com/octokerbs/chronocode-backend/internal/domain"
)

func (p *PostgresClient) GetRepository(ctx context.Context, id int64) (*domain.RepositoryRecord, bool, error) {
	var repo domain.RepositoryRecord
	row := p.db.QueryRowContext(ctx, "SELECT * FROM repository WHERE id = $1", id)
	err := row.Scan(
		&repo.ID,
		&repo.CreatedAt,
		&repo.Name,
		&repo.URL,
		&repo.LastAnalyzedCommit,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &repo, true, nil
}
