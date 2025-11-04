package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// HTTPConfig holds configuration for the HTTP application.
// It loads settings from environment variables.
// Required environment variables:
// - DATABASE_URL: the database connection URL
// - GEMINI_API_KEY: the API key for the Gemini agent
type HTTPConfig struct {
	Port         string
	DatabaseURL  string
	GeminiAPIKey string
}

func NewHTTPConfig() (*HTTPConfig, error) {
	_ = godotenv.Load()

	cfg := &HTTPConfig{
		Port:         getEnv("PORT", ":8080"),
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),
	}

	if cfg.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL is not set")
	}
	if cfg.GeminiAPIKey == "" {
		return nil, errors.New("GEMINI_API_KEY is not set")
	}

	return cfg, nil
}
