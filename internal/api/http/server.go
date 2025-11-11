package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/octokerbs/chronocode-backend/internal/api/http/handler"
	"github.com/octokerbs/chronocode-backend/internal/application"
	"github.com/octokerbs/chronocode-backend/internal/log"
)

type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer(analyzer *application.Analyzer, querier *application.Querier, auth *application.Auth, port string, logger log.Logger) *HTTPServer {
	engine := gin.Default()
	engine.LoadHTMLGlob("web/templates/*")

	server := &http.Server{
		Addr:    port,
		Handler: engine,
	}

	s := &HTTPServer{
		server: server,
	}

	s.registerPublicRoutes(engine, auth, logger)
	s.registerAuthenticatedRoutes(engine, analyzer, querier, logger)

	return s
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

func (s *HTTPServer) registerPublicRoutes(engine *gin.Engine, auth *application.Auth, logger log.Logger) {
	engine.GET("/", s.renderHomePage)

	authHandler := handler.NewAuthHandler(auth, logger)
	engine.GET("/auth/github/login", authHandler.Login)
	engine.GET("/auth/github/callback", authHandler.LoginCallback)
}

func (s *HTTPServer) registerAuthenticatedRoutes(engine *gin.Engine, analyzer *application.Analyzer, querier *application.Querier, logger log.Logger) {
	authenticated := engine.Group("/")
	authenticated.Use(s.authMiddleware())
	{
		analysisHandler := handler.NewAnalyzerHandler(analyzer, logger)
		authenticated.POST("/analyze", analysisHandler.AnalyzeRepository)

		querierHandler := handler.NewQuerierHandler(querier, logger)
		authenticated.GET("/subcommits-timeline", querierHandler.GetSubcommits)
	}
}

// Middleware simple para verificar si el usuario tiene un token
func (s *HTTPServer) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("access_token")
		if err != nil || token == "" {
			c.Redirect(http.StatusTemporaryRedirect, "/auth/github/login")
			c.Abort()
			return
		}
		c.Set("githubToken", token) // Pasa el token al contexto de Gin
		c.Next()
	}
}

func (s *HTTPServer) renderHomePage(c *gin.Context) {
	_, err := c.Cookie("access_token")
	loggedIn := err == nil
	c.HTML(http.StatusOK, "index.html", gin.H{
		"IsLoggedIn": loggedIn,
	})
}
