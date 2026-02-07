package analysis

import (
	"time"
)

type Repository struct {
	ID                 int64      `json:"id"`
	CreatedAt          *time.Time `json:"createdAt"`
	Name               string     `json:"name"`
	URL                string     `json:"url"`
	LastAnalyzedCommit string     `json:"lastAnalyzedCommit"`
}

func (r *Repository) UpdateLastAnalyzedCommit(commitSHA string) {
	r.LastAnalyzedCommit = commitSHA
}
