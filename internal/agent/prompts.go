package agent

var CommitAnalysisPrompt string = `
    	You are a Commit Expert Analyzer specializing in code analysis and software development patterns.
    	You will receive a Git Commit diff.
    	Your task is to given commit, identify the logical units of work ("SubCommits") within this single GitHub commit. 
    	The subcommits will have a title, idea, description, and type.
    
    	Now extract the subcommits from the following diff:
    
    	`
