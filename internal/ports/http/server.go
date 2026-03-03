package http

import (
	"log/slog"
	"net/http"

	"github.com/octokerbs/chronocode/internal/application"
	"github.com/octokerbs/chronocode/internal/ports/http/auth"
	"github.com/octokerbs/chronocode/internal/ports/http/utils"
	"golang.org/x/oauth2"
)

func NewServer(application application.Application, oauthConfig *oauth2.Config, frontendURL string, port string) *http.Server {
	mux := http.NewServeMux()

	authHandler := auth.NewHandler(oauthConfig, frontendURL)
	applicationHandler := NewApplicationHandler(application)

	// Public routes
	mux.HandleFunc("GET /auth/status", authHandler.Status)
	mux.HandleFunc("GET /auth/github/login", authHandler.Login)
	mux.HandleFunc("GET /auth/github/callback", authHandler.Callback)
	mux.HandleFunc("POST /auth/logout", authHandler.Logout)

	// Protected routes
	protected := http.NewServeMux()
	protected.HandleFunc("GET /user/profile", applicationHandler.GetUserProfileQuery)
	protected.HandleFunc("GET /user/repos/search", applicationHandler.SearchReposQuery)
	protected.HandleFunc("GET /repositories", applicationHandler.GetReposQuery)
	protected.HandleFunc("POST /analyze", applicationHandler.AnalyzeRepoCommand)
	protected.HandleFunc("GET /subcommits-timeline", applicationHandler.GetSubcommitsQuery)

	mux.Handle("/", utils.AuthMiddleware(protected))

	handler := utils.RequestLoggingMiddleware(utils.CORSMiddleware(frontendURL)(mux))

	slog.Info("HTTP server configured", "port", port, "frontend_url", frontendURL, "routes", []string{
		"GET /auth/status", "GET /auth/github/login", "GET /auth/github/callback", "POST /auth/logout",
		"GET /user/profile", "GET /user/repos/search", "GET /repositories", "POST /analyze", "GET /subcommits-timeline",
	})

	return &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}
}
