package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	httpApp := NewHTTPApplication()

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
