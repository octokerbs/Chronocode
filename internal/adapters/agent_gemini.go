package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/google/generative-ai-go/genai"
	"github.com/octokerbs/chronocode/internal/domain/agent"
)

type GeminiAgent struct {
	client          *genai.Client
	generativeModel *genai.GenerativeModel
}

func NewGeminiAgent(client *genai.Client, model string) (*GeminiAgent, error) {
	if client == nil {
		return nil, errors.New("missing gemini client")
	}

	if model == "" {
		return nil, errors.New("missing model")
	}

	generativeModel := client.GenerativeModel(model)
	generativeModel.ResponseMIMEType = "application/json"

	slog.Info("Gemini agent initialized", "model", model)
	return &GeminiAgent{client: client, generativeModel: generativeModel}, nil
}

type subcommitResponse struct {
	Title            string   `json:"title"`
	Idea             string   `json:"idea"`
	Description      string   `json:"description"`
	Epic             string   `json:"epic"`
	ModificationType string   `json:"type"`
	Files            []string `json:"files"`
}

type analysisResponse struct {
	Subcommits []subcommitResponse `json:"subcommits"`
}

func (ga *GeminiAgent) AnalyzeDiff(ctx context.Context, diff string) ([]agent.AnalysisResult, error) {
	slog.Debug("Gemini analyzing diff", "diff_length", len(diff))

	prompt := ga.commitAnalysisPrompt() + diff

	text, err := ga.generateStructuredContent(ctx, prompt, ga.analysisSchema())
	if err != nil {
		slog.Error("Gemini content generation failed", "error", err, "diff_length", len(diff))
		return nil, err
	}

	var response analysisResponse
	if err := json.Unmarshal(text, &response); err != nil {
		slog.Error("Failed to unmarshal Gemini response", "error", err, "response_length", len(text))
		return nil, err
	}

	results := make([]agent.AnalysisResult, len(response.Subcommits))
	for i, sc := range response.Subcommits {
		results[i] = agent.AnalysisResult{
			Title:            sc.Title,
			Idea:             sc.Idea,
			Description:      sc.Description,
			Epic:             sc.Epic,
			ModificationType: sc.ModificationType,
			Files:            sc.Files,
		}
	}

	slog.Debug("Gemini analysis completed", "subcommits_produced", len(results), "diff_length", len(diff))
	return results, nil
}

func (ga *GeminiAgent) analysisSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"subcommits": {
				Type:        genai.TypeArray,
				Items:       ga.subcommitSchema(),
				Description: "An array of logical units of work that make up this commit.",
			},
		},
		Required: []string{"subcommits"},
	}
}

func (ga *GeminiAgent) subcommitSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"title": {
				Type:        genai.TypeString,
				Description: "A concise, specific title (5-10 words) that precisely captures what this logical unit of work accomplishes.",
			},
			"idea": {
				Type:        genai.TypeString,
				Description: "A one-sentence thesis explaining the core motivation or reasoning behind this change.",
			},
			"description": {
				Type:        genai.TypeString,
				Description: "A technical explanation detailing implementation specifics and what problem it solves.",
			},
			"epic": {
				Type:        genai.TypeString,
				Description: "A broad initiative or project area label this change belongs to (e.g. 'Authentication', 'Performance', 'CI/CD').",
			},
			"type": {
				Type:        genai.TypeString,
				Description: "The primary category that best represents the nature of this change.",
				Enum:        []string{"FEATURE", "BUG", "REFACTOR", "DOCS", "CHORE", "MILESTONE", "WARNING"},
			},
			"files": {
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeString,
				},
				Description: "An array of file names that are directly related to this subcommit.",
			},
		},
		Required: []string{"title", "idea", "description", "epic", "type", "files"},
	}
}

func (ga *GeminiAgent) commitAnalysisPrompt() string {
	return `You are a Commit Expert Analyzer specializing in code analysis and software development patterns.
You will receive a Git Commit diff.
Your task is to identify the logical units of work ("SubCommits") within this single commit.
Each subcommit should have:
- title: A concise title (5-10 words)
- idea: A one-sentence thesis explaining the core motivation behind this change
- description: A technical explanation of the implementation
- epic: A broad initiative label (e.g. "Authentication", "Performance", "CI/CD")
- type: One of FEATURE, BUG, REFACTOR, DOCS, CHORE, MILESTONE, WARNING
- files: List of related file names

Now extract the subcommits from the following diff:
`
}

func (ga *GeminiAgent) generateStructuredContent(ctx context.Context, prompt string, schema *genai.Schema) ([]byte, error) {
	ga.generativeModel.ResponseSchema = schema

	slog.Debug("Sending request to Gemini API", "prompt_length", len(prompt))

	resp, err := ga.generativeModel.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		slog.Error("Gemini API request failed", "error", err)
		return nil, err
	}

	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			slog.Debug("Gemini API response received", "response_length", len(text))
			return []byte(text), nil
		}
	}

	slog.Error("Gemini API returned no text content in response")
	return nil, errors.New("no text content in response")
}
