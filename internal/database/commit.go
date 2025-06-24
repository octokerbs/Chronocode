package database

import (
	"context"
	"fmt"
	"time"

	"github.com/chrono-code-hackathon/chronocode-go/internal/agent"
	"github.com/chrono-code-hackathon/chronocode-go/internal/sourcecodehost"
)

type CommitRecord struct {
	SHA         string     `json:"sha"`                  // Completed manually via API data
	CreatedAt   *time.Time `json:"created_at,omitempty"` // Completed by Supabase
	Author      string     `json:"author"`               // Completed manually via API data
	Date        string     `json:"date"`                 // Completed manually via API data
	Message     string     `json:"message"`              // Completed manually via API data
	URL         string     `json:"url"`                  // Completed manually via API data
	AuthorEmail string     `json:"author_email"`         // Completed manually via API data
	Description string     `json:"description"`          // Comleted by code analysis
	AuthorURL   string     `json:"author_url"`           // Completed manually via API data
	Files       []string   `json:"files"`                // Completed manually via API data
	RepoID      int64      `json:"repo_id"`              // Completed manually via API data
}

func NewCommitRecord(ctx context.Context, sourceCodeService sourcecodehost.SourcecodeHostService, commitSHA string, commitAnalysis *agent.CommitSchema) (*CommitRecord, error) {
	commitData := sourceCodeService.GetCommitData(ctx, commitSHA)

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
		Description: commitAnalysis.Description,
		AuthorURL:   author_url,
		Files:       files,
		RepoID:      repositoryId,
	}, nil
}

func (cr *CommitRecord) InsertIntoDatabase(ctx context.Context, databaseService DatabaseService) error {
	err := databaseService.InsertCommit(ctx, cr)
	return err
}
