package adapters

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/octokerbs/chronocode/internal/domain/codehost"
	"github.com/octokerbs/chronocode/internal/domain/repo"

	"golang.org/x/oauth2"
)

type GithubCodeHostFactory struct{}

func NewGithubCodeHostFactory() *GithubCodeHostFactory {
	return &GithubCodeHostFactory{}
}

func (f *GithubCodeHostFactory) Create(ctx context.Context, accessToken string) (codehost.CodeHost, error) {
	if accessToken == "" {
		slog.Warn("GitHub code host creation failed - empty access token")
		return nil, codehost.ErrAccessDenied
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	slog.Debug("GitHub code host client created")
	return &GithubCodeHost{client: client}, nil
}

type GithubCodeHost struct {
	client *github.Client
}

func (gc *GithubCodeHost) CanAccessRepo(ctx context.Context, repoURL string) error {
	owner, repoName, err := parseRepoURL(repoURL)
	if err != nil {
		slog.Warn("Invalid repo URL for access check", "repo_url", repoURL, "error", err)
		return codehost.ErrInvalidRepoURL
	}

	slog.Debug("Checking GitHub repo access", "owner", owner, "repo", repoName)
	_, resp, err := gc.client.Repositories.Get(ctx, owner, repoName)
	if err != nil {
		if resp != nil && (resp.StatusCode == 404 || resp.StatusCode == 403 || resp.StatusCode == 401) {
			slog.Warn("GitHub repo access denied", "owner", owner, "repo", repoName, "status", resp.StatusCode)
			return codehost.ErrAccessDenied
		}
		slog.Error("GitHub API error during access check", "owner", owner, "repo", repoName, "error", err)
		return err
	}

	slog.Debug("GitHub repo access confirmed", "owner", owner, "repo", repoName)
	return nil
}

func (gc *GithubCodeHost) CreateRepoFromURL(ctx context.Context, repoURL string) (*repo.Repo, error) {
	owner, repoName, err := parseRepoURL(repoURL)
	if err != nil {
		return nil, codehost.ErrInvalidRepoURL
	}

	slog.Debug("Fetching GitHub repo metadata", "owner", owner, "repo", repoName)
	ghRepo, _, err := gc.client.Repositories.Get(ctx, owner, repoName)
	if err != nil {
		slog.Error("Failed to fetch GitHub repo metadata", "owner", owner, "repo", repoName, "error", err)
		return nil, err
	}

	slog.Info("GitHub repo metadata fetched", "repo_id", *ghRepo.ID, "full_name", *ghRepo.FullName)
	return repo.NewRepo(*ghRepo.ID, *ghRepo.FullName, repoURL, "", time.Now()), nil
}

func (gc *GithubCodeHost) GetAuthenticatedUser(ctx context.Context) (*codehost.UserProfile, error) {
	slog.Debug("Fetching authenticated GitHub user")
	user, _, err := gc.client.Users.Get(ctx, "")
	if err != nil {
		slog.Error("Failed to fetch authenticated GitHub user", "error", err)
		return nil, err
	}

	profile := &codehost.UserProfile{
		ID:    int64(*user.ID),
		Login: *user.Login,
	}
	if user.Name != nil {
		profile.Name = *user.Name
	}
	if user.AvatarURL != nil {
		profile.AvatarURL = *user.AvatarURL
	}
	if user.Email != nil {
		profile.Email = *user.Email
	}
	slog.Info("Authenticated GitHub user fetched", "login", profile.Login, "user_id", profile.ID)
	return profile, nil
}

func (gc *GithubCodeHost) SearchRepositories(ctx context.Context, query string) ([]codehost.RepoSearchResult, error) {
	slog.Debug("Searching GitHub repositories", "query", query)
	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 20},
		Sort:        "updated",
	}

	repos, _, err := gc.client.Repositories.List(ctx, "", opts)
	if err != nil {
		slog.Error("Failed to list GitHub repositories", "query", query, "error", err)
		return nil, err
	}

	var results []codehost.RepoSearchResult
	for _, r := range repos {
		if r.FullName == nil || r.HTMLURL == nil {
			continue
		}
		if query != "" && !strings.Contains(strings.ToLower(*r.FullName), strings.ToLower(query)) {
			continue
		}
		results = append(results, codehost.RepoSearchResult{
			ID:   int64(*r.ID),
			Name: *r.FullName,
			URL:  *r.HTMLURL,
		})
	}
	slog.Info("GitHub repository search completed", "query", query, "fetched", len(repos), "matched", len(results))
	return results, nil
}

