package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/octokerbs/chronocode/internal/app"
	"github.com/octokerbs/chronocode/internal/app/query"
	"github.com/octokerbs/chronocode/internal/domain/repo"
)

type ReposHandler struct {
	app app.Application
}

func NewReposHandler(app app.Application) *ReposHandler {
	return &ReposHandler{app: app}
}

type repoJSON struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	URL                 string `json:"url"`
	AddedAt             string `json:"addedAt"`
	LastAnalyzedCommit  string `json:"lastAnalyzedCommit"`
	CreatedAt           string `json:"createdAt"`
}

func (h *ReposHandler) List(w http.ResponseWriter, r *http.Request) {
	repos, err := h.app.Queries.GetRepos.Handle(r.Context(), query.GetRepos{})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"repositories": mapRepos(repos),
	})
}

func mapRepos(repos []*repo.Repo) []repoJSON {
	result := make([]repoJSON, len(repos))
	for i, r := range repos {
		result[i] = repoJSON{
			ID:                  formatInt64(r.ID()),
			Name:                r.Name(),
			URL:                 r.URL(),
			AddedAt:             r.CreatedAt().Format(time.RFC3339),
			LastAnalyzedCommit:  r.LastAnalyzedCommitSHA(),
			CreatedAt:           r.CreatedAt().Format(time.RFC3339),
		}
	}
	return result
}

func formatInt64(n int64) string {
	return strconv.FormatInt(n, 10)
}
