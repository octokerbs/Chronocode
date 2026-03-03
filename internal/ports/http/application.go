package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/octokerbs/chronocode/internal/application"
	"github.com/octokerbs/chronocode/internal/application/command"
	"github.com/octokerbs/chronocode/internal/application/query"
	"github.com/octokerbs/chronocode/internal/domain/analysis"
	"github.com/octokerbs/chronocode/internal/ports/http/utils"
)

type ApplicationHandler struct {
	application application.Application
}

func NewApplicationHandler(application application.Application) *ApplicationHandler {
	return &ApplicationHandler{application: application}
}

func (h *ApplicationHandler) AnalyzeRepoCommand(w http.ResponseWriter, r *http.Request) {
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
		if errors.Is(err, analysis.ErrAnalysisInProgress) && repoID != 0 {
			slog.Info("Analysis already in progress, returning existing repo", "repo_url", body.RepoURL, "repo_id", repoID)
			utils.WriteJSON(w, http.StatusOK, map[string]any{
				"message": "analysis already in progress",
				"repoId":  repoID,
			})
			return
		}
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

func (h *ApplicationHandler) GetReposQuery(w http.ResponseWriter, r *http.Request) {
	slog.Info("Listing all repositories")

	repos, err := h.application.Queries.GetRepos.Handle(r.Context(), query.GetRepos{})
	if err != nil {
		slog.Error("Failed to list repositories", "error", err)
		utils.WriteError(w, err)
		return
	}

	slog.Info("Repositories listed", "count", len(repos))

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"repositories": utils.MapRepos(repos),
	})
}

func (h *ApplicationHandler) GetSubcommitsQuery(w http.ResponseWriter, r *http.Request) {
	repoIDStr := r.URL.Query().Get("repo_id")
	repoID, err := strconv.ParseInt(repoIDStr, 10, 64)
	if err != nil {
		slog.Warn("Invalid repo_id in subcommits-timeline request", "repo_id_raw", repoIDStr, "error", err)
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid repo_id"})
		return
	}

	slog.Info("Fetching subcommits timeline", "repo_id", repoID)

	token := utils.AccessTokenFromContext(r.Context())
	result, err := h.application.Queries.GetSubcommits.Handle(r.Context(), query.GetSubcommits{
		RepoID:      repoID,
		AccessToken: token,
	})
	if err != nil {
		slog.Error("Failed to fetch subcommits timeline", "repo_id", repoID, "error", err)
		utils.WriteError(w, err)
		return
	}

	isAnalyzing := h.application.Locker.IsLocked(r.Context(), result.RepoURL)

	slog.Info("Subcommits timeline fetched", "repo_id", repoID, "count", len(result.Subcommits), "is_analyzing", isAnalyzing)

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"subcommits":  utils.MapSubcommits(result.Subcommits),
		"repoId":      repoIDStr,
		"repoUrl":     result.RepoURL,
		"isAnalyzing": isAnalyzing,
	})
}

func (h *ApplicationHandler) GetUserProfileQuery(w http.ResponseWriter, r *http.Request) {
	slog.Info("Fetching user profile")

	token := utils.AccessTokenFromContext(r.Context())
	profile, err := h.application.Queries.GetUserProfile.Handle(r.Context(), query.GetUserProfile{
		AccessToken: token,
	})
	if err != nil {
		slog.Error("Failed to fetch user profile", "error", err)
		utils.WriteError(w, err)
		return
	}

	slog.Info("User profile fetched", "user_login", profile.Login, "user_id", profile.ID)

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"id":        profile.ID,
		"login":     profile.Login,
		"name":      profile.Name,
		"avatarUrl": profile.AvatarURL,
		"email":     profile.Email,
	})
}

func (h *ApplicationHandler) SearchReposQuery(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	slog.Info("Searching user repositories", "query", q)

	token := utils.AccessTokenFromContext(r.Context())
	results, err := h.application.Queries.SearchUserRepos.Handle(r.Context(), query.SearchUserRepos{
		AccessToken: token,
		Query:       q,
	})
	if err != nil {
		slog.Error("Failed to search user repositories", "query", q, "error", err)
		utils.WriteError(w, err)
		return
	}

	slog.Info("User repositories search completed", "query", q, "results_count", len(results))

	repos := make([]map[string]any, len(results))
	for i, r := range results {
		repos[i] = map[string]any{
			"id":   r.ID,
			"name": r.Name,
			"url":  r.URL,
		}
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"repositories": repos,
	})
}