func (gc *GithubCodeHost) GetRepoCommitSHAsIntoChannel(ctx context.Context, r *repo.Repo, commits chan<- codehost.CommitReference) (string, error) {
	owner, repoName, err := parseRepoURL(r.URL())
	if err != nil {
		return "", codehost.ErrInvalidRepoURL
	}

	lastSHA := r.LastAnalyzedCommitSHA()
	opts := &github.CommitsListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	slog.Info("Fetching commits from GitHub", "owner", owner, "repo", repoName, "last_analyzed_sha", lastSHA)

	var headSHA string
	var totalFetched, sentCount, mergeSkipped int
	page := 0
	for {
		page++
		pageCommits, resp, err := gc.client.Repositories.ListCommits(ctx, owner, repoName, opts)
		if err != nil {
			slog.Error("Failed to fetch commits page from GitHub", "owner", owner, "repo", repoName, "page", page, "error", err)
			return "", err
		}

		slog.Debug("Fetched commits page", "owner", owner, "repo", repoName, "page", page, "count", len(pageCommits))

		for _, commit := range pageCommits {
			totalFetched++
			if commit.SHA == nil {
				continue
			}

			if *commit.SHA == lastSHA {
				slog.Info("Reached last analyzed commit, stopping fetch", "last_sha", lastSHA, "total_fetched", totalFetched, "sent", sentCount, "merge_skipped", mergeSkipped)
				return headSHA, nil
			}

			if len(commit.Parents) > 1 {
				mergeSkipped++
				continue
			}

			ref := codehost.CommitReference{SHA: *commit.SHA}
			if commit.Commit != nil && commit.Commit.Committer != nil && commit.Commit.Committer.Date != nil {
				ref.CommittedAt = *commit.Commit.Committer.Date
			}

			if headSHA == "" {
				headSHA = ref.SHA
			}
			sentCount++
			commits <- ref
		}

		if resp.NextPage == 0 {
			break
		}
		opts.ListOptions.Page = resp.NextPage
	}

	slog.Info("Commit fetch completed", "owner", owner, "repo", repoName, "total_fetched", totalFetched, "sent", sentCount, "merge_skipped", mergeSkipped, "head_sha", headSHA)
	return headSHA, nil
}

func (gc *GithubCodeHost) GetCommitDiff(ctx context.Context, r *repo.Repo, commitSHA string) (string, error) {
	owner, repoName, err := parseRepoURL(r.URL())
	if err != nil {
		return "", codehost.ErrInvalidRepoURL
	}

	slog.Debug("Fetching commit diff", "owner", owner, "repo", repoName, "commit_sha", commitSHA)

	commit, _, err := gc.client.Repositories.GetCommit(ctx, owner, repoName, commitSHA)
	if err != nil {
		slog.Error("Failed to fetch commit diff from GitHub", "owner", owner, "repo", repoName, "commit_sha", commitSHA, "error", err)
		return "", err
	}

	var diff string
	for _, file := range commit.Files {
		if file.Patch != nil {
			diff += fmt.Sprintf("File: %s\n%s\n\n", *file.Filename, *file.Patch)
		}
	}

	slog.Debug("Commit diff fetched", "commit_sha", commitSHA, "files_count", len(commit.Files), "diff_length", len(diff))
	return diff, nil
}

func parseRepoURL(repoURL string) (owner, repoName string, err error) {
	parsed, err := url.Parse(repoURL)
	if err != nil {
		return "", "", err
	}

	if parsed.Host != "github.com" {
		return "", "", fmt.Errorf("url '%s' is not github.com", repoURL)
	}

	parts := strings.Split(strings.TrimPrefix(parsed.Path, "/"), "/")
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("url '%s' has invalid path", repoURL)
	}

	return parts[0], strings.TrimSuffix(parts[1], ".git"), nil
}
