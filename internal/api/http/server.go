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
	"github.com/octokerbs/chronocode-backend/internal/domain/cache"
	"github.com/octokerbs/chronocode-backend/internal/domain/codehost"
)

type HTTPServer struct {
	server      *http.Server
	frontendURL string
}

func NewHTTPServer(analyzer *application.Analyzer, persistCommits *application.PersistCommits, prepareRepo *application.PrepareRepository, querier *application.Querier, auth *application.Auth, userProfile *application.UserProfile, cachePort cache.Cache, codeHostFactory codehost.CodeHostFactory, port string, frontendURL string) *HTTPServer {
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

	analysisHandler := handler.NewAnalyzerHandler(prepareRepo, analyzer, persistCommits, cachePort, codeHostFactory)
	querierHandler := handler.NewQuerierHandler(querier)
	userHandler := handler.NewUserHandler(userProfile)

	s.registerPublicRoutes(engine, auth)
	s.registerAuthenticatedRoutes(engine, analysisHandler, querierHandler, userHandler)

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

	engine.POST("/auth/logout", func(c *gin.Context) {
		c.SetCookie("access_token", "", -1, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{"message": "logged out"})
	})
}

func (s *HTTPServer) registerAuthenticatedRoutes(engine *gin.Engine, analysisHandler *handler.AnalyzerHandler, querierHandler *handler.QuerierHandler, userHandler *handler.UserHandler) {
	authenticated := engine.Group("/")
	authenticated.Use(s.authMiddleware())
	{
		authenticated.POST("/analyze", analysisHandler.AnalyzeRepository)
		authenticated.GET("/subcommits-timeline", querierHandler.GetSubcommits)
		authenticated.GET("/user/profile", userHandler.GetProfile)
		authenticated.GET("/repositories", userHandler.GetRepositories)
		authenticated.GET("/user/repos/search", userHandler.SearchRepositories)
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
