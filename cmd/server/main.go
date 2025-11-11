package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/octokerbs/chronocode-backend/internal/api/http"
	"github.com/octokerbs/chronocode-backend/internal/application"
	"github.com/octokerbs/chronocode-backend/internal/config"
	"github.com/octokerbs/chronocode-backend/internal/infrastructure/agent/gemini"
	"github.com/octokerbs/chronocode-backend/internal/infrastructure/auth/githubauth"
	"github.com/octokerbs/chronocode-backend/internal/infrastructure/codehost/githubapi"
	"github.com/octokerbs/chronocode-backend/internal/infrastructure/database/postgres"
	"github.com/octokerbs/chronocode-backend/internal/infrastructure/logging/zap"
)

func main() {
	// Load configs
	ctx := context.Background()
	_ = godotenv.Load()
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Build dependencies
	repositoryAnalyzer, querier, authService, logger := buildDependencies(ctx, cfg)
	server := http.NewHTTPServer(repositoryAnalyzer, querier, authService, ":8080", logger)

	// Server start
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("Starting server", "port", "8080:8080")
		if err := server.Run(); err != nil {
			logger.Fatal("Server failed to run", err)
		}
	}()

	<-quit

	// Server shutdown
	logger.Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server forced to shutdown", err)
	}

	logger.Info("Server exited gracefully.")
}

func buildDependencies(ctx context.Context, cfg *config.Config) (*application.Analyzer, *application.Querier, *application.Auth, *zap.ZapLogger) {
	logger, err := zap.NewZapLogger()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	db, err := postgres.NewPostgresDatabase(cfg.DatabaseURL)
	if err != nil {
		logger.Error("Failed to initialize database", err)
	}

	agent, err := gemini.NewGeminiAgent(ctx, cfg.GeminiAPIKey)
	if err != nil {
		logger.Error("Failed to initialize gemini agent", err)
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
		logger,
	)

	analyzer := application.NewAnalyzer(
		ctx,
		agent,
		codeHostFactory,
		db,
		logger,
	)

	return analyzer, querier, auth, logger
}
