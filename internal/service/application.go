package service

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	_ "github.com/lib/pq"
	"github.com/octokerbs/chronocode-backend/internal/adapters"
	"github.com/octokerbs/chronocode-backend/internal/app"
	"google.golang.org/api/option"
)

func NewApplication(ctx context.Context) app.Application {
	geminiClient, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		panic(err)
	}

	agent, err := adapters.NewGeminiAgent(geminiClient, os.Getenv("GEMINI_GENERATIVE_MODEL"))
	if err != nil {
		panic(err)
	}

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

	codeHostFactory := adapters.NewGithubCodeHostFactory()

	fmt.Println(agent)
	fmt.Println(repositoryRepository)
	fmt.Println(codeHostFactory)

	return app.Application{
		Commands: app.Commands{},
		Queries:  app.Queries{},
	}
}
