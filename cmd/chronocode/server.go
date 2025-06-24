package main

import (
	"log"

	"github.com/chrono-code-hackathon/chronocode-go/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	log.SetFlags(log.LstdFlags)

	// Production
	r.GET("/analyze-repository", handlers.AnalyzeRepositoryHandler)

	r.Run("localhost:8080")
}
