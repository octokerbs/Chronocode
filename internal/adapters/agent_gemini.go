package adapters

import (
	"errors"

	"github.com/google/generative-ai-go/genai"
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

	return &GeminiAgent{client: client, generativeModel: generativeModel}, nil
}

// func (ga *GeminiAgent) AnalyzeCommitDiff(ctx context.Context, diff string) (analysis.CommitAnalysis, error) {
// 	tries := 3
// 	prompt := ga.commitAnalysisPrompt() + diff

// 	var text []byte
// 	var err error
// 	for tries > 0 {
// 		text, err = ga.generateStructuredContent(ctx, prompt, ga.commitAnalysisSchema())
// 		if err == nil {
// 			break
// 		}
// 		tries--
// 	}

// 	if err != nil {
// 		return analysis.CommitAnalysis{}, err
// 	}

// 	var commitAnalysis analysis.CommitAnalysis
// 	if err := json.Unmarshal(text, &commitAnalysis); err != nil {
// 		return analysis.CommitAnalysis{}, err
// 	}

// 	return commitAnalysis, nil
// }

// func (ga *GeminiAgent) commitAnalysisSchema() *genai.Schema {
// 	return &genai.Schema{
// 		Type: genai.TypeObject,
// 		Properties: map[string]*genai.Schema{
// 			"commit":     ga.commitSchema(),
// 			"subcommits": ga.subcommitsSchema(),
// 		},
// 		Required: []string{"commit", "subcommits"},
// 	}
// }

// func (ga *GeminiAgent) commitSchema() *genai.Schema {
// 	return &genai.Schema{
// 		Type: genai.TypeObject,
// 		Properties: map[string]*genai.Schema{
// 			"description": {
// 				Type:        genai.TypeString,
// 				Description: "Brief summary of the entire diff, explaining its overall purpose and changes. Don't talk about a commit. Just the diff",
// 			},
// 		},
// 		Required:    []string{"description"},
// 		Description: "Information about the overall commit",
// 	}
// }

// func (ga *GeminiAgent) subcommitsSchema() *genai.Schema {
// 	return &genai.Schema{
// 		Type:        genai.TypeArray,
// 		Items:       ga.subcommitSchema(),
// 		Description: "An array of logical units of work that make up this commit.",
// 	}
// }

// func (ga *GeminiAgent) subcommitSchema() *genai.Schema {
// 	return &genai.Schema{
// 		Type: genai.TypeObject,
// 		Properties: map[string]*genai.Schema{
// 			"title": {
// 				Type:        genai.TypeString,
// 				Description: "A concise, specific title (5-10 words) that precisely captures what this logical unit of work accomplishes.",
// 			},
// 			"idea": {
// 				Type:        genai.TypeString,
// 				Description: "The core concept or purpose (max 15 sentences) explaining why this change was made and what problem it solves. MAKE IT SHORT.",
// 			},
// 			"description": {
// 				Type:        genai.TypeString,
// 				Description: "A comprehensive technical explanation detailing implementation specifics, architectural changes, and potential downstream effects.",
// 			},
// 			"type": {
// 				Type:        genai.TypeString,
// 				Description: "The primary category that best represents the nature of this change.",
// 				Enum:        []string{"FEATURE", "BUG", "REFACTOR", "DOCS", "CHORE", "MILESTONE", "WARNING"},
// 			},
// 			"epic": {
// 				Type:        genai.TypeString,
// 				Description: "If the subcommit is part of a larger epic or feature, mention the epic's name or identifier. If not, leave it blank.",
// 			},
// 			"files": {
// 				Type: genai.TypeArray,
// 				Items: &genai.Schema{
// 					Type: genai.TypeString,
// 				},
// 				Description: "An array of file names that are directly related to this subcommit",
// 			},
// 		},
// 		Required: []string{"title", "idea", "description", "type", "epic", "files"},
// 	}
// }

// func (ga *GeminiAgent) commitAnalysisPrompt() string {
// 	return `
// 	You are a Commit Expert Analyzer specializing in code analysis and software development patterns.
// 	You will receive a Git Commit diff.
// 	Your task is to given commit, identify the logical units of work ("SubCommits") within this single GitHub commit.
// 	The subcommits will have a title, idea, description, and type.

// 	Now extract the subcommits from the following diff:
// 	`
// }

// func (ga *GeminiAgent) generateStructuredContent(ctx context.Context, prompt string, schema *genai.Schema) ([]byte, error) {
// 	ga.generativeModel.ResponseSchema = schema

// 	resp, err := ga.generativeModel.GenerateContent(ctx, genai.Text(prompt))
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, part := range resp.Candidates[0].Content.Parts {
// 		if text, ok := part.(genai.Text); ok {
// 			return []byte(text), nil
// 		}
// 	}

// 	return nil, err
// }
