package domain

import (
	"context"
	"fmt"
	"time"
)

type Database interface {
	InsertRepository(ctx context.Context, repo *RepositoryRecord) error
	InsertCommit(ctx context.Context, commit *CommitRecord) error
	InsertSubcommit(ctx context.Context, subcommit *SubcommitRecord) error

	GetRepository(ctx context.Context, id int64) (*RepositoryRecord, bool, error)
}

type Record interface {
	InsertIntoDatabase(ctx context.Context, database Database) error
}

type CommitRecord struct {
	SHA         string     `json:"sha" sql:"sha"`                         // Completed manually via API data
	CreatedAt   *time.Time `json:"created_at,omitempty" sql:"created_at"` // Completed by Supabase/POSTGRES
	Author      string     `json:"author" sql:"author"`                   // Completed manually via API data
	Date        string     `json:"date" sql:"date"`                       // Completed manually via API data
	Message     string     `json:"message" sql:"message"`                 // Completed manually via API data
	URL         string     `json:"url" sql:"url"`                         // Completed manually via API data
	AuthorEmail string     `json:"author_email" sql:"author_email"`       // Completed manually via API data
	Description string     `json:"description" sql:"description"`         // Completed by code analysis
	AuthorURL   string     `json:"author_url" sql:"author_url"`           // Completed manually via API data
	Files       []string   `json:"files" sql:"files"`                     // Completed manually via API data
	RepoID      int64      `json:"repo_id" sql:"repo_id"`                 // Completed manually via API data
}

func NewCommitRecord(ctx context.Context, repoURL string, codeHost CodeHost, commitSHA string, commit *Commit) (*CommitRecord, error) {
	commitData, err := codeHost.GetCommitData(ctx, repoURL, commitSHA)
	if err != nil {
		return nil, err
	}

	author, ok := commitData["author"].(string)
	if !ok {
		return nil, fmt.Errorf("bad data type from Source Code service when fetching AUTHOR from commit data")
	}

	date, ok := commitData["date"].(string)
	if !ok {
		return nil, fmt.Errorf("bad data type from Source Code service when fetching DATE from commit data")
	}

	message, ok := commitData["message"].(string)
	if !ok {
		return nil, fmt.Errorf("bad data type from Source Code service when fetching MESSAGE from commit data")
	}

	url, ok := commitData["url"].(string)
	if !ok {
		return nil, fmt.Errorf("bad data type from Source Code service when fetching URL from commit data")
	}

	authorEmail, ok := commitData["author_email"].(string)
	if !ok {
		return nil, fmt.Errorf("bad data type from Source Code service when fetching AUTHOR EMAIL from commit data")
	}

	author_url, ok := commitData["author_url"].(string)
	if !ok {
		return nil, fmt.Errorf("bad data type from Source Code service when fetching AUTHOR URL from commit data")
	}

	files, ok := commitData["files"].([]string)
	if !ok {
		return nil, fmt.Errorf("bad data type from Source Code service when fetching FILES from commit data")
	}

	repositoryId, ok := commitData["repository_id"].(int64)
	if !ok {
		return nil, fmt.Errorf("bad data type from Source Code service when fetching REPOSITORY ID from commit data")
	}

	return &CommitRecord{
		SHA:         commitSHA,
		Author:      author,
		Date:        date,
		Message:     message,
		URL:         url,
		AuthorEmail: authorEmail,
		Description: commit.Description,
		AuthorURL:   author_url,
		Files:       files,
		RepoID:      repositoryId,
	}, nil
}

func (cr *CommitRecord) InsertIntoDatabase(ctx context.Context, database Database) error {
	err := database.InsertCommit(ctx, cr)
	return err
}

type RepositoryRecord struct {
	ID                 int64      `json:"id" db:"id"`
	CreatedAt          *time.Time `json:"created_at,omitempty" db:"created_at"`
	Name               string     `json:"name" db:"name"`
	URL                string     `json:"url" db:"url"`
	LastAnalyzedCommit string     `json:"last_analyzed_commit" db:"last_analyzed_commit"`
}

func NewRepositoryRecord(ctx context.Context, repoURL string, codeHost CodeHost) (*RepositoryRecord, error) {
	repoData, err := codeHost.GetRepositoryData(ctx, repoURL)
	if err != nil {
		return nil, err
	}

	id, ok := repoData["id"].(int64)
	if !ok {
		return nil, fmt.Errorf("bad data type from Source Code service when fetching ID from repository data")
	}

	name, ok := repoData["name"].(string)
	if !ok {
		return nil, fmt.Errorf("bad data type from Source Code service when fetching NAME from repository data")

	}

	return &RepositoryRecord{
		ID:                 id,
		Name:               name,
		URL:                repoURL,
		LastAnalyzedCommit: "",
	}, nil
}

func (rr *RepositoryRecord) InsertIntoDatabase(ctx context.Context, database Database) error {
	err := database.InsertRepository(ctx, rr)
	return err
}

type SubcommitRecord struct {
	ID          int64      `json:"id,omitempty" sql:"id"`                 // Automatically set by Supabase
	CreatedAt   *time.Time `json:"created_at,omitempty" sql:"created_at"` // Automatically set by Supabase
	Title       string     `json:"title" sql:"title"`
	Idea        string     `json:"idea" sql:"idea"`
	Description string     `json:"description" sql:"description"`
	CommitSHA   string     `json:"commit_sha" sql:"commit_sha"` // Completed manually when creating the subcommit
	Type        string     `json:"type" sql:"type"`
	Epic        string     `json:"epic" sql:"epic"`
	Files       []string   `json:"files" sql:"files"`
}

func NewSubcommitRecord(commitSHA string, subcommit *Subcommit) *SubcommitRecord {
	return &SubcommitRecord{
		Title:       subcommit.Title,
		Idea:        subcommit.Idea,
		Description: subcommit.Description,
		CommitSHA:   commitSHA,
		Type:        subcommit.Type,
		Epic:        subcommit.Epic,
		Files:       subcommit.Files,
	}
}

func (sr *SubcommitRecord) InsertIntoDatabase(ctx context.Context, database Database) error {
	err := database.InsertSubcommit(ctx, sr)
	return err
}
