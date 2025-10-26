package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/octokerbs/chronocode-backend/internal/application"
	"github.com/octokerbs/chronocode-backend/internal/domain"
)

type PostgresClient struct {
	db *sql.DB
}

func NewPostgresClient(ctx context.Context, dsn string) (*PostgresClient, error) {
	if dsn == "" {
		return nil, fmt.Errorf("postgres dsn is required")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &PostgresClient{db: db}, nil
}

func (p *PostgresClient) GetRepository(ctx context.Context, id int64) (*domain.Repository, bool, error) {
	var repo domain.Repository
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

// Breaks the "Remove IF's heuristic". It is an exception case. Domain entities shouldn't have infrastructure
// methods. They shouldn't know about this stuff.
func (p *PostgresClient) ProcessRecords(ctx context.Context, records <-chan application.DatabaseRecord, errors chan<- string) {
	for record := range records {
		switch rec := record.(type) {
		case *domain.Commit:
			err := p.InsertCommitRecord(ctx, rec)
			if err != nil {
				errors <- fmt.Sprintf("error uploading commit to database: %s", err.Error())
			}
		case *domain.Repository:
			err := p.InsertRepositoryRecord(ctx, rec)
			if err != nil {
				errors <- fmt.Sprintf("error uploading repository to database: %s", err.Error())
			}
		case *domain.Subcommit:
			err := p.InsertSubcommitRecord(ctx, rec)
			if err != nil {
				errors <- fmt.Sprintf("error uploading subcommit to database: %s", err.Error())
			}
		default:
			fmt.Println("Unknown data type")
		}
	}
}

func (p *PostgresClient) InsertRepositoryRecord(ctx context.Context, repo *domain.Repository) error {
	_, err := p.db.ExecContext(ctx, `INSERT INTO repository (id, created_at, name, url, last_analyzed_commit) VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, url = EXCLUDED.url, last_analyzed_commit = EXCLUDED.last_analyzed_commit`,
		repo.ID, repo.CreatedAt, repo.Name, repo.URL, repo.LastAnalyzedCommit)
	return err
}

func (p *PostgresClient) InsertCommitRecord(ctx context.Context, commit *domain.Commit) error {
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

func (p *PostgresClient) InsertSubcommitRecord(ctx context.Context, subcommit *domain.Subcommit) error {
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
