package service

import (
	"context"
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

	agent, err := adapters.NewAgentGemini(geminiClient, os.Getenv("GEMINI_GENERATIVE_MODEL"))
	if err != nil {
		panic(err)
	}

	fmt.Println(agent)

	return app.Application{
		Commands: app.Commands{},
		Queries:  app.Queries{},
	}
}
