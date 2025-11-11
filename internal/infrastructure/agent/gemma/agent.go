package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/octokerbs/chronocode-backend/internal/domain/analysis"
	pkg_errors "github.com/octokerbs/chronocode-backend/pkg/errors"
)

type LLMRequest struct {
	Model       string       `json:"model"`
	Messages    []LLMMessage `json:"messages"`
	Temperature float64      `json:"temperature"`
}

type LLMMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type LLMResponse struct {
	Choices []struct {
		Message LLMMessage `json:"message"`
	} `json:"choices"`
}

type GemmaAgent struct {
	httpClient  *http.Client
	endpointURL string
	modelID     string
}

func NewGemmaAgent(ctx context.Context, endpointURL string) (*GemmaAgent, error) {
	if endpointURL == "" {
		endpointURL = "http://localhost:8080/v1/chat/completions"
	}

	return &GemmaAgent{
		httpClient:  &http.Client{Timeout: 60 * time.Second},
		endpointURL: endpointURL,
		modelID:     "gemma3:4B-Q4_K_M",
	}, nil
}

func (ga *GemmaAgent) AnalyzeCommitDiff(ctx context.Context, diff string) (analysis.CommitAnalysis, error) {
	tries := 3

	systemInstruction := ga.systemInstructionForGemma() + "\n\n" + ga.getSchemaDescription()
	userPrompt := ga.commitAnalysisPrompt() + diff

	var text []byte
	var err error
	for tries > 0 {
		text, err = ga.generateContent(ctx, systemInstruction, userPrompt)
		if err == nil {
			break
		}
		tries--
	}

	if err != nil {
		return analysis.CommitAnalysis{}, pkg_errors.NewError(pkg_errors.ErrInternalFailure, fmt.Errorf("gemma generation failed after %d tries: %w", 3, err))
	}

	var commitAnalysis analysis.CommitAnalysis
	if err := json.Unmarshal(text, &commitAnalysis); err != nil {
		return analysis.CommitAnalysis{}, pkg_errors.NewError(pkg_errors.ErrInternalFailure, fmt.Errorf("failed to unmarshal JSON response from Gemma: %w. Raw output: %s", err, string(text)))
	}

	return commitAnalysis, nil
}

func (ga *GemmaAgent) systemInstructionForGemma() string {
	return `
    You are a Commit Expert Analyzer specializing in code analysis and software development patterns.
    You will receive a Git Commit diff. Your task is to identify the logical units of work ("SubCommits") within this single commit.
    
    You MUST respond with a single valid JSON object that strictly conforms to the provided JSON schema description.
    DO NOT include any explanation, introductory text, or Markdown backticks (e.g., ` + "```json" + `) outside the JSON object.
    The response must start and end with the JSON curly braces {} exactly.
    `
}

func (ga *GemmaAgent) getSchemaDescription() string {
	return `
    JSON Schema to follow:
    {
      "commit": {
        "description": "Brief summary of the entire diff, explaining its overall purpose and changes. Don't talk about a commit. Just the diff"
      },
      "subcommits": [
        {
          "title": "A concise, specific title (5-10 words) that precisely captures what this logical unit of work accomplishes.",
          "idea": "The core concept or purpose (max 15 sentences) explaining why this change was made and what problem it solves. MAKE IT SHORT.",
          "description": "A comprehensive technical explanation detailing implementation specifics, architectural changes, and potential downstream effects.",
          "type": "Must be one of: FEATURE, BUG, REFACTOR, DOCS, CHORE, MILESTONE, WARNING",
          "epic": "If part of a larger epic, mention its name or identifier. If not, leave it blank (empty string).",
          "files": ["An array of file names that are directly related to this subcommit"]
        }
      ]
    }
    `
}

func (ga *GemmaAgent) commitAnalysisPrompt() string {
	return `
	Now extract the subcommits from the following diff:
	`
}

func (ga *GemmaAgent) generateContent(ctx context.Context, systemInstruction, userPrompt string) ([]byte, error) {
	requestBody := LLMRequest{
		Model: ga.modelID,
		Messages: []LLMMessage{
			{Role: "system", Content: systemInstruction},
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.1, // Baja temperatura para análisis estructurado
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", ga.endpointURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := ga.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request to local LLM: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("local LLM returned error status %d: %s", resp.StatusCode, string(responseBody))
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var llmResponse LLMResponse
	if err := json.Unmarshal(responseBody, &llmResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal LLM response structure: %w. Raw API response: %s", err, string(responseBody))
	}

	if len(llmResponse.Choices) == 0 || llmResponse.Choices[0].Message.Content == "" {
		return nil, errors.New("LLM returned empty content")
	}

	return []byte(llmResponse.Choices[0].Message.Content), nil
}
