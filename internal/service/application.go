package service

import (
	"context"
	"database/sql"
	"os"

	"github.com/google/generative-ai-go/genai"
	_ "github.com/lib/pq"
	"github.com/octokerbs/chronocode/internal/adapters"
	"github.com/octokerbs/chronocode/internal/app"
	"github.com/octokerbs/chronocode/internal/app/command"
	"github.com/octokerbs/chronocode/internal/app/query"

	"google.golang.org/api/option"
)

func NewApplication(ctx context.Context) app.Application {
	geminiClient, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		panic(err)
	}

	agent, err := adapters.NewGeminiAgent(geminiClient, os.Getenv("GEMINI_GENERATIVE_MODEL"))
	if err != nil {
		panic(err)
	}

	postgresClient, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	if err := postgresClient.Ping(); err != nil {
		panic(err)
	}

	repoRepository, err := adapters.NewPostgresRepoRepository(postgresClient)
	if err != nil {
		panic(err)
	}

	subcommitRepository, err := adapters.NewPostgresSubcommitRepository(postgresClient)
	if err != nil {
		panic(err)
	}

	codeHostFactory := adapters.NewGithubCodeHostFactory()
	locker := adapters.NewInMemoryLocker()

	return app.Application{
		Commands: app.Commands{
			AnalyzeRepo: command.NewAnalyzeRepoHandler(repoRepository, subcommitRepository, agent, codeHostFactory, locker),
		},
		Queries: app.Queries{
			GetSubcommits:   query.NewGetSubcommitsHandler(repoRepository, subcommitRepository, codeHostFactory),
			GetRepos:        query.NewGetReposHandler(repoRepository),
			GetUserProfile:  query.NewGetUserProfileHandler(codeHostFactory),
			SearchUserRepos: query.NewSearchUserReposHandler(codeHostFactory),
		},
	}
}
