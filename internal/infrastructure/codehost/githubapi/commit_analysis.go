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

type GithubClient struct {
	client  *github.Client
	options *github.CommitsListOptions
}

func NewGithubClient(ctx context.Context, accessToken string) (*GithubClient, error) {
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

	return &GithubClient{client, options}, nil
}

func getRepository(ctx context.Context, client *github.Client, repoURL string) (*github.Repository, error) {
	owner, repo, err := parseRepoURL(repoURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing repository URL: %v", err.Error())
	}

	githubRepository, _, err := client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("error fetching repository: %v", err.Error())
	}

	return githubRepository, nil
}

func parseRepoURL(repoURL string) (string, string, error) {
	parsedURL, err := url.Parse(repoURL)
	if err != nil {
		return "", "", err
	}

	if parsedURL.Host != "github.com" {
		return "", "", fmt.Errorf("not supported version control repository")
	}

	// Remove leading slash and split path
	pathParts := strings.Split(strings.TrimPrefix(parsedURL.Path, "/"), "/")
	if len(pathParts) < 2 {
		return "", "", fmt.Errorf("invalid GitHub repository URL format")
	}

	// Handle potential .git suffix
	repoName := strings.TrimSuffix(pathParts[1], ".git")

	return pathParts[0], repoName, nil
}

func (gc *GithubClient) ProduceCommits(ctx context.Context, repoURL string, lastAnalyzedCommitSHA string, commits chan<- string, errors chan<- string) {
	repository, err := getRepository(ctx, gc.client, repoURL)
	if err != nil {
		errors <- fmt.Sprintf("error finding repository: %v", err.Error())
		return
	}

	gc.options.SHA = lastAnalyzedCommitSHA

	for {
		pageCommits, resp, err := gc.client.Repositories.ListCommits(ctx, *repository.Owner.Login, *repository.Name, gc.options)
		if err != nil {
			errors <- fmt.Sprintf("error fetching commits: %v", err.Error())
			return
		}

		for _, commit := range pageCommits {
			commits <- *commit.SHA
		}

		if resp.NextPage == 0 {
			break
		}

		gc.options.ListOptions.Page = resp.NextPage
	}
}

func (gc *GithubClient) GetCommitDiff(ctx context.Context, repoURL string, commitSHA string) (string, error) {
	repository, err := getRepository(ctx, gc.client, repoURL)
	if err != nil {
		return "", err
	}

	commitFullData, _, err := gc.client.Repositories.GetCommit(ctx, *repository.Owner.Login, *repository.Name, commitSHA) // Github api doesn't fetch the file data when fetching repo commits. We have to do it ourselves.
	if err != nil {
		return "", fmt.Errorf("error fetching commit %s: %v", commitSHA, err)
	}

	diff := ""
	if commitFullData.Files != nil {
		for _, file := range commitFullData.Files {
			if file.Patch != nil {
				diff += fmt.Sprintf("File: %s\n%s\n\n", *file.Filename, *file.Patch)
			}
		}
	}

	return diff, nil
}

func (gc *GithubClient) RepositoryID(ctx context.Context, repoURL string) (int64, error) {
	repository, err := getRepository(ctx, gc.client, repoURL)
	if err != nil {
		return 0, err
	}

	return *repository.ID, nil
}

func (gc *GithubClient) NewRepository(ctx context.Context, repoURL string) (*domain.Repository, error) {
	repository, err := getRepository(ctx, gc.client, repoURL)
	if err != nil {
		return nil, err
	}

	return &domain.Repository{
		ID:                 *repository.ID,
		Name:               *repository.FullName,
		URL:                *repository.HTMLURL,
		LastAnalyzedCommit: "",
	}, nil
}

func (gc *GithubClient) NewCommit(ctx context.Context, repoURL string, commitSHA string) (*domain.Commit, error) {
	repository, err := getRepository(ctx, gc.client, repoURL)
	if err != nil {
		return nil, err
	}

	commit, _, _ := gc.client.Repositories.GetCommit(ctx, *repository.Owner.Login, *repository.Name, commitSHA)

	files := []string{}
	for _, file := range commit.Files {
		files = append(files, *file.Filename)
	}

	return &domain.Commit{
		SHA:         commitSHA,
		Author:      *commit.Commit.Author.Name,
		Date:        commit.Commit.Author.Date.Format(time.RFC3339),
		Message:     *commit.Commit.Message,
		URL:         *commit.HTMLURL,
		AuthorEmail: *commit.Commit.Author.Email,
		Description: "",
		AuthorURL:   *commit.Committer.HTMLURL,
		Files:       files,
		RepoID:      *repository.ID,
	}, nil
}
