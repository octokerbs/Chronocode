package domain

import (
	"context"
	"time"
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
