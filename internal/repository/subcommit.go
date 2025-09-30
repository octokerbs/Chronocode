package repository

import (
	"context"
	"time"

	"github.com/octokerbs/chronocode-go/internal/domain/agent"
)

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

func (sr *SubcommitRecord) InsertIntoDatabase(ctx context.Context, databaseService DatabaseClient) error {
	err := databaseService.InsertSubcommit(ctx, sr)
	return err
}
