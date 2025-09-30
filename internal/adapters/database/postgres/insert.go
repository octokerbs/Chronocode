package postgres

import (
	"context"
	"fmt"

	"github.com/octokerbs/chronocode-go/internal/domain"
)

func (p *PostgresClient) InsertRepository(ctx context.Context, repo *domain.RepositoryRecord) error {
	_, err := p.db.ExecContext(ctx, `INSERT INTO repository (id, created_at, name, url, last_analyzed_commit) VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, url = EXCLUDED.url, last_analyzed_commit = EXCLUDED.last_analyzed_commit`,
		repo.ID, repo.CreatedAt, repo.Name, repo.URL, repo.LastAnalyzedCommit)
	return err
}

func (p *PostgresClient) InsertCommit(ctx context.Context, commit *domain.CommitRecord) error {
	_, err := p.db.ExecContext(ctx, `
		INSERT INTO "commit" (
			sha, author, date, message, url, author_email, description, author_url, files, repo_id
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)
		ON CONFLICT (sha) DO UPDATE SET
			author = EXCLUDED.author,
			date = EXCLUDED.date,
			message = EXCLUDED.message,
			url = EXCLUDED.url,
			author_email = EXCLUDED.author_email,
			description = EXCLUDED.description,
			author_url = EXCLUDED.author_url,
			files = EXCLUDED.files,
			repo_id = EXCLUDED.repo_id
	`, commit.SHA, commit.Author, commit.Date, commit.Message, commit.URL, commit.AuthorEmail, commit.Description, commit.AuthorURL, commit.Files, commit.RepoID)

	if err != nil {
		return fmt.Errorf("insert commit failed: %w", err)
	}
	return nil
}

func (p *PostgresClient) InsertSubcommit(ctx context.Context, subcommit *domain.SubcommitRecord) error {
	_, err := p.db.ExecContext(ctx, `INSERT INTO subcommit (
		id, created_at, title, idea, description, commit_sha, type, epic, files
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9
	) ON CONFLICT (id) DO UPDATE SET
		title = EXCLUDED.title,
		idea = EXCLUDED.idea,
		description = EXCLUDED.description,
		commit_sha = EXCLUDED.commit_sha,
		type = EXCLUDED.type,
		epic = EXCLUDED.epic,
		files = EXCLUDED.files
	`,
		subcommit.ID, subcommit.CreatedAt, subcommit.Title, subcommit.Idea, subcommit.Description, subcommit.CommitSHA, subcommit.Type, subcommit.Epic, subcommit.Files)
	return err
}
