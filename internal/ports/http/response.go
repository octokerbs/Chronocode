package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/octokerbs/chronocode/internal/domain/analysis"
	"github.com/octokerbs/chronocode/internal/domain/codehost"
	"github.com/octokerbs/chronocode/internal/domain/repo"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, err error) {
	status, message := mapDomainError(err)
	writeJSON(w, status, map[string]string{"error": message})
}

func mapDomainError(err error) (int, string) {
	switch {
	case errors.Is(err, codehost.ErrAccessDenied):
		return http.StatusForbidden, "access denied"
	case errors.Is(err, codehost.ErrInvalidRepoURL):
		return http.StatusBadRequest, "invalid repository URL"
	case errors.Is(err, repo.ErrRepositoryNotFound):
		return http.StatusNotFound, "repository not found"
	case errors.Is(err, analysis.ErrAnalysisInProgress):
		return http.StatusConflict, "analysis already in progress"
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
