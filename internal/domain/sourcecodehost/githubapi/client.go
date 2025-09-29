package githubapi

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GithubService struct {
	client     *github.Client
	repository *github.Repository
	options    *github.CommitsListOptions
}

func NewGithubClient(ctx context.Context, accessToken string, repoURL string) (*GithubService, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	repository, err := getRepository(ctx, client, repoURL)
	if err != nil {
		return nil, err
	}

	options := &github.CommitsListOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	return &GithubService{client, repository, options}, nil
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

func (g *GithubService) ProduceCommits(ctx context.Context, lastAnalyzedCommitSHA string, commits chan<- string, errors chan<- error) {
	g.options.SHA = lastAnalyzedCommitSHA

	for {
		pageCommits, resp, err := g.client.Repositories.ListCommits(ctx, *g.repository.Owner.Login, *g.repository.Name, g.options)
		if err != nil {
			errors <- fmt.Errorf("error fetching commits: %v", err.Error())
			return
		}

		for _, commit := range pageCommits {
			commits <- *commit.SHA
		}

		if resp.NextPage == 0 {
			break
		}

		g.options.ListOptions.Page = resp.NextPage
	}
}

func (g *GithubService) GetCommitDiff(ctx context.Context, commitSHA string) (string, error) {
	commitFullData, _, err := g.client.Repositories.GetCommit(ctx, *g.repository.Owner.Login, *g.repository.Name, commitSHA) // Github api doesn't fetch the file data when fetching repo commits. We have to do it ourselves.
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

func (g *GithubService) RepositoryID() int64 {
	return *g.repository.ID
}

func (g *GithubService) GetRepositoryData() map[string]interface{} {
	return map[string]interface{}{
		"id":   *g.repository.ID,
		"name": *g.repository.FullName,
		"url":  *g.repository.HTMLURL,
	}
}

func (g *GithubService) GetCommitData(ctx context.Context, commitSHA string) map[string]interface{} {
	commit, _, _ := g.client.Repositories.GetCommit(ctx, *g.repository.Owner.Login, *g.repository.Name, commitSHA)

	files := []string{}
	for _, file := range commit.Files {
		files = append(files, *file.Filename)
	}

	return map[string]interface{}{
		"author":        *commit.Commit.Author.Name,
		"author_email":  *commit.Commit.Author.Email,
		"author_url":    "",
		"date":          commit.Commit.Author.Date.Format(time.RFC3339),
		"message":       *commit.Commit.Message,
		"url":           *commit.HTMLURL,
		"files":         files,
		"repository_id": *g.repository.ID,
	}
}
