package subcommit

type Subcommit struct {
	title             string
	description       string
	modification_type string
	commit_sha        string
	files             []string
	repo_id           int64 // FK
}

func NewSubcommit(title, description, modificationType, commitSHA string, files []string, repoID int64) Subcommit {
	return Subcommit{
		title:             title,
		description:       description,
		modification_type: modificationType,
		commit_sha:        commitSHA,
		files:             files,
		repo_id:           repoID,
	}
}

func (s *Subcommit) RepoID() int64 {
	return s.repo_id
}
