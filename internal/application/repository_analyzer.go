package application

import (
	"context"
	"fmt"
	"sync"

	"github.com/octokerbs/chronocode-backend/internal/domain"
)

// CodeHost is the service that allocates our code and gives us access to different aspects of it
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

// Agent is the LLM that we use to process our commits giving us descriptions and generating subcommits
type Agent interface {
	AnalyzeDiff(ctx context.Context, diff string) (domain.CommitAnalysis, error)
}

// Database is where we store our repositories, commit data and subcommit data.
type Database interface {
	InsertRepositoryRecord(ctx context.Context, repo *domain.Repository) error
	InsertCommitRecord(ctx context.Context, commit *domain.Commit) error
	InsertSubcommitRecord(ctx context.Context, subcommit *domain.Subcommit) error

	GetRepository(ctx context.Context, id int64) (*domain.Repository, bool, error)
	ProcessRecords(ctx context.Context, records <-chan DatabaseRecord, errors chan<- string)
}

type DatabaseRecord interface {
	IsDatabaseRecord()
}

// RepositoryAnalyzer fetches our repo and gives us an analysis on all the commits, generating also the subcommits.
// This is our principal orquestator for the app.
type RepositoryAnalyzer struct {
	Agent           Agent
	CodeHostFactory CodeHostFactory
	Database        Database

	AnalyzedCommits      []domain.Commit
	AnalyzerSubcommits   []domain.Subcommit
	analyzedCommitsMutex sync.Mutex
}

func NewRepositoryAnalyzer(ctx context.Context, agent Agent, factory CodeHostFactory, database Database) *RepositoryAnalyzer {
	return &RepositoryAnalyzer{
		Agent:              agent,
		CodeHostFactory:    factory,
		Database:           database,
		AnalyzedCommits:    []domain.Commit{},
		AnalyzerSubcommits: []domain.Subcommit{},
	}
}

func (ra *RepositoryAnalyzer) AnalyzeRepository(ctx context.Context, repoURL string, accessToken string) ([]domain.Commit, []domain.Subcommit, []string, error) {
	codeHost, err := ra.CodeHostFactory.Create(ctx, accessToken)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create GitHub client: %w", err)
	}

	commits := make(chan string)
	records := make(chan DatabaseRecord)
	errors := make(chan string)

	repo, err := getOrCreateRepositoryRecord(ctx, repoURL, ra.Database, codeHost)
	if err != nil {
		return nil, nil, nil, err
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

	return ra.AnalyzedCommits, ra.AnalyzerSubcommits, errorsSlice, nil
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

		analysis, err := ra.Agent.AnalyzeDiff(ctx, diff)
		if err != nil {
			errors <- fmt.Sprintf("error unmarshaling response: %s", err.Error())
			continue
		}

		commit, err := codeHost.NewCommit(ctx, repoURL, commitSHA)
		if err != nil {
			errors <- fmt.Sprintf("error creating commit record: %s", err.Error())
			continue
		}

		commit.Description = analysis.Commit.Description

		subcommits := analysis.Subcommits
		for i := range subcommits {
			subcommits[i].CommitSHA = commitSHA
			records <- &subcommits[i]
		}

		ra.saveCommitAnalysis(commit, subcommits, records)

		ra.analyzedCommitsMutex.Lock()
		ra.AnalyzedCommits = append(ra.AnalyzedCommits, *commit)
		ra.AnalyzerSubcommits = append(ra.AnalyzerSubcommits, subcommits...)
		ra.analyzedCommitsMutex.Unlock()
	}
}

func (ra *RepositoryAnalyzer) saveCommitAnalysis(commit *domain.Commit, subcommits []domain.Subcommit, records chan<- DatabaseRecord) {
	records <- commit
	for _, subcommit := range subcommits {
		records <- &subcommit
	}
}

func (ra *RepositoryAnalyzer) databaseInserterWorker(ctx context.Context, records <-chan DatabaseRecord, wg *sync.WaitGroup, errors chan<- string) {
	defer wg.Done()
	ra.Database.ProcessRecords(ctx, records, errors)
}
