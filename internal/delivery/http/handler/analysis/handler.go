package analysis

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/octokerbs/chronocode-backend/internal/application/analysis"
	"github.com/octokerbs/chronocode-backend/internal/delivery/http/handler"
	"github.com/octokerbs/chronocode-backend/pkg/log"
)

type AnalysisHandler struct {
	Analyzer *analysis.RepositoryAnalyzerService
	logger   log.Logger
}

func NewAnalysisHandler(analyzer *analysis.RepositoryAnalyzerService, logger log.Logger) *AnalysisHandler {
	return &AnalysisHandler{
		Analyzer: analyzer,
		logger:   logger,
	}
}

func (h *AnalysisHandler) AnalyzeRepository(c *gin.Context) {
	token, exists := c.Get("githubToken")
	if !exists {
		h.logger.Error("Bad request", fmt.Errorf("github token does not exist: %s", token))
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"Message": "Unauthorized. Please, login again."})
		return
	}
	githubToken := token.(string)

	repoURL := c.PostForm("repoUrl")
	if repoURL == "" {
		h.logger.Error("Bad request", fmt.Errorf("empty repo_url"))
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"Message": "URL del repositorio es requerida."})
		return
	}

	repo, codeHost, err := h.Analyzer.PrepareAnalysis(c.Request.Context(), repoURL, githubToken)
	if err != nil {
		httpErr := handler.FromError(err)

		if httpErr.Status == 0 { // Empty status indicates internal server error
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{"message": httpErr.Message})
			return
		}

		c.HTML(httpErr.Status, "error.html", gin.H{"Message": httpErr.Message})
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

	c.HTML(http.StatusAccepted, "analysis_status.html", gin.H{"message": "Repository analysis has been queued."})
}
