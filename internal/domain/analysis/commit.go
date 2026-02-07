package analysis

import (
	"time"
)

type Commit struct {
	SHA         string       `json:"sha"`
	CreatedAt   *time.Time   `json:"createdAt"`
	Author      string       `json:"author"`
	Date        string       `json:"date"`
	Message     string       `json:"message"`
	URL         string       `json:"url"`
	AuthorEmail string       `json:"authorEmail"`
	Description string       `json:"description"`
	AuthorURL   string       `json:"authorUrl"`
	Files       []string     `json:"files"`
	RepoID      int64        `json:"repoId"`
	Subcommits  []*Subcommit `json:"subcommits"`
}

func (c *Commit) ApplyAnalysis(analysis *CommitAnalysis) {
	c.Description = analysis.Commit.Description

	subcommits := analysis.Subcommits

	for i := range subcommits {
		subcommits[i].ApplyAnalysis(c.SHA, c.CreatedAt)
		c.Subcommits = append(c.Subcommits, &subcommits[i])
	}
}
