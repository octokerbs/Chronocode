package adapters

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/generative-ai-go/genai"
	"github.com/octokerbs/chronocode-backend/internal/domain/repository"
	"google.golang.org/api/option"
)

type GeminiAgent struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

func NewGeminiAgent(ctx context.Context, key string) (*GeminiAgent, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(key))
	if err != nil {
		return nil, err
	}
	model := client.GenerativeModel("gemini-2.0-flash")
	model.ResponseMIMEType = "application/json"
	return &GeminiAgent{client: client, model: model}, nil
}

func (ga *GeminiAgent) AnalyzeCommitDiff(ctx context.Context, diff string) (repository.CommitAnalysis, error) {
	tries := 3
	prompt := ga.commitAnalysisPrompt() + diff

	var text []byte
	var err error
	for tries > 0 {
		text, err = ga.generateStructuredContent(ctx, prompt, ga.commitAnalysisSchema())
		if err == nil {
			break
		}
		tries--
	}

	if err != nil {
		return repository.CommitAnalysis{}, err
	}

	var response struct {
		Commit struct {
			Description string `json:"description"`
		} `json:"commit"`
		Subcommits []struct {
			Title       string   `json:"title"`
			Idea        string   `json:"idea"`
			Description string   `json:"description"`
			Type        string   `json:"type"`
			Epic        string   `json:"epic"`
			Files       []string `json:"files"`
		} `json:"subcommits"`
	}

	if err := json.Unmarshal(text, &response); err != nil {
		return repository.CommitAnalysis{}, err
	}

	analysis := repository.CommitAnalysis{
		Description: response.Commit.Description,
	}

	for _, sc := range response.Subcommits {
		analysis.Subcommits = append(analysis.Subcommits, repository.Subcommit{
			Title:       sc.Title,
			Idea:        sc.Idea,
			Description: sc.Description,
			Type:        sc.Type,
			Epic:        sc.Epic,
			Files:       sc.Files,
		})
	}

	return analysis, nil
}

func (ga *GeminiAgent) commitAnalysisSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"commit":     ga.commitSchema(),
			"subcommits": ga.subcommitsSchema(),
		},
		Required: []string{"commit", "subcommits"},
	}
}

func (ga *GeminiAgent) commitSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"description": {
				Type:        genai.TypeString,
				Description: "Brief summary of the entire diff, explaining its overall purpose and changes.",
			},
		},
		Required:    []string{"description"},
		Description: "Information about the overall commit",
	}
}

func (ga *GeminiAgent) subcommitsSchema() *genai.Schema {
	return &genai.Schema{
		Type:        genai.TypeArray,
		Items:       ga.subcommitSchema(),
		Description: "An array of logical units of work that make up this commit.",
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
				Description: "The core concept or purpose (max 15 sentences) explaining why this change was made.",
			},
			"description": {
				Type:        genai.TypeString,
				Description: "A comprehensive technical explanation detailing implementation specifics.",
			},
			"type": {
				Type:        genai.TypeString,
				Description: "The primary category that best represents the nature of this change.",
				Enum:        []string{"FEATURE", "BUG", "REFACTOR", "DOCS", "CHORE", "MILESTONE", "WARNING"},
			},
			"epic": {
				Type:        genai.TypeString,
				Description: "If the subcommit is part of a larger epic, mention the epic's name. If not, leave it blank.",
			},
			"files": {
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeString,
				},
				Description: "An array of file names that are directly related to this subcommit",
			},
		},
		Required: []string{"title", "idea", "description", "type", "epic", "files"},
	}
}

func (ga *GeminiAgent) commitAnalysisPrompt() string {
	return `
	You are a Commit Expert Analyzer specializing in code analysis and software development patterns.
	You will receive a Git Commit diff.
	Your task is to identify the logical units of work ("SubCommits") within this single GitHub commit.
	The subcommits will have a title, idea, description, and type.

	Now extract the subcommits from the following diff:
	`
}

func (ga *GeminiAgent) generateStructuredContent(ctx context.Context, prompt string, schema *genai.Schema) ([]byte, error) {
	ga.model.ResponseSchema = schema

	resp, err := ga.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, err
	}

	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			return []byte(text), nil
		}
	}

	return nil, errors.New("no text content in response")
}
