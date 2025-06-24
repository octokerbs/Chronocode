package main

import (
	"context"
	"fmt"
	"time"

	"github.com/chrono-code-hackathon/chronocode-go/internal"
	"github.com/chrono-code-hackathon/chronocode-go/internal/services/githubapi"
	"github.com/chrono-code-hackathon/chronocode-go/performance/mocks"
)

func main() {
	analysisDuration := 2 * time.Second
	insertionDuration := 300 * time.Millisecond

	start := time.Now()

	repoURL := "https://github.com/octokerbs/TusLibros"
	accessToken := ""

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	generativeService := &mocks.AgentMock{AnalysisDuration: analysisDuration}
	githubService, err := githubapi.NewGithubClient(ctx, accessToken, repoURL)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	databaseService := &mocks.DatabaseMock{InsertionDuration: insertionDuration}

	repositoryAnalyzer, err := internal.NewRepositoryAnalyzer(ctx, generativeService, githubService, databaseService)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	commits, errors := repositoryAnalyzer.AnalyzeRepository(ctx)

	if 1 == 2 {
		fmt.Println(errors)
		return
	}

	elapsed := time.Since(start)

	for _, commit := range commits {
		fmt.Println(commit.SHA)
		fmt.Println(commit.Date)
		fmt.Println()
	}

	fmt.Println()
	fmt.Printf("Time elapsed analyzing commits concurrently: %v\n", elapsed)
}
