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
}

func (c *Commit) IsDatabaseRecord() {}
