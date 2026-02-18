package subcommit

type Subcommit struct {
	title            string
	description      string
	modificationType string
	commitSHA        string
	files            []string
	repoID           int64
}

func NewSubcommit(title, description, modificationType, commitSHA string, files []string, repoID int64) Subcommit {
	return Subcommit{
		title:            title,
		description:      description,
		modificationType: modificationType,
		commitSHA:        commitSHA,
		files:            files,
		repoID:           repoID,
	}
}

func (s *Subcommit) RepoID() int64 {
	return s.repoID
}
