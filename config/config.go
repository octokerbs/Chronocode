package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL        string
	GeminiAPIKey       string
	GithubClientID     string
	GithubClientSecret string
	RedirectURL        string
	FrontendURL        string
}

func NewConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		GeminiAPIKey:       os.Getenv("GEMINI_API_KEY"),
		GithubClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		GithubClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:        os.Getenv("REDIRECT_URL"),
		FrontendURL:        os.Getenv("FRONTEND_URL"),
	}

	if cfg.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL is not set")
	}

	if cfg.GeminiAPIKey == "" {
		return nil, errors.New("GEMINI_API_KEY is not set")
	}

	if cfg.GithubClientID == "" {
		return nil, errors.New("GITHUB_CLIENT_ID is not set")
	}

	if cfg.GithubClientSecret == "" {
		return nil, errors.New("GITHUB_CLIENT_SECRET is not set")
	}

	if cfg.RedirectURL == "" {
		return nil, errors.New("REDIRECT_URL is not set")
	}

	if cfg.FrontendURL == "" {
		return nil, errors.New("FRONTEND_URL is not set")
	}

	return cfg, nil
}
