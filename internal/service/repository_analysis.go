package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/octokerbs/chronocode-go/internal/domain/agent"
	"github.com/octokerbs/chronocode-go/internal/domain/codehost"
	"github.com/octokerbs/chronocode-go/internal/domain/codehost/githubapi"
	"github.com/octokerbs/chronocode-go/internal/repository"
)

type RepositoryAnalyzer struct {
	AgentClient    agent.AgentClient
	DatabaseClient repository.DatabaseClient

	AnalyzedCommits      []repository.CommitRecord
	analyzedCommitsMutex sync.Mutex
}

func NewRepositoryAnalyzer(ctx context.Context, agent agent.AgentClient, databaseClient repository.DatabaseClient) *RepositoryAnalyzer {
	return &RepositoryAnalyzer{
		AgentClient:     agent,
		DatabaseClient:  databaseClient,
		AnalyzedCommits: []repository.CommitRecord{},
	}
}

func (ra *RepositoryAnalyzer) AnalyzeRepository(ctx context.Context, repoURL string, accessToken string) ([]repository.CommitRecord, []string, error) {
	codeHostClient, err := githubapi.NewGithubClient(ctx, accessToken)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create GitHub client: %w", err)
	}

	commits := make(chan string)
	records := make(chan repository.Record)
	errors := make(chan string)

	repoID, err := codeHostClient.RepositoryID(ctx, repoURL)
	if err != nil {
		return nil, nil, err
	}

	repositoryRecord, err := getOrCreateRepositoryRecord(ctx, repoURL, ra.DatabaseClient, codeHostClient, repoID)
	if err != nil {
		return nil, nil, err
	}

	// Set up workers
	go func() {
		var wg sync.WaitGroup
		for range 200 {
			wg.Add(1)
			go ra.commitAnalyzerWorker(ctx, codeHostClient, repoURL, commits, records, &wg, errors)
		}
		wg.Wait()
		close(records)
	}()

	go func() {
		var wg sync.WaitGroup
		for range 500 {
			wg.Add(1)
			go ra.databaseInserterWorker(ctx, records, &wg, errors)
		}
		wg.Wait()
		close(errors)
	}()

	lastAnalyzedCommitSHA := repositoryRecord.LastAnalyzedCommit

	// Start pipeline
	go func() {
		codeHostClient.ProduceCommits(ctx, repoURL, lastAnalyzedCommitSHA, commits, errors)
		close(commits)
	}()

	// Collect errors
	errorsSlice := []string{}
	for e := range errors {
		errorsSlice = append(errorsSlice, e)
	}

	return ra.AnalyzedCommits, errorsSlice, nil
}

func getOrCreateRepositoryRecord(ctx context.Context, repoURL string, dbClient repository.DatabaseClient, codeHostClient codehost.CodeHostClient, id int64) (*repository.RepositoryRecord, error) {
	repo, ok, err := dbClient.GetRepository(ctx, id)
	if err != nil {
		return nil, err
	}

	if ok {
		return repo, nil
	}

	repo, err = repository.NewRepositoryRecord(ctx, repoURL, codeHostClient)
	if err != nil {
		return nil, err
	}

	dbClient.InsertRepository(ctx, repo)

	return repo, nil
}

func (ra *RepositoryAnalyzer) commitAnalyzerWorker(ctx context.Context, codeHostClient codehost.CodeHostClient, repoURL string, commits <-chan string, records chan<- repository.Record, wg *sync.WaitGroup, errors chan<- string) {
	defer wg.Done()

	for commitSHA := range commits {
		diff, err := codeHostClient.GetCommitDiff(ctx, repoURL, commitSHA)
		if err != nil {
			errors <- fmt.Sprintf("commit diff failed: %s", err.Error())
			continue
		}

		prompt := agent.CommitAnalysisPrompt + diff

		tries := 3
		var text []byte
		for tries > 0 {
			text, err = ra.AgentClient.Generate(ctx, prompt)
			if err != nil {
				tries--
				continue
			}
			break
		}

		if tries == 0 {
			errors <- "no text response parts found for commit"
			continue
		}

		analysis, err := agent.UnmarshalCommitAnalysisSchemaOntoStruct(text)
		if err != nil {
			errors <- fmt.Sprintf("error unmarshaling response: %s", err.Error())
			continue
		}

		commitRecord, err := repository.NewCommitRecord(ctx, repoURL, codeHostClient, commitSHA, &analysis.Commit)
		if err != nil {
			errors <- fmt.Sprintf("error creating commit record: %s", err.Error())
			continue
		}
		records <- commitRecord

		for _, subcommit := range analysis.Subcommits {
			subcommitRecord := repository.NewSubcommitRecord(commitSHA, &subcommit)
			records <- subcommitRecord
		}

		ra.analyzedCommitsMutex.Lock()
		ra.AnalyzedCommits = append(ra.AnalyzedCommits, *commitRecord)
		ra.analyzedCommitsMutex.Unlock()
	}
}

func (ra *RepositoryAnalyzer) databaseInserterWorker(ctx context.Context, records <-chan repository.Record, wg *sync.WaitGroup, errors chan<- string) {
	defer wg.Done()

	for record := range records {
		err := record.InsertIntoDatabase(ctx, ra.DatabaseClient)
		if err != nil {
			errors <- fmt.Sprintf("error uploading object to database: %s", err.Error())
		}
	}
}
