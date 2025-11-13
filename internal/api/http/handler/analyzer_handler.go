package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/octokerbs/chronocode-backend/internal/application"
)

type AnalyzerHandler struct {
	prepareRepo    *application.PrepareRepository
	analyzer       *application.Analyzer
	persistCommits *application.PersistCommits
}

func NewAnalyzerHandler(prepareRepo *application.PrepareRepository, analyzer *application.Analyzer, persistCommits *application.PersistCommits) *AnalyzerHandler {
	return &AnalyzerHandler{
		prepareRepo:    prepareRepo,
		analyzer:       analyzer,
		persistCommits: persistCommits,
	}
}

func (h *AnalyzerHandler) AnalyzeRepository(c *gin.Context) {
	token := c.Query("github_token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized. Please, login again."})
		return
	}

	repoURL := c.Query("repo_url")
	if repoURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Empty repository name."})
		return
	}

	repo, err := h.prepareRepo.Execute(c.Request.Context(), repoURL, token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "URL del repositorio es requerida."})
		return
	}

	events := make(chan application.CommitAnalyzed, 100)
	go h.persistCommits.HandleCommitAnalyzed(context.Background(), events)

	go func() {
		defer close(events)
		h.analyzer.AnalyzeCommits(context.Background(), repo, events, token)
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"Message": "Análisis del repositorio iniciado. Cargando línea de tiempo...",
		"RepoID":  repo.ID,
	})
}
