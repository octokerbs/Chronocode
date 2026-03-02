package service

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/octokerbs/chronocode/internal/application"
	"github.com/octokerbs/chronocode/internal/application/command"
	"github.com/octokerbs/chronocode/internal/ports/http/utils"
)

type AnalyzeHandler struct {
	application application.Application
}

func NewAnalyzeHandler(application application.Application) *AnalyzeHandler {
	return &AnalyzeHandler{application: application}
}

func (h *AnalyzeHandler) Analyze(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RepoURL string `json:"repoUrl"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.Warn("Analyze request failed - invalid request body", "error", err)
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	slog.Info("Starting repository analysis", "repo_url", body.RepoURL)

	token := utils.AccessTokenFromContext(r.Context())
	repoID, err := h.application.Commands.AnalyzeRepo.HandleAsync(r.Context(), command.AnalyzeRepo{
		RepoURL:     body.RepoURL,
		AccessToken: token,
	})
	if err != nil {
		slog.Error("Repository analysis failed", "repo_url", body.RepoURL, "error", err)
		utils.WriteError(w, err)
		return
	}

	slog.Info("Repository analysis started", "repo_url", body.RepoURL, "repo_id", repoID)

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "analysis started",
		"repoId":  repoID,
	})
}
