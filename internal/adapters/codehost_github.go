package adapters

import (
	"context"
	"fmt"
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
		return nil, codehost.ErrAccessDenied
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return &GithubCodeHost{client: client}, nil
}

type GithubCodeHost struct {
	client *github.Client
}

func (gc *GithubCodeHost) CanAccessRepo(ctx context.Context, repoURL string) error {
	owner, repoName, err := parseRepoURL(repoURL)
	if err != nil {
		return codehost.ErrInvalidRepoURL
	}

	_, resp, err := gc.client.Repositories.Get(ctx, owner, repoName)
	if err != nil {
		if resp != nil && (resp.StatusCode == 404 || resp.StatusCode == 403 || resp.StatusCode == 401) {
			return codehost.ErrAccessDenied
		}
		return err
	}

	return nil
}

func (gc *GithubCodeHost) CreateRepoFromURL(ctx context.Context, repoURL string) (*repo.Repo, error) {
	owner, repoName, err := parseRepoURL(repoURL)
	if err != nil {
		return nil, codehost.ErrInvalidRepoURL
	}

	ghRepo, _, err := gc.client.Repositories.Get(ctx, owner, repoName)
	if err != nil {
		return nil, err
	}

	return repo.NewRepo(*ghRepo.ID, *ghRepo.FullName, repoURL, "", time.Now()), nil
}

func (gc *GithubCodeHost) GetAuthenticatedUser(ctx context.Context) (*codehost.UserProfile, error) {
	user, _, err := gc.client.Users.Get(ctx, "")
	if err != nil {
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
	return profile, nil
}

func (gc *GithubCodeHost) SearchRepositories(ctx context.Context, query string) ([]codehost.RepoSearchResult, error) {
	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 20},
		Sort:        "updated",
	}

	repos, _, err := gc.client.Repositories.List(ctx, "", opts)
	if err != nil {
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

	var headSHA string
	for {
		pageCommits, resp, err := gc.client.Repositories.ListCommits(ctx, owner, repoName, opts)
		if err != nil {
			return "", err
		}

		for _, commit := range pageCommits {
			if commit.SHA == nil {
				continue
			}

			if *commit.SHA == lastSHA {
				return headSHA, nil
			}

			if len(commit.Parents) > 1 {
				continue
			}

			ref := codehost.CommitReference{SHA: *commit.SHA}
			if commit.Commit != nil && commit.Commit.Committer != nil && commit.Commit.Committer.Date != nil {
				ref.CommittedAt = *commit.Commit.Committer.Date
			}

			if headSHA == "" {
				headSHA = ref.SHA
			}
			commits <- ref
		}

		if resp.NextPage == 0 {
			break
		}
		opts.ListOptions.Page = resp.NextPage
	}

	return headSHA, nil
}

func (gc *GithubCodeHost) GetCommitDiff(ctx context.Context, r *repo.Repo, commitSHA string) (string, error) {
	owner, repoName, err := parseRepoURL(r.URL())
	if err != nil {
		return "", codehost.ErrInvalidRepoURL
	}

	commit, _, err := gc.client.Repositories.GetCommit(ctx, owner, repoName, commitSHA)
	if err != nil {
		return "", err
	}

	var diff string
	for _, file := range commit.Files {
		if file.Patch != nil {
			diff += fmt.Sprintf("File: %s\n%s\n\n", *file.Filename, *file.Patch)
		}
	}

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
