package http

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/octokerbs/chronocode-backend/internal/api/http/handler"
	"github.com/octokerbs/chronocode-backend/internal/application"
)

type HTTPServer struct {
	server      *http.Server
	frontendURL string
}

func NewHTTPServer(analyzer *application.Analyzer, persistCommits *application.PersistCommits, prepareRepo *application.PrepareRepository, querier *application.Querier, auth *application.Auth, port string, frontendURL string) *HTTPServer {
	engine := gin.Default()

	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{frontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	server := &http.Server{
		Addr:    port,
		Handler: engine,
	}

	s := &HTTPServer{
		server:      server,
		frontendURL: frontendURL,
	}

	analysisHandler := handler.NewAnalyzerHandler(prepareRepo, analyzer, persistCommits)
	querierHandler := handler.NewQuerierHandler(querier)

	s.registerPublicRoutes(engine, auth)
	s.registerAuthenticatedRoutes(engine, analysisHandler, querierHandler)

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

func (s *HTTPServer) registerPublicRoutes(engine *gin.Engine, auth *application.Auth) {
	engine.GET("/auth/status", s.authStatus)

	authHandler := handler.NewAuthHandler(auth, s.frontendURL)
	engine.GET("/auth/github/login", authHandler.Login)
	engine.GET("/auth/github/callback", authHandler.LoginCallback)
}

func (s *HTTPServer) registerAuthenticatedRoutes(engine *gin.Engine, analysisHandler *handler.AnalyzerHandler, querierHandler *handler.QuerierHandler) {
	authenticated := engine.Group("/")
	authenticated.Use(s.authMiddleware())
	{
		authenticated.POST("/analyze", analysisHandler.AnalyzeRepository)
		authenticated.GET("/subcommits-timeline", querierHandler.GetSubcommits)
	}
}

func (s *HTTPServer) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("access_token")
		if err != nil || token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		c.Set("githubToken", token)
		c.Next()
	}
}

func (s *HTTPServer) authStatus(c *gin.Context) {
	_, err := c.Cookie("access_token")
	loggedIn := err == nil
	c.JSON(http.StatusOK, gin.H{
		"isLoggedIn": loggedIn,
	})
}
