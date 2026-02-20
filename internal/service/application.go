package service

import (
	"context"
	"database/sql"
	"log/slog"
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
	slog.Info("Initializing application dependencies")

	slog.Info("Connecting to Gemini AI")
	geminiClient, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		slog.Error("Failed to create Gemini client", "error", err)
		panic(err)
	}
	slog.Info("Gemini AI client connected")

	agent, err := adapters.NewGeminiAgent(geminiClient, os.Getenv("GEMINI_GENERATIVE_MODEL"))
	if err != nil {
		slog.Error("Failed to create Gemini agent", "error", err)
		panic(err)
	}

	slog.Info("Connecting to PostgreSQL")
	postgresClient, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		slog.Error("Failed to open PostgreSQL connection", "error", err)
		panic(err)
	}

	if err := postgresClient.Ping(); err != nil {
		slog.Error("Failed to ping PostgreSQL", "error", err)
		panic(err)
	}
	slog.Info("PostgreSQL connected successfully")

	repoRepository, err := adapters.NewPostgresRepoRepository(postgresClient)
	if err != nil {
		slog.Error("Failed to create repo repository", "error", err)
		panic(err)
	}

	subcommitRepository, err := adapters.NewPostgresSubcommitRepository(postgresClient)
	if err != nil {
		slog.Error("Failed to create subcommit repository", "error", err)
		panic(err)
	}

	codeHostFactory := adapters.NewGithubCodeHostFactory()
	locker := adapters.NewInMemoryLocker()

	slog.Info("All dependencies initialized successfully")

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
		Locker: locker,
	}
}
