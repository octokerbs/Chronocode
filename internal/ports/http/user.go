package http

import (
	"log/slog"
	"net/http"

	"github.com/octokerbs/chronocode/internal/app"
	"github.com/octokerbs/chronocode/internal/app/query"
)

type UserHandler struct {
	app app.Application
}

func NewUserHandler(app app.Application) *UserHandler {
	return &UserHandler{app: app}
}

func (h *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	slog.Info("Fetching user profile")

	token := AccessTokenFromContext(r.Context())
	profile, err := h.app.Queries.GetUserProfile.Handle(r.Context(), query.GetUserProfile{
		AccessToken: token,
	})
	if err != nil {
		slog.Error("Failed to fetch user profile", "error", err)
		writeError(w, err)
		return
	}

	slog.Info("User profile fetched", "user_login", profile.Login, "user_id", profile.ID)

	writeJSON(w, http.StatusOK, map[string]any{
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

	token := AccessTokenFromContext(r.Context())
	results, err := h.app.Queries.SearchUserRepos.Handle(r.Context(), query.SearchUserRepos{
		AccessToken: token,
		Query:       q,
	})
	if err != nil {
		slog.Error("Failed to search user repositories", "query", q, "error", err)
		writeError(w, err)
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

	writeJSON(w, http.StatusOK, map[string]any{
		"repositories": repos,
	})
}
