package repository

import "time"

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

func (s *Subcommit) ApplyCommitInfo(commitSHA string, commitTimestamp *time.Time) {
	s.CommitSHA = commitSHA
	s.CreatedAt = commitTimestamp
}
