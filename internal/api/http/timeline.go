package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/octokerbs/chronocode-backend/internal/application"
	"github.com/octokerbs/chronocode-backend/internal/domain"
)

type TimelineHandler struct {
	timelineService *application.TimelineService
	logger          domain.Logger
}

func NewTimelineHandler(service *application.TimelineService, logger domain.Logger) *TimelineHandler {
	return &TimelineHandler{
		timelineService: service,
		logger:          logger,
	}
}

func (h *TimelineHandler) GetRepositoryTimeline(c *gin.Context) {
	repoID := c.Query("repo_id")

	subcommits, err := h.timelineService.GetSubcommitsFromRepo(c.Request.Context(), repoID)
	if err != nil {
		httpErr := FromError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": httpErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Subcommits": subcommits})
}
