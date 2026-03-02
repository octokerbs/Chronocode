package service

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/octokerbs/chronocode/internal/application"
	"github.com/octokerbs/chronocode/internal/application/query"
	"github.com/octokerbs/chronocode/internal/domain/repo"
	"github.com/octokerbs/chronocode/internal/ports/http/utils"
)

type ReposHandler struct {
	application application.Application
}

func NewReposHandler(application application.Application) *ReposHandler {
	return &ReposHandler{application: application}
}

type repoJSON struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	URL     string `json:"url"`
	AddedAt string `json:"addedAt"`
}

func (h *ReposHandler) List(w http.ResponseWriter, r *http.Request) {
	slog.Info("Listing all repositories")

	repos, err := h.application.Queries.GetRepos.Handle(r.Context(), query.GetRepos{})
	if err != nil {
		slog.Error("Failed to list repositories", "error", err)
		utils.WriteError(w, err)
		return
	}

	slog.Info("Repositories listed", "count", len(repos))

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"repositories": mapRepos(repos),
	})
}

func mapRepos(repos []*repo.Repo) []repoJSON {
	result := make([]repoJSON, len(repos))
	for i, r := range repos {
		result[i] = repoJSON{
			ID:      formatInt64(r.ID()),
			Name:    r.Name(),
			URL:     r.URL(),
			AddedAt: r.CreatedAt().Format(time.RFC3339),
		}
	}
	return result
}

func formatInt64(n int64) string {
	return strconv.FormatInt(n, 10)
}
