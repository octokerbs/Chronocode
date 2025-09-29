package internal

import (
	"context"
	"fmt"
	"sync"

	"github.com/octokerbs/chronocode-go/internal/domain/agent"
	"github.com/octokerbs/chronocode-go/internal/domain/sourcecodehost"
	"github.com/octokerbs/chronocode-go/internal/repository"
)

type RepositoryAnalyzer struct {
	GenerativeAgentService agent.GenerativeAgentService
	SourceCodeHostService  sourcecodehost.SourcecodeHostService
	DatabaseService        repository.DatabaseService
	RepositoryRecord       *repository.RepositoryRecord

	AnalyzedCommits      []repository.CommitRecord
	analyzedCommitsMutex sync.Mutex
}

func NewRepositoryAnalyzer(ctx context.Context, generativeAgentService agent.GenerativeAgentService, sourceCodeHostService sourcecodehost.SourcecodeHostService, databaseService repository.DatabaseService) (*RepositoryAnalyzer, error) {
	repositoryRecord, err := getOrCreateRepositoryRecord(ctx, databaseService, sourceCodeHostService, sourceCodeHostService.RepositoryID())
	if err != nil {
		return nil, err
	}

	return &RepositoryAnalyzer{
		GenerativeAgentService: generativeAgentService,
		SourceCodeHostService:  sourceCodeHostService,
		DatabaseService:        databaseService,
		RepositoryRecord:       repositoryRecord,
		AnalyzedCommits:        []repository.CommitRecord{},
	}, nil
}

func (ra *RepositoryAnalyzer) AnalyzeRepository(ctx context.Context) ([]repository.CommitRecord, []error) {
	commits := make(chan string)
	records := make(chan repository.Record)
	errors := make(chan error)

	// Set up workers
	go func() {
		var wg sync.WaitGroup
		for range 200 {
			wg.Add(1)
			go ra.commitAnalyzerWorker(ctx, commits, records, &wg, errors)
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

	lastAnalyzedCommitSHA := ra.RepositoryRecord.LastAnalyzedCommit

	// Start pipeline
	go func() {
		ra.SourceCodeHostService.ProduceCommits(ctx, lastAnalyzedCommitSHA, commits, errors)
		close(commits)
	}()

	// Collect errors
	errorsSlice := []error{}
	for e := range errors {
		errorsSlice = append(errorsSlice, e)
	}

	return ra.AnalyzedCommits, errorsSlice
}

func getOrCreateRepositoryRecord(ctx context.Context, databaseService repository.DatabaseService, sourceCodeHostService sourcecodehost.SourcecodeHostService, id int64) (*repository.RepositoryRecord, error) {
	repo, ok, err := databaseService.GetRepository(ctx, id)
	if err != nil {
		return nil, err
	}

	if ok {
		return repo, nil
	}

	repo, err = repository.NewRepositoryRecord(sourceCodeHostService)
	if err != nil {
		return nil, err
	}

	databaseService.InsertRepository(ctx, repo)

	return repo, nil
}

func (ra *RepositoryAnalyzer) commitAnalyzerWorker(ctx context.Context, commits <-chan string, records chan<- repository.Record, wg *sync.WaitGroup, errors chan<- error) {
	defer wg.Done()

	for commitSHA := range commits {
		diff, err := ra.SourceCodeHostService.GetCommitDiff(ctx, commitSHA)
		if err != nil {
			errors <- fmt.Errorf("commit diff failed: %s", err.Error())
			continue
		}

		prompt := agent.CommitAnalysisPrompt + diff

		tries := 3
		var text []byte
		for tries > 0 {
			text, err = ra.GenerativeAgentService.Generate(ctx, prompt)
			if err != nil {
				tries--
				continue
			}
			break
		}

		if tries == 0 {
			errors <- fmt.Errorf("no text response parts found for commit")
			continue
		}

		analysis, err := agent.UnmarshalCommitAnalysisSchemaOntoStruct(text)
		if err != nil {
			errors <- fmt.Errorf("error unmarshaling response: %s", err.Error())
			continue
		}

		commitRecord, err := repository.NewCommitRecord(ctx, ra.SourceCodeHostService, commitSHA, &analysis.Commit)
		if err != nil {
			errors <- fmt.Errorf("error creating commit record: %s", err.Error())
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

func (ra *RepositoryAnalyzer) databaseInserterWorker(ctx context.Context, records <-chan repository.Record, wg *sync.WaitGroup, errors chan<- error) {
	defer wg.Done()

	for record := range records {
		err := record.InsertIntoDatabase(ctx, ra.DatabaseService)
		if err != nil {
			errors <- fmt.Errorf("error uploading object to database: %s", err.Error())
		}
	}
}
