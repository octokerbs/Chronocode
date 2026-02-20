package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/octokerbs/chronocode/internal/app"
	"github.com/octokerbs/chronocode/internal/app/command"
)

type AnalyzeHandler struct {
	app app.Application
}

func NewAnalyzeHandler(app app.Application) *AnalyzeHandler {
	return &AnalyzeHandler{app: app}
}

func (h *AnalyzeHandler) Analyze(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RepoURL string `json:"repoUrl"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.Warn("Analyze request failed - invalid request body", "error", err)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	slog.Info("Starting repository analysis", "repo_url", body.RepoURL)

	token := AccessTokenFromContext(r.Context())
	repoID, err := h.app.Commands.AnalyzeRepo.HandleAsync(r.Context(), command.AnalyzeRepo{
		RepoURL:     body.RepoURL,
		AccessToken: token,
	})
	if err != nil {
		slog.Error("Repository analysis failed", "repo_url", body.RepoURL, "error", err)
		writeError(w, err)
		return
	}

	slog.Info("Repository analysis started", "repo_url", body.RepoURL, "repo_id", repoID)

	writeJSON(w, http.StatusOK, map[string]any{
		"message": "analysis started",
		"repoId":  repoID,
	})
}
