package http

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/octokerbs/chronocode/internal/app"
	"github.com/octokerbs/chronocode/internal/app/query"
	"github.com/octokerbs/chronocode/internal/domain/subcommit"
)

type SubcommitsHandler struct {
	app app.Application
}

func NewSubcommitsHandler(app app.Application) *SubcommitsHandler {
	return &SubcommitsHandler{app: app}
}

type subcommitJSON struct {
	ID          int64    `json:"id"`
	CreatedAt   string   `json:"createdAt"`
	Title       string   `json:"title"`
	Idea        string   `json:"idea"`
	Description string   `json:"description"`
	CommitSHA   string   `json:"commitSha"`
	Type        string   `json:"type"`
	Epic        string   `json:"epic"`
	Files       []string `json:"files"`
}

func (h *SubcommitsHandler) GetTimeline(w http.ResponseWriter, r *http.Request) {
	repoIDStr := r.URL.Query().Get("repo_id")
	repoID, err := strconv.ParseInt(repoIDStr, 10, 64)
	if err != nil {
		slog.Warn("Invalid repo_id in subcommits-timeline request", "repo_id_raw", repoIDStr, "error", err)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid repo_id"})
		return
	}

	slog.Info("Fetching subcommits timeline", "repo_id", repoID)

	token := AccessTokenFromContext(r.Context())
	result, err := h.app.Queries.GetSubcommits.Handle(r.Context(), query.GetSubcommits{
		RepoID:      repoID,
		AccessToken: token,
	})
	if err != nil {
		slog.Error("Failed to fetch subcommits timeline", "repo_id", repoID, "error", err)
		writeError(w, err)
		return
	}

	isAnalyzing := h.app.Locker.IsLocked(r.Context(), result.RepoURL)

	slog.Info("Subcommits timeline fetched", "repo_id", repoID, "count", len(result.Subcommits), "is_analyzing", isAnalyzing)

	writeJSON(w, http.StatusOK, map[string]any{
		"subcommits":  mapSubcommits(result.Subcommits),
		"repoId":      repoIDStr,
		"repoUrl":     result.RepoURL,
		"isAnalyzing": isAnalyzing,
	})
}

func mapSubcommits(scs []subcommit.Subcommit) []subcommitJSON {
	result := make([]subcommitJSON, len(scs))
	for i, sc := range scs {
		result[i] = subcommitJSON{
			ID:          sc.ID(),
			CreatedAt:   sc.CommittedAt().Format(time.RFC3339),
			Title:       sc.Title(),
			Idea:        sc.Idea(),
			Description: sc.Description(),
			CommitSHA:   sc.CommitSHA(),
			Type:        sc.ModificationType(),
			Epic:        sc.Epic(),
			Files:       sc.Files(),
		}
	}
	return result
}
