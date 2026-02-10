package service

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/octokerbs/chronocode-backend/internal2/adapters"
	"github.com/octokerbs/chronocode-backend/internal2/app"
	"google.golang.org/api/option"
)

func NewApplication(ctx context.Context) app.Application {

	// Agent setup
	geminiClient, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		panic(err)
	}

	agent, err := adapters.NewGeminiAgent(geminiClient, os.Getenv("GEMINI_GENERATIVE_MODEL"))
	if err != nil {
		panic(err)
	}

	// Repositories setup
	postgresClient, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	if err := postgresClient.Ping(); err != nil {
		panic(err)
	}

	repositoryRepository, err := adapters.NewPostgresRepositoryRepository(postgresClient)
	if err != nil {
		panic(err)
	}

	fmt.Println(agent)
	fmt.Println(repositoryRepository)

	return app.Application{
		Commands: app.Commands{},
		Queries:  app.Queries{},
	}
}
