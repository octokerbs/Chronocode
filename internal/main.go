package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/octokerbs/chronocode/internal/adapters"
	"github.com/octokerbs/chronocode/internal/application"
	"github.com/octokerbs/chronocode/internal/application/command"
	"github.com/octokerbs/chronocode/internal/application/query"

	"github.com/octokerbs/chronocode/internal/ports/http"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"google.golang.org/api/option"
)

func NewApplication(ctx context.Context) application.Application {
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

	return application.Application{
		Commands: application.Commands{
			AnalyzeRepo: command.NewAnalyzeRepoHandler(repoRepository, subcommitRepository, agent, codeHostFactory, locker),
		},
		Queries: application.Queries{
			GetSubcommits:   query.NewGetSubcommitsHandler(repoRepository, subcommitRepository, codeHostFactory),
			GetRepos:        query.NewGetReposHandler(repoRepository),
			GetUserProfile:  query.NewGetUserProfileHandler(codeHostFactory),
			SearchUserRepos: query.NewSearchUserReposHandler(codeHostFactory),
		},
		Locker: locker,
	}
}

func main() {
	logLevel := slog.LevelInfo
	if os.Getenv("LOG_LEVEL") == "debug" {
		logLevel = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})))

	slog.Info("Chronocode server starting", "log_level", logLevel.String())

	ctx := context.Background()
	application := NewApplication(ctx)

	oauthConfig := &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
		Scopes:       []string{"read:user", "user:email", "repo"},
		Endpoint:     github.Endpoint,
	}

	slog.Info("GitHub OAuth configured",
		"client_id", os.Getenv("GITHUB_CLIENT_ID"),
		"redirect_url", os.Getenv("GITHUB_REDIRECT_URL"),
	)

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := http.NewServer(application, oauthConfig, frontendURL, port)

	slog.Info("Chronocode server ready", "port", port, "frontend_url", frontendURL)
	log.Fatal(server.ListenAndServe())
}
