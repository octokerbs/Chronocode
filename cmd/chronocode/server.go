package main

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/octokerbs/chronocode-backend/internal/api"
	"github.com/octokerbs/chronocode-backend/internal/application"
	"github.com/octokerbs/chronocode-backend/internal/infrastructure/agent/gemini"
	"github.com/octokerbs/chronocode-backend/internal/infrastructure/codehost/githubapi"
	"github.com/octokerbs/chronocode-backend/internal/infrastructure/database/mock"
)

func main() {

	// Load env variables
	_ = godotenv.Load()

	// Create the execution context
	ctx := context.Background()

	// Setup dependencies
	geminiClient, err := gemini.NewGeminiClient(ctx, os.Getenv("GEMINI_API_KEY"))
	if err != nil {
		panic(err)
	}

	githubFactory := githubapi.NewGithubFactory()

	// pgClient, err := postgres.NewPostgresClient(ctx, os.Getenv("POSTGRES_DSN"))
	// if err != nil {
	// 	panic(err)
	// }

	mockPostgresClient := &mock.PostgresMock{}

	// Create the application entities
	repoAnalyzer := application.NewRepositoryAnalyzer(ctx, geminiClient, githubFactory, mockPostgresClient)

	// Launch server
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
