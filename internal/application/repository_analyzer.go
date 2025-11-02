package application

import (
	"context"
	"errors"
	"sync"

	"github.com/octokerbs/chronocode-backend/internal/domain"
)

type RepositoryAnalyzer struct {
	Agent           domain.Agent
	CodeHostFactory domain.CodeHostFactory
	Database        domain.Database
	Log             domain.Logger

	newHeadMutex sync.Mutex
	newHeadSHA   string
}

func NewRepositoryAnalyzer(
	ctx context.Context,
	agent domain.Agent,
	codehostFactory domain.CodeHostFactory,
	database domain.Database,
	log domain.Logger,
) *RepositoryAnalyzer {
	ra := &RepositoryAnalyzer{
		Agent:           agent,
		CodeHostFactory: codehostFactory,
		Database:        database,
		Log:             log.With("service", "RepositoryAnalyzer"),
	}

	return ra
}

func (ra *RepositoryAnalyzer) PrepareAnalysis(ctx context.Context, repoURL string, accessToken string) (*domain.Repository, domain.CodeHost, error) {
	log := ra.Log.With("repoURL", repoURL)
	log.Info("Preparing repository analysis")

	codeHost := ra.CodeHostFactory.Create(ctx, accessToken)

	repo, err := ra.fetchOrCreateRepository(ctx, repoURL, codeHost, log)
	if err != nil {
		log.Error("Failed to fetch or create repository", err)
		return nil, nil, err
	}

	log.Info("Repository validated successfully", "repoID", repo.ID)

	return repo, codeHost, nil
}

func (ra *RepositoryAnalyzer) RunAnalysis(ctx context.Context, repo *domain.Repository, codeHost domain.CodeHost) error {
	log := ra.Log.With("repoURL", repo.URL, "repoID", repo.ID)
	log.Info("Starting background analysis")

	ra.clean()

	commitSHAs := make(chan string)
	commits := make(chan *domain.Commit)

	var wgAnalyzers sync.WaitGroup
	var wgPersistency sync.WaitGroup
	const numAnalyzerWorkers = 20
	const numPersistencyWorkers = 40

	log.Info("Starting commit analysis workers", "workerCount", numAnalyzerWorkers)
	for i := 0; i < numAnalyzerWorkers; i++ {
		wgAnalyzers.Add(1)
		workerLog := log.With("analyzerWorkerID", i)
		go ra.commitAnalyzerWorker(ctx, repo.URL, codeHost, commitSHAs, commits, &wgAnalyzers, workerLog)
	}

	log.Info("Starting persistency workers", "workerCount", numPersistencyWorkers)
	for i := 0; i < numPersistencyWorkers; i++ {
		wgPersistency.Add(1)
		workerLog := log.With("persistencyWorkerID", i)
		go ra.commitPersistencyWorker(ctx, commits, &wgPersistency, workerLog)
	}

	go func() {
		defer close(commitSHAs)

		log.Info("Starting commit SHA producer")
		newHeadSHA, err := codeHost.ProduceCommitSHAs(ctx, repo.URL, repo.LastAnalyzedCommit, commitSHAs)
		if err != nil {
			log.Error("Commit SHA producer failed", err)
		} else if newHeadSHA != "" {
			ra.newHeadMutex.Lock()
			ra.newHeadSHA = newHeadSHA
			ra.newHeadMutex.Unlock()
			log.Info("Commit SHA producer identified new head", "newHeadSHA", newHeadSHA)
		}

		log.Info("Commit SHA producer finished")
	}()

	wgAnalyzers.Wait()
	close(commits)
	wgPersistency.Wait()

	log.Info("All commit analysis workers finished")

	ra.newHeadMutex.Lock()
	newHeadSHA := ra.newHeadSHA
	ra.newHeadMutex.Unlock()

	if newHeadSHA != "" && newHeadSHA != repo.LastAnalyzedCommit {
		log.Info("Updating repository's last analyzed commit", "lastCommitSHA", newHeadSHA)
		repo.UpdateLastAnalyzedCommit(newHeadSHA)
		if err := ra.Database.StoreRepository(ctx, repo); err != nil {
			log.Error("Failed to update repository with last analyzed commit", err)
			return err
		}
	} else if newHeadSHA == "" {
		log.Info("No new head commit SHA produced (no new commits or producer error)")
	} else {
		log.Info("New head SHA is the same as the last analyzed commit, no update needed.")
	}

	log.Info("Repository analysis finished successfully")
	return nil
}

func (ra *RepositoryAnalyzer) clean() {
	ra.newHeadMutex.Lock()
	ra.newHeadSHA = "" // Reset for this run
	ra.newHeadMutex.Unlock()
}

func (ra *RepositoryAnalyzer) fetchOrCreateRepository(ctx context.Context, repoURL string, codeHost domain.CodeHost, log domain.Logger) (*domain.Repository, error) {
	fetchedRepository, err := codeHost.FetchRepository(ctx, repoURL)
	if err != nil {
		log.Error("Failed to fetch repository from code host", err)
		return nil, err
	}

	log = log.With("repoID", fetchedRepository.ID)

	repo, err := ra.Database.GetRepository(ctx, fetchedRepository.ID)
	if err == nil {
		log.Info("Repository found in Database")
		return repo, nil
	}

	if !errors.Is(err, domain.ErrNotFound) {
		return nil, err // Maybe a server error
	}

	log.Info("Repository not in Database, storing new repository in database")
	if err := ra.Database.StoreRepository(ctx, fetchedRepository); err != nil {
		log.Error("Failed to store new repository in database", err)
		return nil, err
	}

	repo, err = ra.Database.GetRepository(ctx, fetchedRepository.ID)
	if err != nil {
		log.Error("Failed to store new repository in database", err)
		return nil, err
	}

	return repo, nil
}

func (ra *RepositoryAnalyzer) commitAnalyzerWorker(
	ctx context.Context,
	repoURL string,
	codeHost domain.CodeHost,
	commitSHAs <-chan string,
	commits chan<- *domain.Commit,
	wg *sync.WaitGroup,
	log domain.Logger,
) {
	defer func() {
		wg.Done()
	}()

	for commitSHA := range commitSHAs {
		commitLog := log.With("commitSHA", commitSHA)

		diff, err := codeHost.FetchCommitDiff(ctx, repoURL, commitSHA)
		if err != nil {
			commitLog.Warn("Failed to fetch commit diff, skipping commit", err)
			continue
		}

		analysis, err := ra.Agent.AnalyzeCommitDiff(ctx, diff)
		if err != nil {
			commitLog.Warn("Failed to analyze commit diff, skipping commit", err)
			continue
		}

		commit, err := codeHost.FetchCommit(ctx, repoURL, commitSHA)
		if err != nil {
			commitLog.Warn("Failed to fetch commit details, skipping commit", err)
			continue
		}

		commit.ApplyAnalysis(&analysis)

		commits <- commit
	}
}

func (ra *RepositoryAnalyzer) commitPersistencyWorker(ctx context.Context, commits <-chan *domain.Commit, wg *sync.WaitGroup, log domain.Logger) {
	defer func() {
		wg.Done()
	}()

	for commit := range commits {
		if err := ra.Database.StoreCommit(ctx, commit); err != nil {
			log.Error("Failed to store commit in database", err)
			continue
		}
	}
}
