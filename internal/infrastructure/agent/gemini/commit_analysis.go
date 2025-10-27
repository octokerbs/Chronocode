package gemini

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"github.com/octokerbs/chronocode-backend/internal/domain"
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
	return &GeminiClient{client: client, model: model}, nil
}

func (gc *GeminiClient) AnalyzeDiff(ctx context.Context, diff string) (domain.CommitAnalysis, error) {
	tries := 3
	var text []byte
	var err error
	for tries > 0 {
		text, err = gc.generateCommitAnalysis(ctx, gc.commitAnalysisPrompt()+diff)
		if err != nil {
			tries--
			continue
		}
		break
	}

	if tries == 0 {
		return domain.CommitAnalysis{}, errors.New("no text response parts found for commit")
	}

	return gc.unmarshalCommitAnalysisSchemaOntoStruct(text)
}

func (gc *GeminiClient) generateCommitAnalysis(ctx context.Context, prompt string) ([]byte, error) {
	gc.model.ResponseSchema = gc.commitAnalysisSchema()

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

func (gc *GeminiClient) commitAnalysisSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"commit":     gc.commitSchema(),
			"subcommits": gc.subcommitsSchema(),
		},
		Required: []string{"commit", "subcommits"},
	}
}

func (gc *GeminiClient) commitSchema() *genai.Schema {
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

func (gc *GeminiClient) subcommitsSchema() *genai.Schema {
	return &genai.Schema{
		Type:        genai.TypeArray,
		Items:       gc.subcommitSchema(),
		Description: "An array of logical units of work that make up this commit.",
	}
}

func (gc *GeminiClient) subcommitSchema() *genai.Schema {
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

func (gc *GeminiClient) commitAnalysisPrompt() string {
	return `
	You are a Commit Expert Analyzer specializing in code analysis and software development patterns.
	You will receive a Git Commit diff.
	Your task is to given commit, identify the logical units of work ("SubCommits") within this single GitHub commit. 
	The subcommits will have a title, idea, description, and type.

	Now extract the subcommits from the following diff:
	`
}

func (gc *GeminiClient) unmarshalCommitAnalysisSchemaOntoStruct(text []byte) (domain.CommitAnalysis, error) {
	analysis := &domain.CommitAnalysis{
		Commit:     domain.Commit{},
		Subcommits: []domain.Subcommit{},
	}

	if err := json.Unmarshal(text, &analysis); err != nil {
		var subcommits []domain.Subcommit
		if err := json.Unmarshal(text, &subcommits); err != nil {
			var subcommit domain.Subcommit // Try a single subcommit if array fails
			err := json.Unmarshal(text, &subcommit)
			if err != nil {
				return domain.CommitAnalysis{}, err
			}
			analysis.Subcommits = []domain.Subcommit{subcommit}
		} else {
			analysis.Subcommits = subcommits
		}
	}

	return *analysis, nil
}

func (gc *GeminiClient) Close() {
	gc.client.Close()
}
