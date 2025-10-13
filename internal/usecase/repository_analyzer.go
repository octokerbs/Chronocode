package usecase

import (
	"context"
	"fmt"
	"sync"

	"github.com/octokerbs/chronocode-backend/internal/domain"
)

type RepositoryAnalyzer struct {
	Agent           domain.Agent
	CodeHostFactory domain.CodeHostFactory
	Database        domain.Database

	AnalyzedCommits      []domain.CommitRecord
	analyzedCommitsMutex sync.Mutex
}

func NewRepositoryAnalyzer(ctx context.Context, agent domain.Agent, factory domain.CodeHostFactory, database domain.Database) *RepositoryAnalyzer {
	return &RepositoryAnalyzer{
		Agent:           agent,
		CodeHostFactory: factory,
		Database:        database,
		AnalyzedCommits: []domain.CommitRecord{},
	}
}

func (ra *RepositoryAnalyzer) AnalyzeRepository(ctx context.Context, repoURL string, accessToken string) ([]domain.CommitRecord, []string, error) {
	codeHost, err := ra.CodeHostFactory.Create(ctx, accessToken)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create GitHub client: %w", err)
	}

	commits := make(chan string)
	records := make(chan domain.Record)
	errors := make(chan string)

	repoID, err := codeHost.RepositoryID(ctx, repoURL)
	if err != nil {
		return nil, nil, err
	}

	repositoryRecord, err := getOrCreateRepositoryRecord(ctx, repoURL, ra.Database, codeHost, repoID)
	if err != nil {
		return nil, nil, err
	}

	// Set up workers
	go func() {
		var wg sync.WaitGroup
		for range 200 {
			wg.Add(1)
			go ra.commitAnalyzerWorker(ctx, codeHost, repoURL, commits, records, &wg, errors)
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
		codeHost.ProduceCommits(ctx, repoURL, lastAnalyzedCommitSHA, commits, errors)
		close(commits)
	}()

	// Collect errors
	errorsSlice := []string{}
	for e := range errors {
		errorsSlice = append(errorsSlice, e)
	}

	return ra.AnalyzedCommits, errorsSlice, nil
}

func getOrCreateRepositoryRecord(ctx context.Context, repoURL string, database domain.Database, codeHost domain.CodeHost, id int64) (*domain.RepositoryRecord, error) {
	repo, ok, err := database.GetRepository(ctx, id)
	if err != nil {
		return nil, err
	}

	if ok {
		return repo, nil
	}

	repo, err = domain.NewRepositoryRecord(ctx, repoURL, codeHost)
	if err != nil {
		return nil, err
	}

	database.InsertRepository(ctx, repo)

	return repo, nil
}

func (ra *RepositoryAnalyzer) commitAnalyzerWorker(ctx context.Context, codeHost domain.CodeHost, repoURL string, commits <-chan string, records chan<- domain.Record, wg *sync.WaitGroup, errors chan<- string) {
	defer wg.Done()

	for commitSHA := range commits {
		diff, err := codeHost.GetCommitDiff(ctx, repoURL, commitSHA)
		if err != nil {
			errors <- fmt.Sprintf("commit diff failed: %s", err.Error())
			continue
		}

		prompt := domain.CommitAnalysisPrompt + diff

		tries := 3
		var text []byte
		for tries > 0 {
			text, err = ra.Agent.Generate(ctx, prompt)
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

		analysis, err := domain.UnmarshalCommitAnalysisSchemaOntoStruct(text)
		if err != nil {
			errors <- fmt.Sprintf("error unmarshaling response: %s", err.Error())
			continue
		}

		commitRecord, err := domain.NewCommitRecord(ctx, repoURL, codeHost, commitSHA, &analysis.Commit)
		if err != nil {
			errors <- fmt.Sprintf("error creating commit record: %s", err.Error())
			continue
		}

		records <- commitRecord

		for _, subcommit := range analysis.Subcommits {
			subcommitRecord := domain.NewSubcommitRecord(commitSHA, &subcommit)
			records <- subcommitRecord
		}

		ra.analyzedCommitsMutex.Lock()
		ra.AnalyzedCommits = append(ra.AnalyzedCommits, *commitRecord)
		ra.analyzedCommitsMutex.Unlock()
	}
}

func (ra *RepositoryAnalyzer) databaseInserterWorker(ctx context.Context, records <-chan domain.Record, wg *sync.WaitGroup, errors chan<- string) {
	defer wg.Done()

	for record := range records {
		err := record.InsertIntoDatabase(ctx, ra.Database)
		if err != nil {
			errors <- fmt.Sprintf("error uploading object to database: %s", err.Error())
		}
	}
}
