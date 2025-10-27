package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/octokerbs/chronocode-backend/internal/application"
	"github.com/octokerbs/chronocode-backend/internal/domain"
)

type PostgresClient struct {
	DB *sql.DB
}

func NewPostgresClient(ctx context.Context, dsn string) (*PostgresClient, error) {
	if dsn == "" {
		return nil, fmt.Errorf("postgres dsn is required")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &PostgresClient{DB: db}, nil
}

func (p *PostgresClient) GetRepository(ctx context.Context, id int64) (*domain.Repository, bool, error) {
	var repo domain.Repository
	row := p.DB.QueryRowContext(ctx, "SELECT * FROM repository WHERE id = $1", id)
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
	query := `
        INSERT INTO repositories (name, url, last_analyzed_commit, created_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at
    `
	if repo.CreatedAt == nil {
		t := time.Now()
		repo.CreatedAt = &t
	}

	var returnedID int64
	var returnedCreatedAt time.Time

	err := p.DB.QueryRowContext(ctx, query,
		repo.Name,
		repo.URL,
		repo.LastAnalyzedCommit,
		repo.CreatedAt,
	).Scan(&returnedID, &returnedCreatedAt)

	if err != nil {
		return fmt.Errorf("postgres: error al insertar repositorio: %w", err)
	}

	repo.ID = returnedID
	repo.CreatedAt = &returnedCreatedAt

	return nil
}

func (p *PostgresClient) InsertCommitRecord(ctx context.Context, commit *domain.Commit) error {
	query := `
        INSERT INTO commits (sha, created_at, author, message, url, author_email, description, author_url, files, repo_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `
	_, err := p.DB.ExecContext(ctx, query,
		commit.SHA,
		commit.CreatedAt,
		commit.Author,
		commit.Message,
		commit.URL,
		commit.AuthorEmail,
		commit.Description,
		commit.AuthorURL,
		commit.Files,
		commit.RepoID,
	)

	if err != nil {
		return fmt.Errorf("postgres: error al insertar commit: %w", err)
	}
	return nil
}

func (p *PostgresClient) InsertSubcommitRecord(ctx context.Context, subcommit *domain.Subcommit) error {
	query := `
        INSERT INTO subcommits (title, idea, description, commit_sha, type, epic, files)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at
    `
	var returnedID int64
	var returnedCreatedAt time.Time

	err := p.DB.QueryRowContext(ctx, query,
		subcommit.Title,
		subcommit.Idea,
		subcommit.Description,
		subcommit.CommitSHA,
		subcommit.Type,
		subcommit.Epic,
		subcommit.Files,
	).Scan(&returnedID, &returnedCreatedAt)

	if err != nil {
		return fmt.Errorf("postgres: error al insertar subcommit: %w", err)
	}

	subcommit.ID = returnedID
	subcommit.CreatedAt = &returnedCreatedAt

	return nil
}
