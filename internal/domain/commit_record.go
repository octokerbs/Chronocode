package domain

import (
	"context"
	"fmt"
	"time"
)

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
