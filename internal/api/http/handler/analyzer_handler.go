package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/octokerbs/chronocode-backend/internal/application"
	"github.com/octokerbs/chronocode-backend/internal/domain/cache"
	"github.com/octokerbs/chronocode-backend/internal/domain/codehost"
)

type AnalyzerHandler struct {
	prepareRepo     *application.PrepareRepository
	analyzer        *application.Analyzer
	persistCommits  *application.PersistCommits
	cache           cache.Cache
	codeHostFactory codehost.CodeHostFactory
}

func NewAnalyzerHandler(prepareRepo *application.PrepareRepository, analyzer *application.Analyzer, persistCommits *application.PersistCommits, cachePort cache.Cache, codeHostFactory codehost.CodeHostFactory) *AnalyzerHandler {
	return &AnalyzerHandler{
		prepareRepo:     prepareRepo,
		analyzer:        analyzer,
		persistCommits:  persistCommits,
		cache:           cachePort,
		codeHostFactory: codeHostFactory,
	}
}

type analyzeRequest struct {
	RepoURL string `json:"repoUrl"`
}

func (h *AnalyzerHandler) AnalyzeRepository(c *gin.Context) {
	token, exists := c.Get("githubToken")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	githubToken := token.(string)

	var req analyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.RepoURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "repository URL is required"})
		return
	}

	repo, err := h.prepareRepo.Execute(c.Request.Context(), req.RepoURL, githubToken)
	if err != nil {
		httpErr := FromError(err)
		c.JSON(httpErr.Status, gin.H{"error": httpErr.Message})
		return
	}

	events := make(chan application.CommitAnalyzed, 100)
	go h.persistCommits.HandleCommitAnalyzed(context.Background(), events)

	go func() {
		defer close(events)
		h.analyzer.AnalyzeCommits(context.Background(), repo, events, githubToken)
	}()

	go func() {
		ch := h.codeHostFactory.Create(context.Background(), githubToken)
		profile, err := ch.FetchAuthenticatedUser(context.Background())
		if err != nil {
			return
		}
		userID := fmt.Sprintf("%d", profile.ID)
		h.cache.AddUserRepository(context.Background(), userID, cache.UserRepository{
			ID:      fmt.Sprintf("%d", repo.ID),
			Name:    repo.Name,
			URL:     repo.URL,
			AddedAt: time.Now().Format(time.RFC3339),
		})
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Repository analysis started. Loading timeline...",
		"repoId":  repo.ID,
	})
}
