package domain

import (
	"encoding/json"
)

type CommitAnalysis struct {
	Commit     Commit      `json:"commit"`
	Subcommits []Subcommit `json:"subcommits"`
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

var CommitAnalysisPrompt string = `
    	You are a Commit Expert Analyzer specializing in code analysis and software development patterns.
    	You will receive a Git Commit diff.
    	Your task is to given commit, identify the logical units of work ("SubCommits") within this single GitHub commit. 
    	The subcommits will have a title, idea, description, and type.

    	Now extract the subcommits from the following diff:
    	`
