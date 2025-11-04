package domain

// DTO for the Agent port.
// The only purpose of this entity is to be a shell to trasnfer data from the api response to the application.
type CommitAnalysis struct {
	Commit     Commit
	Subcommits []Subcommit
}
