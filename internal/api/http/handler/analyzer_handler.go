package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/octokerbs/chronocode-backend/internal/application"
	"github.com/octokerbs/chronocode-backend/internal/log"
)

type AnalyzerHandler struct {
	analyzer *application.Analyzer
	logger   log.Logger
}

func NewAnalyzerHandler(analyzer *application.Analyzer, logger log.Logger) *AnalyzerHandler {
	return &AnalyzerHandler{
		analyzer: analyzer,
		logger:   logger,
	}
}

func (h *AnalyzerHandler) AnalyzeRepository(c *gin.Context) {
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

	repo, codeHost, err := h.analyzer.PrepareAnalysis(c.Request.Context(), repoURL, githubToken)
	if err != nil {
		httpErr := FromError(err)

		if httpErr.Status == 0 { // Empty status indicates internal server error
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{"message": httpErr.Message})
			return
		}

		c.HTML(httpErr.Status, "error.html", gin.H{"Message": httpErr.Message})
		return
	}

	repoID := repo.ID
	go func() {
		analysisCtx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		if err := h.analyzer.RunAnalysis(analysisCtx, repo, codeHost); err != nil {
			h.logger.Error("Background analysis failed", err, "repoURL", repoURL)
		} else {
			h.logger.Info("Background analysis complete", "repoURL", repoURL)
		}
	}()

	c.HTML(http.StatusAccepted, "analysis_status.html", gin.H{
		"Message": "Análisis del repositorio iniciado. Cargando línea de tiempo...",
		"RepoID":  repoID,
	})
}
