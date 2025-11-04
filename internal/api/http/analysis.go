package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/octokerbs/chronocode-backend/internal/domain"
	"github.com/octokerbs/chronocode-backend/internal/service"
)

type AnalysisHandler struct {
	repoAnalyzer *service.RepositoryAnalyzerService
	logger       domain.Logger
}

func NewAnalysisHandler(repoAnalyzer *service.RepositoryAnalyzerService, logger domain.Logger) *AnalysisHandler {
	return &AnalysisHandler{
		repoAnalyzer: repoAnalyzer,
		logger:       logger,
	}
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
	if _, err := fmt.Sscanf(authHeader, "Bearer %s", &accessToken); err != nil { //...
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Authorization header format"})
		return
	}

	repo, codeHost, err := h.repoAnalyzer.PrepareAnalysis(c.Request.Context(), repoURL, accessToken)
	if err != nil {
		httpErr := FromError(err)

		if httpErr.Status == 0 { // Empty status indicates internal server error
			c.JSON(http.StatusInternalServerError, gin.H{"message": httpErr.Message})
			return
		}

		c.JSON(httpErr.Status, gin.H{"Message": httpErr.Message})
		return
	}

	go func() {
		analysisCtx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		if err := h.repoAnalyzer.RunAnalysis(analysisCtx, repo, codeHost); err != nil {
			h.logger.Error("Background analysis failed", err, "repoURL", repoURL)
		} else {
			h.logger.Info("Background analysis complete", "repoURL", repoURL)
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"status":  "pending",
		"message": "Repository analysis has been queued.",
	})
}
