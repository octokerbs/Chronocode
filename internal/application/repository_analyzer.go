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

	ra.cleanCommits()

	commits := make(chan string)
	var wg sync.WaitGroup
	const numWorkers = 200
	log.Info("Starting commit analysis workers", "workerCount", numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		workerLog := log.With("workerID", i)
		go ra.commitAnalyzerWorker(ctx, repo.URL, codeHost, commits, &wg, workerLog)
	}

	go func() {
		log.Info("Starting commit SHA producer")
		codeHost.ProduceCommitSHAs(ctx, repo.URL, repo.LastAnalyzedCommit, commits)
		log.Info("Commit SHA producer finished")
	}()

	wg.Wait()
	log.Info("All commit analysis workers finished")

	if len(ra.analyzedCommits) == 0 {
		log.Info("No new commits found to store")
		return nil
	}

	log.Info("Storing analyzed commits", "commitCount", len(ra.analyzedCommits))
	if err := ra.Database.StoreCommits(ctx, ra.analyzedCommits); err != nil {
		log.Error("Failed to store commits in database", err)
		return err
	}

	lastCommitSHA := ra.analyzedCommits[len(ra.analyzedCommits)-1].SHA
	repo.UpdateLastAnalyzedCommit(lastCommitSHA)
	log.Info("Updating repository's last analyzed commit", "lastCommitSHA", lastCommitSHA)

	if err := ra.Database.StoreRepository(ctx, repo); err != nil {
		log.Error("Failed to update repository with last analyzed commit", err)
		return err
	}

	log.Info("Repository analysis finished successfully")
	return nil
}

func (ra *RepositoryAnalyzer) cleanCommits() {
	ra.analyzedCommits = []*domain.Commit{}
}

func (ra *RepositoryAnalyzer) fetchOrCreateRepository(ctx context.Context, repoURL string, codeHost domain.CodeHost, log Logger) (*domain.Repository, error) {
	id, err := codeHost.FetchRepositoryID(ctx, repoURL)
	if err != nil {
		log.Error("Failed to fetch repository ID from code host", err)
		return nil, err
	}

	log = log.With("repoID", id)
	log.Info("Fetched repository ID")

	repo, err := ra.Database.GetRepository(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrRepositoryNotFound) {
			log.Info("Repository not in Database, fetching from code host")
			repo, err = codeHost.FetchRepository(ctx, repoURL)
			if err != nil {
				log.Error("Failed to fetch new repository from code host", err)
				return nil, err
			}

			log.Info("Storing new repository in database")
			if err := ra.Database.StoreRepository(ctx, repo); err != nil {
				log.Error("Failed to store new repository in database", err)
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return repo, nil
}

func (ra *RepositoryAnalyzer) commitAnalyzerWorker(
	ctx context.Context,
	repoURL string,
	codeHost domain.CodeHost,
	commits <-chan string,
	wg *sync.WaitGroup,
	log Logger,
) {
	defer func() {
		wg.Done()
	}()

	for commitSHA := range commits {
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
