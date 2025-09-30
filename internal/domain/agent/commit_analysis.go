package agent

import (
	"encoding/json"
)

type CommitAnalysisSchema struct {
	Commit     CommitSchema      `json:"commit"`
	Subcommits []SubcommitSchema `json:"subcommits"`
}

type CommitSchema struct {
	Description string `json:"description"`
}

type SubcommitSchema struct {
	Title       string   `json:"title"`
	Idea        string   `json:"idea"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Epic        string   `json:"epic"`
	Files       []string `json:"files"`
}

func UnmarshalCommitAnalysisSchemaOntoStruct(text []byte) (CommitAnalysisSchema, error) {
	analysis := &CommitAnalysisSchema{
		Commit:     CommitSchema{},
		Subcommits: []SubcommitSchema{},
	}

	if err := json.Unmarshal(text, &analysis); err != nil {
		var subcommits []SubcommitSchema
		if err := json.Unmarshal(text, &subcommits); err != nil {
			var subcommit SubcommitSchema // Try a single subcommit if array fails
			err := json.Unmarshal(text, &subcommit)
			if err != nil {
				return CommitAnalysisSchema{}, err
			}
			analysis.Subcommits = []SubcommitSchema{subcommit}
		} else {
			analysis.Subcommits = subcommits
		}
	}

	return *analysis, nil
}
