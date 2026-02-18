package subcommit

import "time"

type Subcommit struct {
	title            string
	description      string
	modificationType string
	commitSHA        string
	files            []string
	repoID           int64
	committedAt      time.Time
}

func NewSubcommit(title, description, modificationType, commitSHA string, files []string, repoID int64, committedAt time.Time) Subcommit {
	return Subcommit{
		title:            title,
		description:      description,
		modificationType: modificationType,
		commitSHA:        commitSHA,
		files:            files,
		repoID:           repoID,
		committedAt:      committedAt,
	}
}

func (s *Subcommit) Title() string {
	return s.title
}

func (s *Subcommit) Description() string {
	return s.description
}

func (s *Subcommit) ModificationType() string {
	return s.modificationType
}

func (s *Subcommit) Files() []string {
	return s.files
}

func (s *Subcommit) RepoID() int64 {
	return s.repoID
}

func (s *Subcommit) CommitSHA() string {
	return s.commitSHA
}

func (s *Subcommit) CommittedAt() time.Time {
	return s.committedAt
}
