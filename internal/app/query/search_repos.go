package query

import (
	"context"

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
	ch, err := h.codeHostFactory.Create(ctx, cmd.AccessToken)
	if err != nil {
		return nil, err
	}

	return ch.SearchRepositories(ctx, cmd.Query)
}
