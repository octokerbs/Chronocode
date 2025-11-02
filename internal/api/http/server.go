package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/octokerbs/chronocode-backend/internal/application"
	"github.com/octokerbs/chronocode-backend/internal/domain"
)

type ChronocodeServer struct {
	httpServer *http.Server
	engine     *gin.Engine
	logger     domain.Logger
}

func NewChronocodeServer(port string, logger domain.Logger, repoAnalyzer *application.RepositoryAnalyzer) *ChronocodeServer {
	engine := gin.Default()

	repoAnalyzerHandler := NewAnalysisHandler(repoAnalyzer, logger)

	engine.GET("/analyze-repository", repoAnalyzerHandler.AnalyzeRepository)

	httpServer := &http.Server{
		Addr:    port,
		Handler: engine, // Use gin engine as the HTTP handler
	}

	return &ChronocodeServer{
		httpServer: httpServer,
		engine:     engine,
		logger:     logger,
	}
}

func (s *ChronocodeServer) Run() error {
	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *ChronocodeServer) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
