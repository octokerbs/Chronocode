package subcommit

type Subcommit struct {
	title             string
	description       string
	modification_type string
	commit_sha        string
	files             []string
	repo_id           int64 // FK
}

func (s *Subcommit) RepoID() int64 {
	return s.repo_id
}
