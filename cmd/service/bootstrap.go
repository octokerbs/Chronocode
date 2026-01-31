package main

import (
	"context"
	"database/sql"
	"os"

	"github.com/octokerbs/chronocode-backend/internal/adapters"
	"github.com/octokerbs/chronocode-backend/internal/app"
	"github.com/octokerbs/chronocode-backend/internal/app/command"
	"github.com/octokerbs/chronocode-backend/internal/app/query"
	"go.uber.org/zap"
)

func NewApplication(ctx context.Context, logger *zap.Logger) app.Application {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	repoRepository := adapters.NewRepoPostgresRepository(db)
	agent, err := adapters.NewGeminiAgent(ctx, os.Getenv("GEMINI_API_KEY"))
	if err != nil {
		panic(err)
	}

	codeHostFactory := adapters.NewGitHubCodeHostFactory()

	return app.Application{
		Commands: app.Commands{
			AnalyzeRepo: command.NewAnalyzeRepoHandler(repoRepository, agent, codeHostFactory, logger),
		},
		Queries: app.Queries{
			IsRepoAnalyzed: query.NewIsRepoAnalyzedHandler(logger),
			RepoSubcommits: query.NewRepoSubcommitsHandler(logger),
		},
	}
}
