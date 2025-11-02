package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	DatabaseURL  string
	GeminiAPIKey string
}

func New() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
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

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
