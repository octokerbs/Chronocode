package main

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/octokerbs/chronocode-go/internal/adapters/agent/gemini"
	githubapi "github.com/octokerbs/chronocode-go/internal/adapters/codehost/github"
	"github.com/octokerbs/chronocode-go/internal/adapters/database/postgres"
	"github.com/octokerbs/chronocode-go/internal/api"
	"github.com/octokerbs/chronocode-go/internal/application"
)

type Server struct {
	anEngine *gin.Engine
	aPort    string
}

func NewServer(port string) (*Server, error) {
	_ = godotenv.Load()

	engine := gin.Default()

	ctx := context.Background()

	geminiClient, err := gemini.NewGeminiClient(ctx, os.Getenv("GEMINI_API_KEY"))
	if err != nil {
		return nil, err
	}

	githubFactory := githubapi.NewGithubFactory()

	pgClient, err := postgres.NewPostgresClient(ctx, os.Getenv("POSTGRES_DSN"))
	if err != nil {
		return nil, err
	}

	repoAnalyzer := application.NewRepositoryAnalyzer(ctx, geminiClient, githubFactory, pgClient)

	handler := api.NewAnalysisHandler(repoAnalyzer)

	engine.GET("/analyze-repository", handler.AnalyzeRepository)

	return &Server{engine, port}, nil
}

func (s *Server) Run() {
	s.anEngine.Run(s.aPort)
}
