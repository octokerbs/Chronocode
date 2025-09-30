package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/octokerbs/chronocode-go/internal/service"
)

type AnalysisHandler struct {
	repoAnalyzer *service.RepositoryAnalyzer
}

func NewAnalysisHandler(repoAnalyzer *service.RepositoryAnalyzer) *AnalysisHandler {
	return &AnalysisHandler{repoAnalyzer}
}

func (h *AnalysisHandler) AnalyzeRepository(c *gin.Context) {
	repoURL := c.Query("repo_url")
	accessToken := c.Query("access_token")

	if repoURL == "" || accessToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "repo_url and access_token are required",
		})
		return
	}

	commits, errors, err := h.repoAnalyzer.AnalyzeRepository(c.Request.Context(), repoURL, accessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	advisory := ""
	if len(commits) > 20 {
		advisory = "Not all commits were analyzed due to repository analysis limit reached"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":         "success",
		"analyses_count": len(commits),
		"commits":        commits,
		"advisory":       advisory,
		"errors":         errors,
	})
}
