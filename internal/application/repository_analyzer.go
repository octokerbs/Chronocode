package application

import (
	"context"
	"fmt"
	"sync"

	"github.com/octokerbs/chronocode-backend/internal/domain"
)

type CodeHostFactory interface {
	Create(ctx context.Context, accessToken string) (CodeHost, error)
}

type CodeHost interface {
	NewRepository(ctx context.Context, repoURL string) (*domain.Repository, error)
	NewCommit(ctx context.Context, repoURL string, commitSHA string) (*domain.Commit, error)

	ProduceCommits(ctx context.Context, repoURL string, lastAnalyzedCommitSHA string, commits chan<- string, errors chan<- string)
	GetCommitDiff(ctx context.Context, repoURL string, commitSHA string) (string, error)

	RepositoryID(ctx context.Context, repoURL string) (int64, error)
}

type Agent interface {
	Generate(ctx context.Context, prompt string) ([]byte, error)
}

type DatabaseRecord interface {
	IsDatabaseRecord()
}

type Database interface {
	InsertRepositoryRecord(ctx context.Context, repo *domain.Repository) error
	InsertCommitRecord(ctx context.Context, commit *domain.Commit) error
	InsertSubcommitRecord(ctx context.Context, subcommit *domain.Subcommit) error

	GetRepository(ctx context.Context, id int64) (*domain.Repository, bool, error)
	ProcessRecords(ctx context.Context, records <-chan DatabaseRecord, errors chan<- string)
}

type RepositoryAnalyzer struct {
	Agent           Agent
	CodeHostFactory CodeHostFactory
	Database        Database

	AnalyzedCommits      []domain.Commit
	analyzedCommitsMutex sync.Mutex
}

func NewRepositoryAnalyzer(ctx context.Context, agent Agent, factory CodeHostFactory, database Database) *RepositoryAnalyzer {
	return &RepositoryAnalyzer{
		Agent:           agent,
		CodeHostFactory: factory,
		Database:        database,
		AnalyzedCommits: []domain.Commit{},
	}
}

func (ra *RepositoryAnalyzer) AnalyzeRepository(ctx context.Context, repoURL string, accessToken string) ([]domain.Commit, []string, error) {
	codeHost, err := ra.CodeHostFactory.Create(ctx, accessToken)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create GitHub client: %w", err)
	}

	commits := make(chan string)
	records := make(chan DatabaseRecord)
	errors := make(chan string)

	repo, err := getOrCreateRepositoryRecord(ctx, repoURL, ra.Database, codeHost)
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

	lastAnalyzedCommitSHA := repo.LastAnalyzedCommit

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

func getOrCreateRepositoryRecord(ctx context.Context, repoURL string, database Database, codeHost CodeHost) (*domain.Repository, error) {
	id, err := codeHost.RepositoryID(ctx, repoURL)
	if err != nil {
		return nil, err
	}

	repo, ok, err := database.GetRepository(ctx, id)
	if err != nil {
		return nil, err
	}

	if ok {
		return repo, nil
	}

	repo, err = codeHost.NewRepository(ctx, repoURL)
	if err != nil {
		return nil, err
	}

	database.InsertRepositoryRecord(ctx, repo)

	return repo, nil
}

func (ra *RepositoryAnalyzer) commitAnalyzerWorker(ctx context.Context, codeHost CodeHost, repoURL string, commits <-chan string, records chan<- DatabaseRecord, wg *sync.WaitGroup, errors chan<- string) {
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

		// TODO: Primero crear el commit, despues cargarle la data como la descripcion o los subcommits.
		analysis, err := domain.UnmarshalCommitAnalysisSchemaOntoStruct(text)
		if err != nil {
			errors <- fmt.Sprintf("error unmarshaling response: %s", err.Error())
			continue
		}

		commitRecord, err := codeHost.NewCommit(ctx, repoURL, commitSHA)
		if err != nil {
			errors <- fmt.Sprintf("error creating commit record: %s", err.Error())
			continue
		}

		records <- commitRecord

		for _, subcommit := range analysis.Subcommits {
			subcommit.CommitSHA = commitSHA
			records <- &subcommit
		}

		ra.analyzedCommitsMutex.Lock()
		ra.AnalyzedCommits = append(ra.AnalyzedCommits, *commitRecord)
		ra.analyzedCommitsMutex.Unlock()
	}
}

func (ra *RepositoryAnalyzer) databaseInserterWorker(ctx context.Context, records <-chan DatabaseRecord, wg *sync.WaitGroup, errors chan<- string) {
	defer wg.Done()
	ra.Database.ProcessRecords(ctx, records, errors)
}
