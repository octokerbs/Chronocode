package analysis

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/octokerbs/chronocode-backend/internal/api/http/handler"
	"github.com/octokerbs/chronocode-backend/internal/domain"
	"github.com/octokerbs/chronocode-backend/internal/service/analysis"
	"github.com/octokerbs/chronocode-backend/internal/service/query"
)

type AnalysisHandler struct {
	Analyzer *analysis.Analyzer
	Querier  *query.Querier
	logger   domain.Logger
}

func NewAnalysisHandler(analyzer *analysis.Analyzer, querier *query.Querier, logger domain.Logger) *AnalysisHandler {
	return &AnalysisHandler{
		Analyzer: analyzer,
		Querier:  querier,
		logger:   logger,
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

	repo, codeHost, err := h.Analyzer.PrepareAnalysis(c.Request.Context(), repoURL, accessToken)
	if err != nil {
		httpErr := handler.FromError(err)

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

		if err := h.Analyzer.RunAnalysis(analysisCtx, repo, codeHost); err != nil {
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

func (h *AnalysisHandler) GetSubcommits(c *gin.Context) {
	repoID := c.Query("repo_id")

	subcommits, err := h.Querier.GetSubcommitsFromRepo(c.Request.Context(), repoID)
	if err != nil {
		httpErr := handler.FromError(err)
		c.JSON(httpErr.Status, gin.H{"message": httpErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"subcommits": subcommits})
}
