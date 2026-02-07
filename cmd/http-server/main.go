package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/octokerbs/chronocode-backend/config"
	"github.com/octokerbs/chronocode-backend/internal/api/http"
	"github.com/octokerbs/chronocode-backend/internal/application"
	"github.com/octokerbs/chronocode-backend/internal/infrastructure/agent/gemini"
	"github.com/octokerbs/chronocode-backend/internal/infrastructure/auth/githubauth"
	"github.com/octokerbs/chronocode-backend/internal/infrastructure/codehost/githubapi"
	"github.com/octokerbs/chronocode-backend/internal/infrastructure/database/postgres"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	_ = godotenv.Load()
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	repositoryAnalyzer, persistCommits, prepareRepo, querier, authService := buildDependencies(ctx, cfg)
	server := http.NewHTTPServer(repositoryAnalyzer, persistCommits, prepareRepo, querier, authService, ":8080", cfg.FrontendURL)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.Run(); err != nil {
			fmt.Printf("Server failed to start: %e", err)
		}
	}()

	<-quit

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
	}
}

func buildDependencies(ctx context.Context, cfg *config.Config) (*application.Analyzer, *application.PersistCommits, *application.PrepareRepository, *application.Querier, *application.Auth) {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("Error building logger")
	}

	db, err := postgres.NewPostgresDatabase(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Couldn't initialize postgres db", zap.String("buildDependencies", err.Error()))
	}

	agent, err := gemini.NewGeminiAgent(ctx, cfg.GeminiAPIKey)
	if err != nil {
		logger.Fatal("Couldn't initialize postgres gemini agent", zap.String("buildDependencies", err.Error()))
	}

	githubAuth := githubauth.NewGitHubAuth(
		cfg.GithubClientID,
		cfg.GithubClientSecret,
		cfg.RedirectURL,
	)

	auth := application.NewAuth(githubAuth)

	codeHostFactory := githubapi.NewGitHubCodeHostFactory()

	querier := application.NewQuerier(
		db,
	)

	analyzer := application.NewAnalyzer(
		agent,
		codeHostFactory,
	)

	persistCommits := &application.PersistCommits{Database: db}
	prepareRepo := &application.PrepareRepository{CodeHostFactory: codeHostFactory, Database: db}

	return analyzer, persistCommits, prepareRepo, querier, auth
}
