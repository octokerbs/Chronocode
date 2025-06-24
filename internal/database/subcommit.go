package database

import (
	"context"
	"time"

	"github.com/chrono-code-hackathon/chronocode-go/internal/agent"
)

type SubcommitRecord struct {
	ID          int64      `json:"id,omitempty"`         // Automatically set by Supabase
	CreatedAt   *time.Time `json:"created_at,omitempty"` // Automatically set by Supabase
	Title       string     `json:"title"`
	Idea        string     `json:"idea"`
	Description string     `json:"description"`
	CommitSHA   string     `json:"commit_sha"` // Completed manually when creating the subcommit
	Type        string     `json:"type"`
	Epic        string     `json:"epic"`
	Files       []string   `json:"files"`
}

func NewSubcommitRecord(commitSHA string, subcommitAnalysis *agent.SubcommitSchema) *SubcommitRecord {
	return &SubcommitRecord{
		Title:       subcommitAnalysis.Title,
		Idea:        subcommitAnalysis.Idea,
		Description: subcommitAnalysis.Description,
		CommitSHA:   commitSHA,
		Type:        subcommitAnalysis.Type,
		Epic:        subcommitAnalysis.Epic,
		Files:       subcommitAnalysis.Files,
	}
}

func (sr *SubcommitRecord) InsertIntoDatabase(ctx context.Context, databaseService DatabaseService) error {
	err := databaseService.InsertSubcommit(ctx, sr)
	return err
}
