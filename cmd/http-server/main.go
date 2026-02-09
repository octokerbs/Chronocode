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
	"github.com/octokerbs/chronocode-backend/internal/domain/cache"
	"github.com/octokerbs/chronocode-backend/internal/domain/codehost"
	"github.com/octokerbs/chronocode-backend/internal/infrastructure/agent/gemini"
	"github.com/octokerbs/chronocode-backend/internal/infrastructure/auth/githubauth"
	rediscache "github.com/octokerbs/chronocode-backend/internal/infrastructure/cache/redis"
	"github.com/octokerbs/chronocode-backend/internal/infrastructure/codehost/githubapi"
	"github.com/octokerbs/chronocode-backend/internal/infrastructure/database/postgres"
	"go.uber.org/zap"
)

type dependencies struct {
	analyzer        *application.Analyzer
	persistCommits  *application.PersistCommits
	prepareRepo     *application.PrepareRepository
	querier         *application.Querier
	auth            *application.Auth
	userProfile     *application.UserProfile
	cache           cache.Cache
	codeHostFactory codehost.CodeHostFactory
}

func main() {
	ctx := context.Background()
	_ = godotenv.Load()
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	deps := buildDependencies(ctx, cfg)
	server := http.NewHTTPServer(deps.analyzer, deps.persistCommits, deps.prepareRepo, deps.querier, deps.auth, deps.userProfile, deps.cache, deps.codeHostFactory, ":8080", cfg.FrontendURL)

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

func buildDependencies(ctx context.Context, cfg *config.Config) *dependencies {
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
		logger.Fatal("Couldn't initialize gemini agent", zap.String("buildDependencies", err.Error()))
	}

	rc, err := rediscache.NewRedisCache(cfg.RedisURL)
	if err != nil {
		logger.Fatal("Couldn't initialize Redis cache", zap.String("buildDependencies", err.Error()))
	}

	githubAuth := githubauth.NewGitHubAuth(
		cfg.GithubClientID,
		cfg.GithubClientSecret,
		cfg.RedirectURL,
	)

	auth := application.NewAuth(githubAuth)

	codeHostFactory := githubapi.NewGitHubCodeHostFactory()

	querier := application.NewQuerier(db)

	analyzer := application.NewAnalyzer(agent, codeHostFactory)

	persistCommits := &application.PersistCommits{Database: db}
	prepareRepo := &application.PrepareRepository{CodeHostFactory: codeHostFactory, Database: db}

	userProfile := &application.UserProfile{
		CodeHostFactory: codeHostFactory,
		Cache:           rc,
	}

	return &dependencies{
		analyzer:        analyzer,
		persistCommits:  persistCommits,
		prepareRepo:     prepareRepo,
		querier:         querier,
		auth:            auth,
		userProfile:     userProfile,
		cache:           rc,
		codeHostFactory: codeHostFactory,
	}
}
