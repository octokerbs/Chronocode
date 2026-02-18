package adapters

import (
	"context"
	"errors"

	"github.com/octokerbs/chronocode/internal/domain/codehost"
	"github.com/octokerbs/chronocode/internal/domain/repo"
)

var (
	ValidRepoURL             = "https/validRepo"
	ValidRepoID        int64 = 123456789
	ValidRepoCommitSHA       = "CommitSHA-1"
	ValidEmptyRepoURL        = "https/emptyRepo"
	ValidEmptyRepoID   int64 = 9876543221
	InvalidRepoURL           = "https/invalidRepo"
	ForbiddenRepoURL         = "https/forbiddenRepo"
	FailingAgentRepoURL        = "https/failingAgentRepo"
	FailingAgentRepoID   int64 = 111111111
	ValidAccessToken           = "valid-token"
	InvalidAccessToken         = "invalid-token"
	ValidCommitDiff            = "diff --git a/main.go b/main.go\n+func main() {}"
	FailingDiff                = "failing-diff"
)

type CodeHostFactory struct{}

func NewCodeHostFactory() *CodeHostFactory {
	return &CodeHostFactory{}
}

func (f *CodeHostFactory) Create(ctx context.Context, accessToken string) (codehost.CodeHost, error) {
	if accessToken == "" || accessToken == InvalidAccessToken {
		return nil, errors.New("invalid access token")
	}

	return NewCodeHost(), nil
}

type CodeHost struct{}

func NewCodeHost() *CodeHost {
	return &CodeHost{}
}

func (c *CodeHost) CanAccessRepo(ctx context.Context, repoURL string) error {
	if repoURL == ForbiddenRepoURL {
		return codehost.ErrAccessDenied
	}

	return nil
}

func (c *CodeHost) CreateRepoFromURL(ctx context.Context, url string) (*repo.Repo, error) {
	if url == InvalidRepoURL {
		return nil, codehost.ErrInvalidRepoURL
	}

	if url == ValidEmptyRepoURL {
		return repo.NewRepo(ValidEmptyRepoID, "empty-repo", ValidEmptyRepoURL, ""), nil
	}

	if url == FailingAgentRepoURL {
		return repo.NewRepo(FailingAgentRepoID, "failing-agent", FailingAgentRepoURL, ""), nil
	}

	return repo.NewRepo(ValidRepoID, "chronocode", ValidRepoURL, "FFFFFF"), nil
}

func (c *CodeHost) GetRepoCommitSHAsIntoChannel(ctx context.Context, repo *repo.Repo, commitSHAs chan<- string) error {
	if repo.URL() == ValidEmptyRepoURL {
		return nil
	}

	commitSHAs <- ValidRepoCommitSHA
	return nil
}

func (c *CodeHost) GetCommitDiff(ctx context.Context, r *repo.Repo, commitSHA string) (string, error) {
	if r.URL() == FailingAgentRepoURL {
		return FailingDiff, nil
	}

	return ValidCommitDiff, nil
}
