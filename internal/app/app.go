package app

import (
	"github.com/octokerbs/chronocode/internal/app/command"
	"github.com/octokerbs/chronocode/internal/app/query"
	"github.com/octokerbs/chronocode/internal/domain/analysis"
)

type Application struct {
	Commands Commands
	Queries  Queries
	Locker   analysis.Locker
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
