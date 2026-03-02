package service

import (
	"log/slog"
	"net/http"

	"github.com/octokerbs/chronocode/internal/application"
	"github.com/octokerbs/chronocode/internal/application/query"
	"github.com/octokerbs/chronocode/internal/ports/http/utils"
)

type UserHandler struct {
	application application.Application
}

func NewUserHandler(application application.Application) *UserHandler {
	return &UserHandler{application: application}
}

func (h *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
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

func (h *UserHandler) SearchRepos(w http.ResponseWriter, r *http.Request) {
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
