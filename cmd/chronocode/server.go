package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/octokerbs/chronocode-backend/internal/api"
	"github.com/octokerbs/chronocode-backend/internal/application"
	"github.com/octokerbs/chronocode-backend/internal/infrastructure/agent/gemini"
	"github.com/octokerbs/chronocode-backend/internal/infrastructure/codehost/githubapi"
	"github.com/octokerbs/chronocode-backend/internal/infrastructure/database/postgres"
)

func main() {
	_ = godotenv.Load()
	ctx := context.Background()

	geminiClient, err := gemini.NewGeminiAgent(ctx, os.Getenv("GEMINI_API_KEY"))
	if err != nil {
		panic(err)
	}

	githubClient := githubapi.NewGitHubFactory()

	dsn := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	pgClient, err := postgres.NewPostgresDatabase(dsn)
	if err != nil {
		panic(err)
	}

	repoAnalyzer := application.NewRepositoryAnalyzer(ctx, geminiClient, githubClient, pgClient)

	server, err := NewServer(":8080", repoAnalyzer)
	if err != nil {
		panic(err)
	}

	server.Run()
}

type Server struct {
	anEngine *gin.Engine
	aPort    string
}

func NewServer(port string, repoAnalyzer *application.RepositoryAnalyzer) (*Server, error) {
	engine := gin.Default()

	repoAnalyzerHandler := api.NewAnalysisHandler(repoAnalyzer)

	engine.GET("/analyze-repository", repoAnalyzerHandler.AnalyzeRepository)

	return &Server{engine, port}, nil
}

func (s *Server) Run() {
	s.anEngine.Run(s.aPort)
}
