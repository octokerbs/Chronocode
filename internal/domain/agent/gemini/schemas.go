package gemini

import "github.com/google/generative-ai-go/genai"

var CommitAnalysisSchema = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"commit":     commitSchema,
		"subcommits": subcommitsSchema,
	},
	Required: []string{"commit", "subcommits"},
}

var commitSchema = &genai.Schema{
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

var subcommitsSchema = &genai.Schema{
	Type:        genai.TypeArray,
	Items:       subcommitSchema,
	Description: "An array of logical units of work that make up this commit.",
}

var subcommitSchema = &genai.Schema{
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
