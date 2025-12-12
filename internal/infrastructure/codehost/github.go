package codehost

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/octokerbs/chronocode-backend/internal/domain"
	pkg_errors "github.com/octokerbs/chronocode-backend/pkg/errs"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

type GithubClient struct {
	client  *github.Client
	options *github.CommitsListOptions
}

func NewGithubClient(ctx context.Context, accessToken string) *GithubClient {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	options := &github.CommitsListOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	return &GithubClient{client, options}
}

func (gc *GithubClient) FetchRepository(ctx context.Context, repoURL string) (*domain.Repository, error) {
	githubRepository, err := gc.fetchRepository(ctx, repoURL)
	if err != nil {
		return nil, gc.translateGithubError(err)
	}

	now := time.Now()

	repository := &domain.Repository{
		ID:                 *githubRepository.ID,
		CreatedAt:          &now,
		Name:               *githubRepository.FullName,
		URL:                *githubRepository.HTMLURL,
		LastAnalyzedCommit: "",
	}

	return repository, nil
}

func (gc *GithubClient) FetchRepositoryID(ctx context.Context, repoURL string) (int64, error) {
	githubRepository, err := gc.fetchRepository(ctx, repoURL)
	if err != nil {
		return 0, gc.translateGithubError(err)
	}

	return *githubRepository.ID, nil
}

func (gc *GithubClient) FetchCommit(ctx context.Context, repoURL string, commitSHA string) (*domain.Commit, error) {
	githubCommit, err := gc.fetchCommit(ctx, repoURL, commitSHA)
	if err != nil {
		return nil, gc.translateGithubError(err)
	}

	files := []string{}
	for _, file := range githubCommit.Files {
		files = append(files, *file.Filename)
	}

	repoID, err := gc.FetchRepositoryID(ctx, repoURL)
	if err != nil {
		return nil, gc.translateGithubError(err)
	}

	now := time.Now()

	return &domain.Commit{
		SHA:         commitSHA,
		CreatedAt:   &now,
		Author:      *githubCommit.Commit.Author.Name,
		Date:        githubCommit.Commit.Author.Date.Format(time.RFC3339),
		Message:     *githubCommit.Commit.Message,
		URL:         *githubCommit.HTMLURL,
		AuthorEmail: *githubCommit.Commit.Author.Email,
		Description: "",
		AuthorURL:   *githubCommit.Committer.HTMLURL,
		Files:       files,
		RepoID:      repoID,
	}, nil
}

func (gc *GithubClient) FetchCommitDiff(ctx context.Context, repoURL string, commitSHA string) (string, error) {
	githubCommit, err := gc.fetchCommit(ctx, repoURL, commitSHA)
	if err != nil {
		return "", gc.translateGithubError(err)
	}

	diff := ""
	if githubCommit.Files != nil {
		for _, file := range githubCommit.Files {
			if file.Patch != nil {
				diff += fmt.Sprintf("File: %s\n%s\n\n", *file.Filename, *file.Patch)
			}
		}
	}

	return diff, nil
}

func (gc *GithubClient) ProduceCommitSHAs(ctx context.Context, repoURL string, lastAnalyzedCommitSHA string, commitSHAs chan<- string) (string, error) {
	repository, err := gc.fetchRepository(ctx, repoURL)
	if err != nil {
		return "", gc.translateGithubError(err)
	}

	gc.setCommitOffset(lastAnalyzedCommitSHA)

	var newHeadSHA string
	for {
		pageCommits, resp, err := gc.fetchCommits(ctx, repository)
		if err != nil {
			return "", gc.translateGithubError(err)
		}

		if newHeadSHA == "" && len(pageCommits) > 0 {
			newHeadSHA = *pageCommits[0].SHA
		}

		for _, commit := range pageCommits {
			commitSHAs <- *commit.SHA
		}

		if resp.NextPage == 0 {
			break
		}

		gc.nextPage(resp)
	}

	return newHeadSHA, nil
}

func (gc *GithubClient) translateGithubError(err error) error {
	var githubErr *github.ErrorResponse
	if errors.As(err, &githubErr) {
		switch githubErr.Response.StatusCode {
		case 400:
			return pkg_errors.ErrBadRequest
		case 404:
			return pkg_errors.ErrNotFound
		case 401, 403:
			return pkg_errors.ErrUnauthorized
		}
	}

	return pkg_errors.NewError(pkg_errors.ErrInternalFailure, err)
}

func (gc *GithubClient) fetchRepository(ctx context.Context, repoURL string) (*github.Repository, error) {
	owner, repo, err := gc.parseRepoURL(repoURL)
	if err != nil {
		return nil, err
	}

	githubRepository, _, err := gc.client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return nil, err
	}

	return githubRepository, nil
}

func (gc *GithubClient) parseRepoURL(repoURL string) (string, string, error) {
	parsedURL, err := url.Parse(repoURL)
	if err != nil {
		return "", "", err
	}

	if parsedURL.Host != "github.com" {
		return "", "", fmt.Errorf("url '%s' is not github.com", repoURL)
	}

	pathParts := strings.Split(strings.TrimPrefix(parsedURL.Path, "/"), "/")
	if len(pathParts) < 2 {
		return "", "", fmt.Errorf("url '%s' has invalid path", repoURL)
	}

	repoName := strings.TrimSuffix(pathParts[1], ".git")

	return pathParts[0], repoName, nil
}

func (gc *GithubClient) nextPage(resp *github.Response) {
	gc.options.ListOptions.Page = resp.NextPage
}

func (gc *GithubClient) setCommitOffset(commitSHA string) {
	gc.options.SHA = commitSHA
}

func (gc *GithubClient) fetchCommits(ctx context.Context, repository *github.Repository) ([]*github.RepositoryCommit, *github.Response, error) {
	pageCommits, resp, err := gc.client.Repositories.ListCommits(ctx, *repository.Owner.Login, *repository.Name, gc.options)
	if err != nil {
		return nil, nil, err
	}
	return pageCommits, resp, nil
}

func (gc *GithubClient) fetchCommit(ctx context.Context, repoURL string, commitSHA string) (*github.RepositoryCommit, error) {
	githubRepository, err := gc.fetchRepository(ctx, repoURL)
	if err != nil {
		return nil, err
	}

	githubCommit, _, err := gc.client.Repositories.GetCommit(ctx, *githubRepository.Owner.Login, *githubRepository.Name, commitSHA)
	if err != nil {
		return nil, err
	}

	return githubCommit, nil
}

type GitHubAuth struct {
	config *oauth2.Config
}

func NewGitHubAuth(clientID, clientSecret, redirectURL string) *GitHubAuth {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"repo", "read:org", "user:email"},
		Endpoint:     githuboauth.Endpoint,
	}
	return &GitHubAuth{config: config}
}

func (a *GitHubAuth) GetAuthURL(state string) string {
	return a.config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

func (a *GitHubAuth) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return a.config.Exchange(ctx, code)
}
