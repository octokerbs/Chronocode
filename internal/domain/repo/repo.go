package repo

type Repo struct {
	id                    int64
	name                  string
	url                   string
	lastAnalyzedCommitSHA string
}

func NewRepo(id int64, name, url, lastAnalyzedCommit string) *Repo {
	return &Repo{id, name, url, lastAnalyzedCommit}
}

// IsURL
// Testing method to avoid breaking encapsulation
func (r *Repo) IsURL(url string) bool {
	return r.url == url
}

func (r *Repo) URL() string {
	return r.url
}

func (r *Repo) ID() int64 {
	return r.id
}

func (r *Repo) LastAnalyzedCommitSHA() string {
	return r.lastAnalyzedCommitSHA
}

func (r *Repo) SetLastAnalyzedCommitSHA(sha string) {
	r.lastAnalyzedCommitSHA = sha
}
