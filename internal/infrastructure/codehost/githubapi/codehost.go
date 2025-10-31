package githubapi

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/octokerbs/chronocode-backend/internal/domain"
	"golang.org/x/oauth2"
)

type CodeHost struct {
	github *githubClient
}

func NewGithubCodeHost(ctx context.Context, accessToken string) *CodeHost {
	github := newGithubClient(ctx, accessToken)
	return &CodeHost{github}
}

func (ch *CodeHost) FetchRepository(ctx context.Context, repoURL string) (*domain.Repository, error) {
	githubRepository, err := ch.github.fetchRepository(ctx, repoURL)
	if err != nil {
		return nil, err
	}

	repository := &domain.Repository{
		ID:                 *githubRepository.ID,
		Name:               *githubRepository.FullName,
		URL:                *githubRepository.HTMLURL,
		LastAnalyzedCommit: "",
	}

	return repository, nil
}

func (ch *CodeHost) FetchRepositoryID(ctx context.Context, repoURL string) (int64, error) {
	githubRepository, err := ch.github.fetchRepository(ctx, repoURL)
	if err != nil {
		return 0, err
	}

	return *githubRepository.ID, nil
}

func (ch *CodeHost) FetchCommit(ctx context.Context, repoURL string, commitSHA string) (*domain.Commit, error) {
	githubCommit, err := ch.github.fetchCommit(ctx, repoURL, commitSHA)
	if err != nil {
		return nil, err
	}

	files := []string{}
	for _, file := range githubCommit.Files {
		files = append(files, *file.Filename)
	}

	repoID, err := ch.FetchRepositoryID(ctx, repoURL)
	if err != nil {
		return nil, err
	}

	return &domain.Commit{
		SHA:         commitSHA,
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

func (ch *CodeHost) FetchCommitDiff(ctx context.Context, repoURL string, commitSHA string) (string, error) {
	githubCommit, err := ch.github.fetchCommit(ctx, repoURL, commitSHA)
	if err != nil {
		return "", err
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

func (ch *CodeHost) ProduceCommitSHAs(ctx context.Context, repoURL string, lastAnalyzedCommitSHA string, commits chan<- string) {
	repository, err := ch.github.fetchRepository(ctx, repoURL)
	if err != nil {
		return
	}

	ch.github.setCommitOffset(lastAnalyzedCommitSHA)

	for {
		pageCommits, resp, err := ch.github.fetchCommits(ctx, repository)
		if err != nil {
			return
		}

		for _, commit := range pageCommits {
			commits <- *commit.SHA
		}

		if resp.NextPage == 0 {
			break
		}

		ch.github.nextPage(resp)
	}

	close(commits)
}

type githubClient struct {
	client  *github.Client
	options *github.CommitsListOptions
}

func newGithubClient(ctx context.Context, accessToken string) *githubClient {
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

	return &githubClient{client, options}
}

func (gc *githubClient) fetchRepository(ctx context.Context, repoURL string) (*github.Repository, error) {
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

func (gc *githubClient) parseRepoURL(repoURL string) (string, string, error) {
	parsedURL, err := url.Parse(repoURL)
	if err != nil {
		return "", "", err
	}

	if parsedURL.Host != "github.com" {
		return "", "", fmt.Errorf("not supported version control repository")
	}

	pathParts := strings.Split(strings.TrimPrefix(parsedURL.Path, "/"), "/")
	if len(pathParts) < 2 {
		return "", "", fmt.Errorf("invalid GitHub repository URL format")
	}

	repoName := strings.TrimSuffix(pathParts[1], ".git")

	return pathParts[0], repoName, nil
}

func (gc *githubClient) nextPage(resp *github.Response) {
	gc.options.ListOptions.Page = resp.NextPage
}

func (gc *githubClient) setCommitOffset(commitSHA string) {
	gc.options.SHA = commitSHA
}

func (gc *githubClient) fetchCommits(ctx context.Context, repository *github.Repository) ([]*github.RepositoryCommit, *github.Response, error) {
	pageCommits, resp, err := gc.client.Repositories.ListCommits(ctx, *repository.Owner.Login, *repository.Name, gc.options)
	if err != nil {
		return nil, nil, err
	}
	return pageCommits, resp, nil
}

func (gc *githubClient) fetchCommit(ctx context.Context, repoURL string, commitSHA string) (*github.RepositoryCommit, error) {
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
