package service

import (
	"context"
	"database/sql"
	"os"

	_ "github.com/lib/pq"
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

	return app.Application{
		Commands: app.Commands{
			AnalyzeRepo: command.NewAnalyzeRepoHandler(logger),
		},
		Queries: app.Queries{
			IsRepoAnalyzed: query.NewIsRepoAnalyzedHandler(logger),
			RepoSubcommits: query.NewRepoSubcommitsHandler(logger),
		},
	}
}
