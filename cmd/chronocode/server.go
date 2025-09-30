package main

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/octokerbs/chronocode-go/internal/api"
	"github.com/octokerbs/chronocode-go/internal/domain/agent/gemini"
	"github.com/octokerbs/chronocode-go/internal/repository/postgres"
	"github.com/octokerbs/chronocode-go/internal/service"
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

	pgClient, err := postgres.NewPostgresClient(ctx, os.Getenv("POSTGRES_DSN"))
	if err != nil {
		return nil, err
	}

	repoAnalyzer := service.NewRepositoryAnalyzer(ctx, geminiClient, pgClient)

	handler := api.NewAnalysisHandler(repoAnalyzer)

	engine.GET("/analyze-repository", handler.AnalyzeRepository)

	return &Server{engine, port}, nil
}

func (s *Server) Run() {
	s.anEngine.Run(s.aPort)
}
