package transport

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/octokerbs/chronocode-backend/config"
	"github.com/octokerbs/chronocode-backend/pkg/logger"
	"go.uber.org/zap"
)

type HTTPServer struct {
	r    *gin.Engine
	port string
}

func NewHTTPServer(port string) HTTPServer {
	r := gin.New()
	r.Use(ginzap.Ginzap(logger.Log, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger.Log, true))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://frontend:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	userHandler, noteHandler := setupHandlers(db)
	defineRoutes(r, userHandler, noteHandler)

	return HTTPServer{r, port}
}
