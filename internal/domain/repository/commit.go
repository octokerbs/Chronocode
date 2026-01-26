package repository

import "time"

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

type CommitAnalysis struct {
	Description string
	Subcommits  []Subcommit
}

func (c *Commit) ApplyAnalysis(analysis *CommitAnalysis) {
	c.Description = analysis.Description

	for i := range analysis.Subcommits {
		sc := &analysis.Subcommits[i]
		sc.ApplyCommitInfo(c.SHA, c.CreatedAt)
		c.Subcommits = append(c.Subcommits, sc)
	}
}
