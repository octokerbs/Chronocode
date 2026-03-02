package command

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"

	"github.com/octokerbs/chronocode/internal/domain/agent"
	"github.com/octokerbs/chronocode/internal/domain/analysis"
	"github.com/octokerbs/chronocode/internal/domain/codehost"
	"github.com/octokerbs/chronocode/internal/domain/repo"
	"github.com/octokerbs/chronocode/internal/domain/subcommit"
)

const maxConcurrentAnalyses = 10

type AnalyzeRepo struct {
	RepoURL     string
	AccessToken string
}

type AnalyzeRepoHandler struct {
	repoRepository      repo.Repository
	subcommitRepository subcommit.Repository
	agent               agent.Agent
	codeHostFactory     codehost.CodeHostFactory
	locker              analysis.Locker
}

func NewAnalyzeRepoHandler(repoRepository repo.Repository, subcommitRepository subcommit.Repository, agent agent.Agent, codeHostFactory codehost.CodeHostFactory, locker analysis.Locker) AnalyzeRepoHandler {
	return AnalyzeRepoHandler{repoRepository: repoRepository, subcommitRepository: subcommitRepository, agent: agent, codeHostFactory: codeHostFactory, locker: locker}
}

func (s *AnalyzeRepoHandler) Handle(ctx context.Context, cmd AnalyzeRepo) (int64, error) {
	slog.Info("AnalyzeRepo command received", "repo_url", cmd.RepoURL)

	codeHost, err := s.codeHostFactory.Create(ctx, cmd.AccessToken)
	if err != nil {
		slog.Error("Failed to create code host client", "repo_url", cmd.RepoURL, "error", err)
		return 0, err
	}

	slog.Debug("Checking repository access", "repo_url", cmd.RepoURL)
	if err := codeHost.CanAccessRepo(ctx, cmd.RepoURL); err != nil {
		slog.Warn("Repository access denied", "repo_url", cmd.RepoURL, "error", err)
		return 0, err
	}

	slog.Debug("Acquiring analysis lock", "repo_url", cmd.RepoURL)
	release, err := s.locker.Acquire(ctx, cmd.RepoURL)
	if err != nil {
		slog.Warn("Failed to acquire analysis lock - analysis already in progress", "repo_url", cmd.RepoURL)
		return 0, err
	}
	defer release()
	slog.Debug("Analysis lock acquired", "repo_url", cmd.RepoURL)

	newRepo, err := s.repoRepository.GetRepo(ctx, cmd.RepoURL)
	if err != nil {
		if !errors.Is(err, repo.ErrRepositoryNotFound) {
			slog.Error("Failed to look up repository", "repo_url", cmd.RepoURL, "error", err)
			return 0, err
		}

		slog.Info("Repository not found in database, creating from GitHub", "repo_url", cmd.RepoURL)
		newRepo, err = codeHost.CreateRepoFromURL(ctx, cmd.RepoURL)
		if err != nil {
			slog.Error("Failed to create repository from URL", "repo_url", cmd.RepoURL, "error", err)
			return 0, err
		}
		slog.Info("Repository created from GitHub", "repo_id", newRepo.ID(), "repo_name", newRepo.Name())
	} else {
		slog.Info("Existing repository found", "repo_id", newRepo.ID(), "repo_name", newRepo.Name(), "last_analyzed_sha", newRepo.LastAnalyzedCommitSHA())
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	commitRefs := make(chan codehost.CommitReference, 100)
	subcommits := make(chan subcommit.Subcommit, 100)

	var fetchErr, analysisErr, storageErr error
	var headSHA string
	wg.Add(3)

	slog.Info("Starting analysis pipeline", "repo_id", newRepo.ID(), "repo_url", cmd.RepoURL)

	go func() {
		defer wg.Done()
		defer close(commitRefs)
		headSHA, fetchErr = codeHost.GetRepoCommitSHAsIntoChannel(ctx, newRepo, commitRefs)
		if fetchErr != nil {
			slog.Error("Commit fetch pipeline failed", "repo_id", newRepo.ID(), "error", fetchErr)
			cancel()
		}
	}()

	go func() {
		defer wg.Done()
		defer close(subcommits)
		analysisErr = s.analyzeCommits(ctx, codeHost, newRepo, commitRefs, subcommits)
	}()

	go func() {
		defer wg.Done()
		storageErr = s.subcommitRepository.StoreSubcommits(ctx, subcommits)
		if storageErr != nil {
			slog.Error("Subcommit storage pipeline failed", "repo_id", newRepo.ID(), "error", storageErr)
		}
	}()

	wg.Wait()

	if fetchErr != nil {
		return 0, fetchErr
	}

	if analysisErr == nil && storageErr == nil && headSHA != "" {
		slog.Info("All commits analyzed successfully, updating last analyzed SHA", "repo_id", newRepo.ID(), "head_sha", headSHA)
		newRepo.SetLastAnalyzedCommitSHA(headSHA)
	} else if analysisErr != nil {
		slog.Warn("Analysis completed with errors, not updating last analyzed SHA", "repo_id", newRepo.ID(), "error", analysisErr)
	}

	if err := s.repoRepository.StoreRepo(ctx, newRepo); err != nil {
		slog.Error("Failed to store repository after analysis", "repo_id", newRepo.ID(), "error", err)
		return 0, err
	}

	slog.Info("AnalyzeRepo command completed", "repo_id", newRepo.ID(), "repo_url", cmd.RepoURL, "head_sha", headSHA)

	return newRepo.ID(), errors.Join(analysisErr, storageErr)
}

func (s *AnalyzeRepoHandler) HandleAsync(ctx context.Context, cmd AnalyzeRepo) (int64, error) {
	slog.Info("AnalyzeRepo async command received", "repo_url", cmd.RepoURL)

	codeHost, err := s.codeHostFactory.Create(ctx, cmd.AccessToken)
	if err != nil {
		slog.Error("Failed to create code host client", "repo_url", cmd.RepoURL, "error", err)
		return 0, err
	}

	slog.Debug("Checking repository access", "repo_url", cmd.RepoURL)
	if err := codeHost.CanAccessRepo(ctx, cmd.RepoURL); err != nil {
		slog.Warn("Repository access denied", "repo_url", cmd.RepoURL, "error", err)
		return 0, err
	}

	slog.Debug("Acquiring analysis lock", "repo_url", cmd.RepoURL)
	release, err := s.locker.Acquire(ctx, cmd.RepoURL)
	if err != nil {
		slog.Warn("Failed to acquire analysis lock - analysis already in progress", "repo_url", cmd.RepoURL)
		return 0, err
	}
	slog.Debug("Analysis lock acquired", "repo_url", cmd.RepoURL)

	newRepo, err := s.repoRepository.GetRepo(ctx, cmd.RepoURL)
	if err != nil {
		if !errors.Is(err, repo.ErrRepositoryNotFound) {
			slog.Error("Failed to look up repository", "repo_url", cmd.RepoURL, "error", err)
			release()
			return 0, err
		}

		slog.Info("Repository not found in database, creating from GitHub", "repo_url", cmd.RepoURL)
		newRepo, err = codeHost.CreateRepoFromURL(ctx, cmd.RepoURL)
		if err != nil {
			slog.Error("Failed to create repository from URL", "repo_url", cmd.RepoURL, "error", err)
			release()
			return 0, err
		}
		slog.Info("Repository created from GitHub", "repo_id", newRepo.ID(), "repo_name", newRepo.Name())
	} else {
		slog.Info("Existing repository found", "repo_id", newRepo.ID(), "repo_name", newRepo.Name(), "last_analyzed_sha", newRepo.LastAnalyzedCommitSHA())
	}

	if err := s.repoRepository.StoreRepo(ctx, newRepo); err != nil {
		slog.Error("Failed to store repository before analysis", "repo_id", newRepo.ID(), "error", err)
		release()
		return 0, err
	}

	go func() {
		defer release()
		s.runPipeline(codeHost, newRepo)
	}()

	slog.Info("AnalyzeRepo async command returning immediately", "repo_id", newRepo.ID(), "repo_url", cmd.RepoURL)
	return newRepo.ID(), nil
}

func (s *AnalyzeRepoHandler) runPipeline(codeHost codehost.CodeHost, targetRepo *repo.Repo) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	commitRefs := make(chan codehost.CommitReference, 100)
	subcommits := make(chan subcommit.Subcommit, 100)

	var fetchErr, analysisErr, storageErr error
	var headSHA string
	wg.Add(3)

	slog.Info("Starting async analysis pipeline", "repo_id", targetRepo.ID(), "repo_url", targetRepo.URL())

	go func() {
		defer wg.Done()
		defer close(commitRefs)
		headSHA, fetchErr = codeHost.GetRepoCommitSHAsIntoChannel(ctx, targetRepo, commitRefs)
		if fetchErr != nil {
			slog.Error("Commit fetch pipeline failed", "repo_id", targetRepo.ID(), "error", fetchErr)
			cancel()
		}
	}()

	go func() {
		defer wg.Done()
		defer close(subcommits)
		analysisErr = s.analyzeCommits(ctx, codeHost, targetRepo, commitRefs, subcommits)
	}()

	go func() {
		defer wg.Done()
		storageErr = s.subcommitRepository.StoreSubcommits(ctx, subcommits)
		if storageErr != nil {
			slog.Error("Subcommit storage pipeline failed", "repo_id", targetRepo.ID(), "error", storageErr)
		}
	}()

	wg.Wait()

	if fetchErr != nil {
		slog.Error("Async analysis pipeline failed during fetch", "repo_id", targetRepo.ID(), "error", fetchErr)
		return
	}

	if analysisErr == nil && storageErr == nil && headSHA != "" {
		slog.Info("All commits analyzed successfully, updating last analyzed SHA", "repo_id", targetRepo.ID(), "head_sha", headSHA)
		targetRepo.SetLastAnalyzedCommitSHA(headSHA)
	} else if analysisErr != nil {
		slog.Warn("Async analysis completed with errors, not updating last analyzed SHA", "repo_id", targetRepo.ID(), "error", analysisErr)
	}

	if err := s.repoRepository.StoreRepo(ctx, targetRepo); err != nil {
		slog.Error("Failed to store repository after async analysis", "repo_id", targetRepo.ID(), "error", err)
	}

	slog.Info("Async analysis pipeline completed", "repo_id", targetRepo.ID(), "head_sha", headSHA)
}

