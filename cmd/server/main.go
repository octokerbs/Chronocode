package main

import (
	"context"
	"log"
	"os"

	httpport "github.com/octokerbs/chronocode/internal/ports/http"
	"github.com/octokerbs/chronocode/internal/service"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func main() {
	ctx := context.Background()
	application := service.NewApplication(ctx)

	oauthConfig := &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
		Scopes:       []string{"read:user", "user:email", "repo"},
		Endpoint:     github.Endpoint,
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := httpport.NewServer(application, oauthConfig, frontendURL, port)

	log.Printf("Starting server on :%s", port)
	log.Fatal(server.ListenAndServe())
}
