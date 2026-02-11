package adapters

import (
	"context"
	"errors"

	"github.com/google/go-github/github"

	"golang.org/x/oauth2"
)

type GithubCodeHostFactory struct {
}

func NewGithubCodeHostFactory() *GithubCodeHostFactory {
	return &GithubCodeHostFactory{}
}

func (f *GithubCodeHostFactory) Create(ctx context.Context, accessToken string) (*GithubCodeHost, error) {
	codeHost, err := NewGithubCodeHost(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	return codeHost, nil
}

type GithubCodeHost struct {
	client  *github.Client
	options *github.CommitsListOptions
}

func NewGithubCodeHost(ctx context.Context, accessToken string) (*GithubCodeHost, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	if client == nil {
		return nil, errors.New("missing codehost client")
	}

	options := &github.CommitsListOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	return &GithubCodeHost{client, options}, nil
}

// func (gc *GithubCodeHost) FetchRepository(ctx context.Context, repoURL string) (*analysis.Repository, error) {
// 	githubRepository, err := gc.fetchRepository(ctx, repoURL)
// 	if err != nil {
// 		return nil, err
// 	}

// 	now := time.Now()

// 	repository := &analysis.Repository{
// 		ID:                 *githubRepository.ID,
// 		CreatedAt:          &now,
// 		Name:               *githubRepository.FullName,
// 		URL:                *githubRepository.HTMLURL,
// 		LastAnalyzedCommit: "",
// 	}

// 	return repository, nil
// }

// func (gc *GithubCodeHost) FetchRepositoryID(ctx context.Context, repoURL string) (int64, error) {
// 	githubRepository, err := gc.fetchRepository(ctx, repoURL)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return *githubRepository.ID, nil
// }

// func (gc *GithubCodeHost) FetchCommit(ctx context.Context, repoURL string, commitSHA string) (*analysis.Commit, error) {
// 	githubCommit, err := gc.fetchCommit(ctx, repoURL, commitSHA)
// 	if err != nil {
// 		return nil, err
// 	}

// 	files := []string{}
// 	for _, file := range githubCommit.Files {
// 		files = append(files, *file.Filename)
// 	}

// 	repoID, err := gc.FetchRepositoryID(ctx, repoURL)
// 	if err != nil {
// 		return nil, err
// 	}

// 	now := time.Now()

// 	return &analysis.Commit{
// 		SHA:         commitSHA,
// 		CreatedAt:   &now,
// 		Author:      *githubCommit.Commit.Author.Name,
// 		Date:        githubCommit.Commit.Author.Date.Format(time.RFC3339),
// 		Message:     *githubCommit.Commit.Message,
// 		URL:         *githubCommit.HTMLURL,
// 		AuthorEmail: *githubCommit.Commit.Author.Email,
// 		Description: "",
// 		AuthorURL:   *githubCommit.Committer.HTMLURL,
// 		Files:       files,
// 		RepoID:      repoID,
// 	}, nil
// }

// func (gc *GithubCodeHost) FetchCommitDiff(ctx context.Context, repoURL string, commitSHA string) (string, error) {
// 	githubCommit, err := gc.fetchCommit(ctx, repoURL, commitSHA)
// 	if err != nil {
// 		return "", err
// 	}

// 	diff := ""
// 	if githubCommit.Files != nil {
// 		for _, file := range githubCommit.Files {
// 			if file.Patch != nil {
// 				diff += fmt.Sprintf("File: %s\n%s\n\n", *file.Filename, *file.Patch)
// 			}
// 		}
// 	}

// 	return diff, nil
// }

// func (gc *GithubCodeHost) ProduceCommitSHAs(ctx context.Context, repoURL string, lastAnalyzedCommitSHA string, commitSHAs chan<- string) (string, error) {
// 	repository, err := gc.fetchRepository(ctx, repoURL)
// 	if err != nil {
// 		return "", err
// 	}

// 	gc.setCommitOffset(lastAnalyzedCommitSHA)

// 	var newHeadSHA string
// 	for {
// 		pageCommits, resp, err := gc.fetchCommits(ctx, repository)
// 		if err != nil {
// 			return "", err
// 		}

// 		if newHeadSHA == "" && len(pageCommits) > 0 {
// 			newHeadSHA = *pageCommits[0].SHA
// 		}

// 		for _, commit := range pageCommits {
// 			commitSHAs <- *commit.SHA
// 		}

// 		if resp.NextPage == 0 {
// 			break
// 		}

// 		gc.nextPage(resp)
// 	}

// 	return newHeadSHA, nil
// }

// func (gc *GithubCodeHost) FetchAuthenticatedUser(ctx context.Context) (*user.GitHubProfile, error) {
// 	ghUser, _, err := gc.client.Users.Get(ctx, "")
// 	if err != nil {
// 		return nil, err
// 	}

// 	profile := &user.GitHubProfile{
// 		ID:    *ghUser.ID,
// 		Login: *ghUser.Login,
// 	}

// 	if ghUser.Name != nil {
// 		profile.Name = *ghUser.Name
// 	}
// 	if ghUser.AvatarURL != nil {
// 		profile.AvatarURL = *ghUser.AvatarURL
// 	}
// 	if ghUser.Email != nil {
// 		profile.Email = *ghUser.Email
// 	}

// 	return profile, nil
// }

// func (gc *GithubCodeHost) SearchUserRepositories(ctx context.Context, query string) ([]*analysis.Repository, error) {
// 	opts := &github.RepositoryListOptions{
// 		Sort:      "updated",
// 		Direction: "desc",
// 		ListOptions: github.ListOptions{
// 			PerPage: 20,
// 		},
// 	}

// 	ghRepos, _, err := gc.client.Repositories.List(ctx, "", opts)
// 	if err != nil {
// 		return nil, err
// 	}

// 	queryLower := strings.ToLower(query)
// 	var repos []*analysis.Repository
// 	for _, r := range ghRepos {
// 		if r.FullName == nil {
// 			continue
// 		}
// 		if query == "" || strings.Contains(strings.ToLower(*r.FullName), queryLower) {
// 			now := time.Now()
// 			repo := &analysis.Repository{
// 				ID:        *r.ID,
// 				CreatedAt: &now,
// 				Name:      *r.FullName,
// 			}
// 			if r.HTMLURL != nil {
// 				repo.URL = *r.HTMLURL
// 			}
// 			repos = append(repos, repo)
// 		}
// 	}

// 	return repos, nil
// }

// func (gc *GithubCodeHost) fetchRepository(ctx context.Context, repoURL string) (*github.Repository, error) {
// 	owner, repo, err := gc.parseRepoURL(repoURL)
// 	if err != nil {
// 		return nil, err
// 	}

// 	githubRepository, _, err := gc.client.Repositories.Get(ctx, owner, repo)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return githubRepository, nil
// }

// func (gc *GithubCodeHost) parseRepoURL(repoURL string) (string, string, error) {
// 	parsedURL, err := url.Parse(repoURL)
// 	if err != nil {
// 		return "", "", err
// 	}

// 	if parsedURL.Host != "github.com" {
// 		return "", "", fmt.Errorf("url '%s' is not github.com", repoURL)
// 	}

// 	pathParts := strings.Split(strings.TrimPrefix(parsedURL.Path, "/"), "/")
// 	if len(pathParts) < 2 {
// 		return "", "", fmt.Errorf("url '%s' has invalid path", repoURL)
// 	}

// 	repoName := strings.TrimSuffix(pathParts[1], ".git")

// 	return pathParts[0], repoName, nil
// }

// func (gc *GithubCodeHost) nextPage(resp *github.Response) {
// 	gc.options.ListOptions.Page = resp.NextPage
// }

// func (gc *GithubCodeHost) setCommitOffset(commitSHA string) {
// 	gc.options.SHA = commitSHA
// }

// func (gc *GithubCodeHost) fetchCommits(ctx context.Context, repository *github.Repository) ([]*github.RepositoryCommit, *github.Response, error) {
// 	pageCommits, resp, err := gc.client.Repositories.ListCommits(ctx, *repository.Owner.Login, *repository.Name, gc.options)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	return pageCommits, resp, nil
// }

// func (gc *GithubCodeHost) fetchCommit(ctx context.Context, repoURL string, commitSHA string) (*github.RepositoryCommit, error) {
// 	githubRepository, err := gc.fetchRepository(ctx, repoURL)
// 	if err != nil {
// 		return nil, err
// 	}

// 	githubCommit, _, err := gc.client.Repositories.GetCommit(ctx, *githubRepository.Owner.Login, *githubRepository.Name, commitSHA)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return githubCommit, nil
// }