func (s *AnalyzeRepoHandler) analyzeCommits(ctx context.Context, codeHost codehost.CodeHost, r *repo.Repo, commitRefs <-chan codehost.CommitReference, subcommits chan<- subcommit.Subcommit) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error
	sem := make(chan struct{}, maxConcurrentAnalyses)

	var totalCommits, analyzedCommits, skippedCommits, failedCommits atomic.Int64

	for ref := range commitRefs {
		if ctx.Err() != nil {
			break
		}

		totalCommits.Add(1)
		sem <- struct{}{}

		wg.Add(1)
		go func(ref codehost.CommitReference) {
			defer func() { <-sem }()
			defer wg.Done()

			if ctx.Err() != nil {
				return
			}

			alreadyAnalyzed, err := s.subcommitRepository.HasSubcommitsForCommit(ctx, r.ID(), ref.SHA)
			if err != nil {
				failedCommits.Add(1)
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
				return
			}
			if alreadyAnalyzed {
				skippedCommits.Add(1)
				slog.Debug("Commit already analyzed, skipping", "repo_id", r.ID(), "commit_sha", ref.SHA)
				return
			}

			slog.Debug("Analyzing commit", "repo_id", r.ID(), "commit_sha", ref.SHA)

			diff, err := codeHost.GetCommitDiff(ctx, r, ref.SHA)
			if err != nil {
				failedCommits.Add(1)
				slog.Error("Failed to fetch commit diff", "repo_id", r.ID(), "commit_sha", ref.SHA, "error", err)
				mu.Lock()
				errs = append(errs, fmt.Errorf("%w: %v", codehost.ErrDiffFetchFailed, err))
				mu.Unlock()
				return
			}

			results, err := s.agent.AnalyzeDiff(ctx, diff)
			if err != nil {
				failedCommits.Add(1)
				slog.Error("Agent failed to analyze commit diff", "repo_id", r.ID(), "commit_sha", ref.SHA, "error", err)
				mu.Lock()
				errs = append(errs, fmt.Errorf("%w: %v", agent.ErrAnalysisFailed, err))
				mu.Unlock()
				return
			}

			analyzedCommits.Add(1)
			slog.Debug("Commit analyzed", "repo_id", r.ID(), "commit_sha", ref.SHA, "subcommits_produced", len(results))

			for _, result := range results {
				subcommits <- subcommit.NewSubcommit(result.Title, result.Idea, result.Description, result.Epic, result.ModificationType, ref.SHA, result.Files, r.ID(), ref.CommittedAt)
			}
		}(ref)
	}

	wg.Wait()

	slog.Info("Commit analysis pipeline completed",
		"repo_id", r.ID(),
		"total_commits", totalCommits.Load(),
		"analyzed", analyzedCommits.Load(),
		"skipped", skippedCommits.Load(),
		"failed", failedCommits.Load(),
	)

	return errors.Join(errs...)
}
