package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/octokerbs/chronocode-backend/internal/application"
	"github.com/octokerbs/chronocode-backend/internal/domain"
)

type Server struct {
	httpServer *http.Server
	engine     *gin.Engine
	logger     domain.Logger
}

func NewServer(port string, logger domain.Logger, repoAnalyzer *application.RepositoryAnalyzer, timeline *application.TimelineService) *Server {
	engine := gin.Default()

	repoAnalyzerHandler := NewAnalysisHandler(repoAnalyzer, logger)
	engine.POST("/analyze-repository", repoAnalyzerHandler.AnalyzeRepository)

	timelineHandler := NewTimelineHandler(timeline, logger)
	engine.GET("/repository-timeline", timelineHandler.GetRepositoryTimeline)

	httpServer := &http.Server{
		Addr:    port,
		Handler: engine,
	}

	return &Server{
		httpServer: httpServer,
		engine:     engine,
		logger:     logger,
	}
}

func (s *Server) Run() error {
	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
