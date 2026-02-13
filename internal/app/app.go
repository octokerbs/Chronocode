package app

import "github.com/octokerbs/chronocode/internal/app/command"

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	AnalyzeRepo command.AnalyzeRepo
}

type Queries struct {
}
