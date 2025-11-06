package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	http_analysis "github.com/octokerbs/chronocode-backend/internal/api/http/handler/analysis"
	"github.com/octokerbs/chronocode-backend/internal/domain"
	"github.com/octokerbs/chronocode-backend/internal/service/analysis"
	"github.com/octokerbs/chronocode-backend/internal/service/query"
)

type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer(analyzer *analysis.Analyzer, querier *query.Querier, port string, logger domain.Logger) *HTTPServer {
	engine := gin.Default()

	analysisHandler := http_analysis.NewAnalysisHandler(analyzer, querier, logger)
	engine.POST("/analyze", analysisHandler.AnalyzeRepository)
	engine.GET("/subcommits", analysisHandler.GetSubcommits)

	server := &http.Server{
		Addr:    port,
		Handler: engine,
	}

	return &HTTPServer{
		server: server,
	}
}

func (s *HTTPServer) Run() error {
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
