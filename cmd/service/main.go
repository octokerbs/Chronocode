package main

import (
	"context"
	"net/http"

	"github.com/octokerbs/chronocode-backend/common/logs"
	"github.com/octokerbs/chronocode-backend/common/server"
	"github.com/octokerbs/chronocode-backend/internal/ports"
)

func main() {
	logger := logs.Init()
	ctx := context.Background()

	app := NewApplication(ctx, logger)

	server.RunHTTPServer(func(mux *http.ServeMux) http.Handler {
		httpServer := ports.NewHttpServer(app)
		mux.HandleFunc("/subcommits", httpServer.GetSubcommits)
		return mux
	}, logger)
}
