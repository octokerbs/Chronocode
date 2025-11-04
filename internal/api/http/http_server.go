package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/octokerbs/chronocode-backend/internal/domain"
	"github.com/octokerbs/chronocode-backend/internal/service"
)

type HttpServer struct {
	httpServer *http.Server
	engine     *gin.Engine
	logger     domain.Logger
}

func NewServer(port string, logger domain.Logger, repoAnalyzer *service.RepositoryAnalyzerService, timeline *service.TimelineService) *HttpServer {
	engine := gin.Default()

	repoAnalyzerHandler := NewAnalysisHandler(repoAnalyzer, logger)
	engine.POST("/analyze-repository", repoAnalyzerHandler.AnalyzeRepository)

	timelineHandler := NewTimelineHandler(timeline, logger)
	engine.GET("/repository-timeline", timelineHandler.GetRepositoryTimeline)

	httpServer := &http.Server{
		Addr:    port,
		Handler: engine,
	}

	return &HttpServer{
		httpServer: httpServer,
		engine:     engine,
		logger:     logger,
	}
}

func (s *HttpServer) Run() error {
	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *HttpServer) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
