package domain

import (
	"time"
)

type Subcommit struct {
	ID          int64
	CreatedAt   *time.Time
	Title       string
	Idea        string
	Description string
	CommitSHA   string
	Type        string
	Epic        string
	Files       []string
}

func (s *Subcommit) ApplyAnalysis(commitSHA string, commitTimeStamp *time.Time) {
	s.CommitSHA = commitSHA
	s.CreatedAt = commitTimeStamp
}
