package gemini

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiClient struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

func NewGeminiClient(ctx context.Context, key string) (*GeminiClient, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(key))
	if err != nil {
		return nil, err
	}
	model := client.GenerativeModel("gemini-2.0-flash")
	model.ResponseMIMEType = "application/json"
	model.ResponseSchema = CommitAnalysisSchema
	return &GeminiClient{client: client, model: model}, nil
}

func (gc *GeminiClient) Generate(ctx context.Context, prompt string) ([]byte, error) {
	resp, err := gc.model.GenerateContent(ctx, genai.Text(prompt))
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

func (gc *GeminiClient) Close() {
	gc.client.Close()
}
