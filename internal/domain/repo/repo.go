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
