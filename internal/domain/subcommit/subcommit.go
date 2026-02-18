package subcommit

import "time"

type Subcommit struct {
	id               int64
	title            string
	idea             string
	description      string
	epic             string
	modificationType string
	commitSHA        string
	files            []string
	repoID           int64
	committedAt      time.Time
}

func NewSubcommit(title, idea, description, epic, modificationType, commitSHA string, files []string, repoID int64, committedAt time.Time) Subcommit {
	return Subcommit{
		title:            title,
		idea:             idea,
		description:      description,
		epic:             epic,
		modificationType: modificationType,
		commitSHA:        commitSHA,
		files:            files,
		repoID:           repoID,
		committedAt:      committedAt,
	}
}

func NewSubcommitFromDB(id int64, title, idea, description, epic, modificationType, commitSHA string, files []string, repoID int64, committedAt time.Time) Subcommit {
	sc := NewSubcommit(title, idea, description, epic, modificationType, commitSHA, files, repoID, committedAt)
	sc.id = id
	return sc
}

func (s *Subcommit) ID() int64 {
	return s.id
}

func (s *Subcommit) Title() string {
	return s.title
}

func (s *Subcommit) Idea() string {
	return s.idea
}

func (s *Subcommit) Description() string {
	return s.description
}

func (s *Subcommit) Epic() string {
	return s.epic
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
