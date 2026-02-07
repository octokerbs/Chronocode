package analysis

import (
	"time"
)

type Subcommit struct {
	ID          int64      `json:"id"`
	CreatedAt   *time.Time `json:"createdAt"`
	Title       string     `json:"title"`
	Idea        string     `json:"idea"`
	Description string     `json:"description"`
	CommitSHA   string     `json:"commitSha"`
	Type        string     `json:"type"`
	Epic        string     `json:"epic"`
	Files       []string   `json:"files"`
}

func (s *Subcommit) ApplyAnalysis(commitSHA string, commitTimeStamp *time.Time) {
	s.CommitSHA = commitSHA
	s.CreatedAt = commitTimeStamp
}
