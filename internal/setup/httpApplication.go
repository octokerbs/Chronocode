package setup

import (
	"context"
	"fmt"

	"github.com/octokerbs/chronocode-backend/internal/api/http"
	"github.com/octokerbs/chronocode-backend/internal/application"
	"github.com/octokerbs/chronocode-backend/internal/config"
	"github.com/octokerbs/chronocode-backend/internal/domain"
	"github.com/octokerbs/chronocode-backend/internal/infrastructure/agent/gemini"
	"github.com/octokerbs/chronocode-backend/internal/infrastructure/codehost/githubapi"
	"github.com/octokerbs/chronocode-backend/internal/infrastructure/database/postgres"
	"github.com/octokerbs/chronocode-backend/internal/infrastructure/logging/zap"
)

type HTTPApplication struct {
	Config *config.HTTPConfig
	Server *http.Server
	Logger domain.Logger
}

func NewHTTPApplication(cfg *config.HTTPConfig) (*HTTPApplication, error) {
	ctx := context.Background()

	logger, err := zap.NewLogger()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}
	logger.Info("Logger initialized")

	logger.Info("Initializing infrastructure adapters...")
	geminiClient, err := gemini.NewGeminiAgent(ctx, cfg.GeminiAPIKey)
	if err != nil {
		return nil, fmt.Errorf("failed to init gemini client: %w", err)
	}

	githubClient := githubapi.NewGitHubFactory()

	postgresClient, err := postgres.NewPostgresDatabase(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to init postgres database: %w", err)
	}

	logger.Info("Initializing application layer...")
	repoAnalyzer := application.NewRepositoryAnalyzer(ctx, geminiClient, githubClient, postgresClient, logger)

	logger.Info("Initializing API server...")
	server := http.NewServer(cfg.Port, logger, repoAnalyzer)

	return &HTTPApplication{
		Config: cfg,
		Server: server,
		Logger: logger,
	}, nil
}

func (ha *HTTPApplication) Run() error {
	ha.Logger.Info("Starting server...", "port", ha.Config.Port)
	return ha.Server.Run()
}

func (ha *HTTPApplication) Shutdown(ctx context.Context) error {
	ha.Logger.Info("Shutting down server...")
	return ha.Server.Shutdown(ctx)
}
