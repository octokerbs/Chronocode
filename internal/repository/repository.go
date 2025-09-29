package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/octokerbs/chronocode-go/internal/domain/sourcecodehost"
)

type RepositoryRecord struct {
	ID                 int64      `json:"id"`
	CreatedAt          *time.Time `json:"created_at,omitempty"` // Completed via supabase
	Name               string     `json:"name"`
	URL                string     `json:"url"`
	LastAnalyzedCommit string     `json:"last_analyzed_commit"`
}

func NewRepositoryRecord(sourceCodeService sourcecodehost.SourcecodeHostService) (*RepositoryRecord, error) {
	repoData := sourceCodeService.GetRepositoryData()

	id, ok := repoData["id"].(int64)
	if !ok {
		return nil, fmt.Errorf("bad data type from Source Code service when fetching ID from repository data")
	}

	name, ok := repoData["name"].(string)
	if !ok {
		return nil, fmt.Errorf("bad data type from Source Code service when fetching NAME from repository data")

	}

	url, ok := repoData["url"].(string)
	if !ok {
		return nil, fmt.Errorf("bad data type from Source Code service when fetching URL from repository data")
	}

	return &RepositoryRecord{
		ID:                 id,
		Name:               name,
		URL:                url,
		LastAnalyzedCommit: "",
	}, nil
}

func (rr *RepositoryRecord) InsertIntoDatabase(ctx context.Context, databaseService DatabaseService) error {
	err := databaseService.InsertRepository(ctx, rr)
	return err
}
