package app

import (
	"github.com/octokerbs/chronocode-backend/internal/app/command"
	"github.com/octokerbs/chronocode-backend/internal/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	AnalyzeRepo command.AnalyzeRepoHandler
}

type Queries struct {
	IsRepoAnalyzed query.IsRepoAnalyzedHandler
	RepoSubcommits query.RepoSubcommitsHandler
}
