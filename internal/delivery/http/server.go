package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	application_analysis "github.com/octokerbs/chronocode-backend/internal/application/analysis"
	application_identity "github.com/octokerbs/chronocode-backend/internal/application/identity"
	"github.com/octokerbs/chronocode-backend/internal/application/query"
	http_analysis "github.com/octokerbs/chronocode-backend/internal/delivery/http/handler/analysis"
	http_identity "github.com/octokerbs/chronocode-backend/internal/delivery/http/handler/identity"
	http_querier "github.com/octokerbs/chronocode-backend/internal/delivery/http/handler/query"
	"github.com/octokerbs/chronocode-backend/pkg/log"
)

type HTTPServer struct {
	server      *http.Server
	authService *application_identity.AuthService
}

func NewHTTPServer(analyzer *application_analysis.RepositoryAnalyzerService, querier *query.QuerierService, authService *application_identity.AuthService, port string, logger log.Logger) *HTTPServer {
	engine := gin.Default()
	engine.LoadHTMLGlob("web/templates/*")

	server := &http.Server{
		Addr:    port,
		Handler: engine,
	}

	s := &HTTPServer{
		server:      server,
		authService: authService,
	}

	s.registerPublicRoutes(engine, logger)
	s.registerAuthenticatedRoutes(engine, analyzer, querier, logger)

	return s
}

func (s *HTTPServer) registerPublicRoutes(engine *gin.Engine, logger log.Logger) {
	authHandler := http_identity.NewAuthHandler(s.authService, logger)

	engine.GET("/", authHandler.RenderHomePage)
	engine.GET("/auth/github/login", authHandler.GithubLogin)
	engine.GET("/auth/github/callback", authHandler.GithubCallback)
}

func (s *HTTPServer) registerAuthenticatedRoutes(engine *gin.Engine, analyzer *application_analysis.RepositoryAnalyzerService, querier *query.QuerierService, logger log.Logger) {
	authenticated := engine.Group("/")
	authenticated.Use(s.authMiddleware())
	{
		analysisHandler := http_analysis.NewAnalysisHandler(analyzer, logger)
		authenticated.POST("/analyze", analysisHandler.AnalyzeRepository)

		querierHandler := http_querier.NewQuerierHandler(querier, logger)
		authenticated.GET("/subcommits-timeline", querierHandler.GetSubcommits())
	}
}

// Middleware simple para verificar si el usuario tiene un token
func (s *HTTPServer) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("github_access_token")
		if err != nil || token == "" {
			c.Redirect(http.StatusTemporaryRedirect, "/auth/github/login")
			c.Abort()
			return
		}
		c.Set("githubToken", token) // Pasa el token al contexto de Gin
		c.Next()
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
