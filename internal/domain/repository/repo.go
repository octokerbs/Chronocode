package repository

import "time"

type Repo struct {
	ID                 int64
	CreatedAt          *time.Time
	Name               string
	URL                string
	LastAnalyzedCommit string
}

func (r *Repo) UpdateLastAnalyzedCommit(commitSHA string) {
	r.LastAnalyzedCommit = commitSHA
}
