package gemini

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiGenerativeService struct {
	geminiService *genai.Client
	model         *genai.GenerativeModel
}

func NewGeminiGenerativeService(ctx context.Context, key string) (*GeminiGenerativeService, error) {
	geminiService, err := genai.NewClient(ctx, option.WithAPIKey(key))
	if err != nil {
		return nil, err
	}
	model := geminiService.GenerativeModel("gemini-2.0-flash")
	model.ResponseMIMEType = "application/json"
	model.ResponseSchema = commitSchema
	return &GeminiGenerativeService{geminiService: geminiService, model: model}, nil
}

func (g *GeminiGenerativeService) Generate(ctx context.Context, prompt string) ([]byte, error) {
	resp, err := g.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("error generating content %v", err.Error())
	}

	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			return []byte(text), nil
		}
	}

	return nil, fmt.Errorf("no text response parts found for commit")
}

func (g *GeminiGenerativeService) Close() {
	g.geminiService.Close()
}
