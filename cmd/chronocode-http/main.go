package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/octokerbs/chronocode-backend/internal/config"
	"github.com/octokerbs/chronocode-backend/internal/setup"
)

func main() {
	cfg, err := config.NewHTTPConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	httpApp, err := setup.NewHTTPApplication(cfg)
	if err != nil {
		log.Fatalf("Failed to bootstrap application: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := httpApp.Run(); err != nil {
			httpApp.Logger.Error("Server failed to start", err)
			stop()
		}
	}()

	<-ctx.Done() // Wait for shutdown signal

	httpApp.Logger.Info("Shutdown signal received, starting graceful shutdown...")
	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpApp.Shutdown(shutdownCtx); err != nil {
		httpApp.Logger.Error("Graceful shutdown failed", err)
	} else {
		httpApp.Logger.Info("Application shutdown complete.")
	}
}
