package domain

type CommitAnalysis struct {
	Commit     Commit
	Subcommits []Subcommit
}
