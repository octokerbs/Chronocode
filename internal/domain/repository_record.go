package domain

import (
	"context"
	"fmt"
	"time"
)

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
