package main

import (
	"context"

	"github.com/octokerbs/chronocode-backend/common/logs"
	"github.com/octokerbs/chronocode-backend/internal/ports"
	"github.com/octokerbs/chronocode-backend/internal/service"
)

func main() {
	logger := logs.Init()
	ctx := context.Background()

	app := service.NewApplication(ctx, logger)
	server := ports.NewHttpServer(app)
	server.RunHTTPServer()
}
