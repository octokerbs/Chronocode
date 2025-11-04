package setup

import (
	"context"
	"log"

	"github.com/octokerbs/chronocode-backend/internal/api/http"
	"github.com/octokerbs/chronocode-backend/internal/config"
	"github.com/octokerbs/chronocode-backend/internal/domain"
	"github.com/octokerbs/chronocode-backend/internal/infrastructure/agent/gemini"
	"github.com/octokerbs/chronocode-backend/internal/infrastructure/codehost/githubapi"
	"github.com/octokerbs/chronocode-backend/internal/infrastructure/database/postgres"
	"github.com/octokerbs/chronocode-backend/internal/infrastructure/logging/zap"
	"github.com/octokerbs/chronocode-backend/internal/service"
)

type HTTPApplication struct {
	Config   *config.HTTPConfig
	Logger   domain.Logger
	Server   *http.Server
	DB       domain.Database
	Analyzer *service.RepositoryAnalyzerService
	Timeline *service.TimelineService
}

func NewHTTPApplication() *HTTPApplication {
	cfg, err := config.NewHTTPConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger, err := zap.NewLogger()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	ctx := context.Background()

	db, err := postgres.NewPostgresDatabase(cfg.DatabaseURL)
	if err != nil {
		logger.Error("Failed to initialize database", err)
	}

	agent, err := gemini.NewGeminiAgent(ctx, cfg.GeminiAPIKey)
	if err != nil {
		logger.Error("Failed to initialize gemini agent", err)
	}

	codeHostFactory := githubapi.NewGitHubFactory()

	analyzerService := service.NewRepositoryAnalyzer(
		ctx,
		agent,
		codeHostFactory,
		db,
		logger,
	)

	timelineService := service.NewTimelineService(
		db,
		logger,
	)

	server := http.NewServer(
		cfg.Port,
		logger,
		analyzerService,
		timelineService,
	)

	return &HTTPApplication{
		Config:   cfg,
		Logger:   logger,
		Server:   server,
		DB:       db,
		Analyzer: analyzerService,
		Timeline: timelineService,
	}
}

func (ha *HTTPApplication) Run() error {
	ha.Logger.Info("Starting server...", "port", ha.Config.Port)
	return ha.Server.Run()
}

func (ha *HTTPApplication) Shutdown(ctx context.Context) error {
	ha.Logger.Info("Shutting down server...")
	return ha.Server.Shutdown(ctx)
}
