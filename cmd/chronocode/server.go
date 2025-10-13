package main

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/octokerbs/chronocode-backend/internal/handler"
	"github.com/octokerbs/chronocode-backend/internal/infrastructure/agent/gemini"
	"github.com/octokerbs/chronocode-backend/internal/infrastructure/codehost/githubapi"
	"github.com/octokerbs/chronocode-backend/internal/infrastructure/database/postgres"
	"github.com/octokerbs/chronocode-backend/internal/usecase"
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

	repoAnalyzer := usecase.NewRepositoryAnalyzer(ctx, geminiClient, githubFactory, pgClient)

	handler := handler.NewAnalysisHandler(repoAnalyzer)

	engine.GET("/analyze-repository", handler.AnalyzeRepository)

	return &Server{engine, port}, nil
}

func (s *Server) Run() {
	s.anEngine.Run(s.aPort)
}
