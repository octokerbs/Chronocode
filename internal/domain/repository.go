package domain

import (
	"time"
)

type Repository struct {
	ID                 int64
	CreatedAt          *time.Time
	Name               string
	URL                string
	LastAnalyzedCommit string
}

func (r *Repository) IsDatabaseRecord() {}
