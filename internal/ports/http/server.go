package http

import (
	"net/http"

	"github.com/octokerbs/chronocode/internal/app"
	"golang.org/x/oauth2"
)

func NewServer(application app.Application, oauthConfig *oauth2.Config, frontendURL string, port string) *http.Server {
	mux := http.NewServeMux()

	authHandler := NewAuthHandler(oauthConfig, frontendURL)
	analyzeHandler := NewAnalyzeHandler(application)
	subcommitsHandler := NewSubcommitsHandler(application)
	reposHandler := NewReposHandler(application)
	userHandler := NewUserHandler(application)

	// Public routes
	mux.HandleFunc("GET /auth/status", authHandler.Status)
	mux.HandleFunc("GET /auth/github/login", authHandler.Login)
	mux.HandleFunc("GET /auth/github/callback", authHandler.Callback)
	mux.HandleFunc("POST /auth/logout", authHandler.Logout)

	// Protected routes
	protected := http.NewServeMux()
	protected.HandleFunc("GET /user/profile", userHandler.Profile)
	protected.HandleFunc("GET /user/repos/search", userHandler.SearchRepos)
	protected.HandleFunc("GET /repositories", reposHandler.List)
	protected.HandleFunc("POST /analyze", analyzeHandler.Analyze)
	protected.HandleFunc("GET /subcommits-timeline", subcommitsHandler.GetTimeline)

	mux.Handle("/", AuthMiddleware(protected))

	handler := CORSMiddleware(frontendURL)(mux)

	return &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}
}
