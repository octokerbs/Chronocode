package adapters

import (
	"context"
	"errors"
	"time"

	"github.com/octokerbs/chronocode/internal/domain/codehost"
	"github.com/octokerbs/chronocode/internal/domain/repo"
)

var (
	ValidRepoURL        = "https/validRepo"
	ValidRepoID   int64 = 123456789
	ValidRepoCommitSHA  = "CommitSHA-1"
	ValidRepoCommitDate = time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)

	ValidRepoCommitSHA2  = "CommitSHA-2"
	ValidRepoCommitDate2 = time.Date(2025, 1, 10, 10, 0, 0, 0, time.UTC)

	ValidEmptyRepoURL      = "https/emptyRepo"
	ValidEmptyRepoID int64 = 9876543221

	InvalidRepoURL   = "https/invalidRepo"
	ForbiddenRepoURL = "https/forbiddenRepo"

	FailingAgentRepoURL      = "https/failingAgentRepo"
	FailingAgentRepoID int64 = 111111111
	FailingCommitSHA         = "FailingCommitSHA-1"
	FailingCommitDate        = time.Date(2025, 1, 14, 10, 0, 0, 0, time.UTC)

	PartialFailureRepoURL      = "https/partialFailureRepo"
	PartialFailureRepoID int64 = 222222222

	ValidAccessToken   = "valid-token"
	InvalidAccessToken = "invalid-token"
	ValidCommitDiff    = "diff --git a/main.go b/main.go\n+func main() {}"
	FailingDiff        = "failing-diff"
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

	if url == PartialFailureRepoURL {
		return repo.NewRepo(PartialFailureRepoID, "partial-failure", PartialFailureRepoURL, ""), nil
	}

	return repo.NewRepo(ValidRepoID, "chronocode", ValidRepoURL, ""), nil
}

// commitsForRepo returns non-merge commits newest-first for the given repo.
func (c *CodeHost) commitsForRepo(r *repo.Repo) []codehost.CommitReference {
	switch r.URL() {
	case ValidEmptyRepoURL:
		return nil
	case FailingAgentRepoURL:
		return []codehost.CommitReference{
			{SHA: FailingCommitSHA, CommittedAt: FailingCommitDate},
		}
	case PartialFailureRepoURL:
		return []codehost.CommitReference{
			{SHA: ValidRepoCommitSHA, CommittedAt: ValidRepoCommitDate},
			{SHA: FailingCommitSHA, CommittedAt: FailingCommitDate},
		}
	default:
		return []codehost.CommitReference{
			{SHA: ValidRepoCommitSHA, CommittedAt: ValidRepoCommitDate},
			{SHA: ValidRepoCommitSHA2, CommittedAt: ValidRepoCommitDate2},
		}
	}
}

func (c *CodeHost) GetRepoCommitSHAsIntoChannel(ctx context.Context, r *repo.Repo, commits chan<- codehost.CommitReference) (string, error) {
	allCommits := c.commitsForRepo(r)
	lastSHA := r.LastAnalyzedCommitSHA()

	var headSHA string
	for _, ref := range allCommits {
		if ref.SHA == lastSHA {
			break
		}
		if headSHA == "" {
			headSHA = ref.SHA
		}
		commits <- ref
	}

	return headSHA, nil
}

func (c *CodeHost) GetCommitDiff(ctx context.Context, r *repo.Repo, commitSHA string) (string, error) {
	if commitSHA == FailingCommitSHA {
		return FailingDiff, nil
	}

	return ValidCommitDiff, nil
}
