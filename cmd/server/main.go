package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	httpport "github.com/octokerbs/chronocode/internal/ports/http"
	"github.com/octokerbs/chronocode/internal/service"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func main() {
	logLevel := slog.LevelInfo
	if os.Getenv("LOG_LEVEL") == "debug" {
		logLevel = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})))

	slog.Info("Chronocode server starting", "log_level", logLevel.String())

	ctx := context.Background()
	application := service.NewApplication(ctx)

	oauthConfig := &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
		Scopes:       []string{"read:user", "user:email", "repo"},
		Endpoint:     github.Endpoint,
	}

	slog.Info("GitHub OAuth configured",
		"client_id", os.Getenv("GITHUB_CLIENT_ID"),
		"redirect_url", os.Getenv("GITHUB_REDIRECT_URL"),
	)

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := httpport.NewServer(application, oauthConfig, frontendURL, port)

	slog.Info("Chronocode server ready", "port", port, "frontend_url", frontendURL)
	log.Fatal(server.ListenAndServe())
}
