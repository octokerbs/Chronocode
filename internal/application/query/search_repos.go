package query

import (
	"context"
	"log/slog"

	"github.com/octokerbs/chronocode/internal/domain/codehost"
)

type SearchUserRepos struct {
	AccessToken string
	Query       string
}

type SearchUserReposHandler struct {
	codeHostFactory codehost.CodeHostFactory
}

func NewSearchUserReposHandler(codeHostFactory codehost.CodeHostFactory) SearchUserReposHandler {
	return SearchUserReposHandler{codeHostFactory: codeHostFactory}
}

func (h *SearchUserReposHandler) Handle(ctx context.Context, cmd SearchUserRepos) ([]codehost.RepoSearchResult, error) {
	slog.Info("SearchUserRepos query received", "query", cmd.Query)

	ch, err := h.codeHostFactory.Create(ctx, cmd.AccessToken)
	if err != nil {
		slog.Error("Failed to create code host client for repo search", "error", err)
		return nil, err
	}

	results, err := ch.SearchRepositories(ctx, cmd.Query)
	if err != nil {
		slog.Error("Failed to search repositories", "query", cmd.Query, "error", err)
		return nil, err
	}

	slog.Info("SearchUserRepos query completed", "query", cmd.Query, "results_count", len(results))
	return results, nil
}
