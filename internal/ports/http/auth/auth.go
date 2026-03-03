package auth

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"time"

	"github.com/octokerbs/chronocode/internal/ports/http/utils"
	"golang.org/x/oauth2"
)

type Handler struct {
	oauthConfig *oauth2.Config
	frontendURL string
}

func NewHandler(oauthConfig *oauth2.Config, frontendURL string) *Handler {
	return &Handler{oauthConfig: oauthConfig, frontendURL: frontendURL}
}

func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("access_token")
	isLoggedIn := err == nil
	slog.Info("Auth status check", "is_logged_in", isLoggedIn, "remote_addr", r.RemoteAddr)
	utils.WriteJSON(w, http.StatusOK, map[string]bool{"isLoggedIn": isLoggedIn})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	state := generateState()
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   300,
	})

	url := h.oauthConfig.AuthCodeURL(state)
	slog.Info("OAuth login initiated, redirecting to GitHub", "remote_addr", r.RemoteAddr)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) Callback(w http.ResponseWriter, r *http.Request) {
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil || stateCookie.Value != r.URL.Query().Get("state") {
		slog.Warn("OAuth callback failed - invalid state parameter", "has_cookie", err == nil, "remote_addr", r.RemoteAddr)
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid state"})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "oauth_state",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	code := r.URL.Query().Get("code")
	slog.Info("OAuth callback received, exchanging code for token", "remote_addr", r.RemoteAddr)

	token, err := h.oauthConfig.Exchange(r.Context(), code)
	if err != nil {
		slog.Error("OAuth token exchange failed", "error", err, "remote_addr", r.RemoteAddr)
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "token exchange failed"})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token.AccessToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		MaxAge:   int(24 * time.Hour / time.Second),
	})

	slog.Info("OAuth login successful, redirecting to frontend", "redirect_url", h.frontendURL+"/home", "remote_addr", r.RemoteAddr)
	http.Redirect(w, r, h.frontendURL+"/home", http.StatusTemporaryRedirect)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "access_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	slog.Info("User logged out", "remote_addr", r.RemoteAddr)
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}

func generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
