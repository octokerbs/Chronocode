package app

import (
	"github.com/octokerbs/chronocode/internal/app/command"
	"github.com/octokerbs/chronocode/internal/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	AnalyzeRepo command.AnalyzeRepoHandler
}

type Queries struct {
	GetSubcommits   query.GetSubcommitsHandler
	GetRepos        query.GetReposHandler
	GetUserProfile  query.GetUserProfileHandler
	SearchUserRepos query.SearchUserReposHandler
}
