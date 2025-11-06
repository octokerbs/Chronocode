package gemini

import (
	"context"
	"encoding/json"

	"github.com/google/generative-ai-go/genai"
	"github.com/octokerbs/chronocode-backend/internal/domain"
	"github.com/octokerbs/chronocode-backend/internal/domain/analysis"
	"google.golang.org/api/option"
)

type Agent struct {
	gemini *geminiClient
}

func NewGeminiAgent(ctx context.Context, apiKey string) (*Agent, error) {
	gemini, err := newgeminiClient(ctx, apiKey)
	if err != nil {
		return nil, err
	}

	return &Agent{gemini}, nil
}

func (a *Agent) AnalyzeCommitDiff(ctx context.Context, diff string) (analysis.CommitAnalysis, error) {
	tries := 3
	prompt := a.gemini.commitAnalysisPrompt() + diff

	var text []byte
	var err error
	for tries > 0 {
		text, err = a.gemini.generateStructuredContent(ctx, prompt, a.gemini.commitAnalysisSchema())
		if err == nil {
			break
		}
		tries--
	}

	if err != nil {
		return analysis.CommitAnalysis{}, domain.NewError(domain.ErrInternalFailure, err)
	}

	var commitAnalysis analysis.CommitAnalysis
	if err := json.Unmarshal(text, &commitAnalysis); err != nil {
		return analysis.CommitAnalysis{}, domain.NewError(domain.ErrInternalFailure, err)
	}

	return commitAnalysis, nil
}

type geminiClient struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

func newgeminiClient(ctx context.Context, key string) (*geminiClient, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(key))
	if err != nil {
		return nil, err
	}
	model := client.GenerativeModel("gemini-2.0-flash")
	model.ResponseMIMEType = "application/json"
	return &geminiClient{client: client, model: model}, nil
}

func (gc *geminiClient) commitAnalysisSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"commit":     gc.commitSchema(),
			"subcommits": gc.subcommitsSchema(),
		},
		Required: []string{"commit", "subcommits"},
	}
}

func (gc *geminiClient) commitSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"description": {
				Type:        genai.TypeString,
				Description: "Brief summary of the entire diff, explaining its overall purpose and changes. Don't talk about a commit. Just the diff",
			},
		},
		Required:    []string{"description"},
		Description: "Information about the overall commit",
	}
}

func (gc *geminiClient) subcommitsSchema() *genai.Schema {
	return &genai.Schema{
		Type:        genai.TypeArray,
		Items:       gc.subcommitSchema(),
		Description: "An array of logical units of work that make up this commit.",
	}
}

func (gc *geminiClient) subcommitSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"title": {
				Type:        genai.TypeString,
				Description: "A concise, specific title (5-10 words) that precisely captures what this logical unit of work accomplishes.",
			},
			"idea": {
				Type:        genai.TypeString,
				Description: "The core concept or purpose (max 15 sentences) explaining why this change was made and what problem it solves. MAKE IT SHORT.",
			},
			"description": {
				Type:        genai.TypeString,
				Description: "A comprehensive technical explanation detailing implementation specifics, architectural changes, and potential downstream effects.",
			},
			"type": {
				Type:        genai.TypeString,
				Description: "The primary category that best represents the nature of this change.",
				Enum:        []string{"FEATURE", "BUG", "REFACTOR", "DOCS", "CHORE", "MILESTONE", "WARNING"},
			},
			"epic": {
				Type:        genai.TypeString,
				Description: "If the subcommit is part of a larger epic or feature, mention the epic's name or identifier. If not, leave it blank.",
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

func (gc *geminiClient) commitAnalysisPrompt() string {
	return `
	You are a Commit Expert Analyzer specializing in code analysis and software development patterns.
	You will receive a Git Commit diff.
	Your task is to given commit, identify the logical units of work ("SubCommits") within this single GitHub commit. 
	The subcommits will have a title, idea, description, and type.

	Now extract the subcommits from the following diff:
	`
}

func (gc *geminiClient) generateStructuredContent(ctx context.Context, prompt string, schema *genai.Schema) ([]byte, error) {
	gc.model.ResponseSchema = schema

	resp, err := gc.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, err
	}

	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			return []byte(text), nil
		}
	}

	return nil, err
}
