package repo

import "time"

type Repo struct {
	id                    int64
	name                  string
	url                   string
	lastAnalyzedCommitSHA string
	createdAt             time.Time
}

func NewRepo(id int64, name, url, lastAnalyzedCommit string, createdAt time.Time) *Repo {
	return &Repo{id, name, url, lastAnalyzedCommit, createdAt}
}

// IsURL
// Testing method to avoid breaking encapsulation
func (r *Repo) IsURL(url string) bool {
	return r.url == url
}

func (r *Repo) Name() string {
	return r.name
}

func (r *Repo) URL() string {
	return r.url
}

func (r *Repo) ID() int64 {
	return r.id
}

func (r *Repo) CreatedAt() time.Time {
	return r.createdAt
}

func (r *Repo) LastAnalyzedCommitSHA() string {
	return r.lastAnalyzedCommitSHA
}

func (r *Repo) SetLastAnalyzedCommitSHA(sha string) {
	r.lastAnalyzedCommitSHA = sha
}
