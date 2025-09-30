package domain

import (
	"encoding/json"
)

type CommitAnalysis struct {
	Commit     Commit      `json:"commit"`
	Subcommits []Subcommit `json:"subcommits"`
}

type Commit struct {
	Description string `json:"description"`
}

type Subcommit struct {
	Title       string   `json:"title"`
	Idea        string   `json:"idea"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Epic        string   `json:"epic"`
	Files       []string `json:"files"`
}

func UnmarshalCommitAnalysisSchemaOntoStruct(text []byte) (CommitAnalysis, error) {
	analysis := &CommitAnalysis{
		Commit:     Commit{},
		Subcommits: []Subcommit{},
	}

	if err := json.Unmarshal(text, &analysis); err != nil {
		var subcommits []Subcommit
		if err := json.Unmarshal(text, &subcommits); err != nil {
			var subcommit Subcommit // Try a single subcommit if array fails
			err := json.Unmarshal(text, &subcommit)
			if err != nil {
				return CommitAnalysis{}, err
			}
			analysis.Subcommits = []Subcommit{subcommit}
		} else {
			analysis.Subcommits = subcommits
		}
	}

	return *analysis, nil
}
