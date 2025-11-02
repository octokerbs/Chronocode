package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/octokerbs/chronocode-backend/internal/api/httperror"
	"github.com/octokerbs/chronocode-backend/internal/application"
)

type AnalysisHandler struct {
	repoAnalyzer *application.RepositoryAnalyzer
}

func NewAnalysisHandler(repoAnalyzer *application.RepositoryAnalyzer) *AnalysisHandler {
	return &AnalysisHandler{repoAnalyzer}
}

func (h *AnalysisHandler) AnalyzeRepository(c *gin.Context) {
	repoURL := c.Query("repo_url")
	authHeader := c.GetHeader("Authorization")

	if repoURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "repo_url query parameter is required"})
		return
	}

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is required"})
		return
	}

	var accessToken string
	if _, err := fmt.Sscanf(authHeader, "Bearer %s", &accessToken); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Authorization header format"})
		return
	}

	repo, codeHost, err := h.repoAnalyzer.PrepareAnalysis(c.Request.Context(), repoURL, accessToken)
	if err != nil {
		httpErr := httperror.FromError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": httpErr.Status, "message": httpErr.Message})
		return
	}

	go func() {
		analysisCtx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		if err := h.repoAnalyzer.RunAnalysis(analysisCtx, repo, codeHost); err != nil {
			log.Printf("Background analysis failed for %s: %v", repoURL, err)
		} else {
			log.Printf("Background analysis complete for %s", repoURL)
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"status":  "pending",
		"message": "Repository analysis has been queued.",
	})
}
