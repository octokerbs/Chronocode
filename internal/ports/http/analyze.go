package http

import (
	"encoding/json"
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
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	token := AccessTokenFromContext(r.Context())
	repoID, err := h.app.Commands.AnalyzeRepo.Handle(r.Context(), command.AnalyzeRepo{
		RepoURL:     body.RepoURL,
		AccessToken: token,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"message": "analysis complete",
		"repoId":  repoID,
	})
}
