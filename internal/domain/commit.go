package domain

import (
	"time"
)

type Commit struct {
	SHA         string
	CreatedAt   *time.Time
	Author      string
	Date        string
	Message     string
	URL         string
	AuthorEmail string
	Description string
	AuthorURL   string
	Files       []string
	RepoID      int64
	Subcommits  []*Subcommit
}

func (c *Commit) ApplyAnalysis(analysis *CommitAnalysis) {
	c.Description = analysis.Commit.Description

	now := time.Now()
	subcommits := analysis.Subcommits

	for i := range subcommits {
		subcommits[i].ApplyAnalysis(c.SHA, &now)
		c.Subcommits = append(c.Subcommits, &subcommits[i])
	}
}
