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
	Log             Logger

	resultsMutex    sync.Mutex
	analyzedCommits []*domain.Commit

	newHeadMutex sync.Mutex
	newHeadSHA   string
}

func NewRepositoryAnalyzer(
	ctx context.Context,
	agent domain.Agent,
	codehostFactory domain.CodeHostFactory,
	database domain.Database,
	log Logger,
) *RepositoryAnalyzer {
	ra := &RepositoryAnalyzer{
		Agent:           agent,
		CodeHostFactory: codehostFactory,
		Database:        database,
		Log:             log.With("service", "RepositoryAnalyzer"),
		analyzedCommits: []*domain.Commit{},
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
	var wg sync.WaitGroup
	const numWorkers = 200
	log.Info("Starting commit analysis workers", "workerCount", numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		workerLog := log.With("workerID", i)
		go ra.commitAnalyzerWorker(ctx, repo.URL, codeHost, commitSHAs, &wg, workerLog)
	}

	go func() {
		log.Info("Starting commit SHA producer")
		newHeadSHA, err := codeHost.ProduceCommitSHAs(ctx, repo.URL, repo.LastAnalyzedCommit, commitSHAs)
		if err != nil {
			log.Error("Commit SHA producer failed", err)
			// Producer is responsible for closing commitSHAs chan even on error
		} else if newHeadSHA != "" {
			ra.newHeadMutex.Lock()
			ra.newHeadSHA = newHeadSHA
			ra.newHeadMutex.Unlock()
			log.Info("Commit SHA producer identified new head", "newHeadSHA", newHeadSHA)
		}
		// Producer must close the channel to stop the workers
		log.Info("Commit SHA producer finished")
	}()

	wg.Wait()
	log.Info("All commit analysis workers finished")

	// This part can be removed when you move to persistence workers
	if len(ra.analyzedCommits) == 0 {
		log.Info("No new commits were successfully analyzed and collected")
	} else {
		log.Info("Storing analyzed commits", "commitCount", len(ra.analyzedCommits))
		if err := ra.Database.StoreCommits(ctx, ra.analyzedCommits); err != nil {
			log.Error("Failed to store commits in database", err)
			return err
		}
	}

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
	ra.analyzedCommits = []*domain.Commit{}
	ra.newHeadMutex.Lock()
	ra.newHeadSHA = "" // Reset for this run
	ra.newHeadMutex.Unlock()
}

func (ra *RepositoryAnalyzer) fetchOrCreateRepository(ctx context.Context, repoURL string, codeHost domain.CodeHost, log Logger) (*domain.Repository, error) {
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

	if !errors.Is(err, domain.ErrRepositoryNotFound) {
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
	wg *sync.WaitGroup,
	log Logger,
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

		ra.resultsMutex.Lock()
		ra.analyzedCommits = append(ra.analyzedCommits, commit)
		ra.resultsMutex.Unlock()
	}
}
