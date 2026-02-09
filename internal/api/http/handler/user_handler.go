package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/octokerbs/chronocode-backend/internal/application"
)

type UserHandler struct {
	userProfile *application.UserProfile
}

func NewUserHandler(userProfile *application.UserProfile) *UserHandler {
	return &UserHandler{userProfile: userProfile}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	token := c.GetString("githubToken")

	profile, err := h.userProfile.GetProfile(c.Request.Context(), token)
	if err != nil {
		httpErr := FromError(err)
		c.JSON(httpErr.Status, gin.H{"error": httpErr.Message})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *UserHandler) GetRepositories(c *gin.Context) {
	token := c.GetString("githubToken")

	repos, err := h.userProfile.GetUserRepositories(c.Request.Context(), token)
	if err != nil {
		httpErr := FromError(err)
		c.JSON(httpErr.Status, gin.H{"error": httpErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"repositories": repos})
}

func (h *UserHandler) SearchRepositories(c *gin.Context) {
	token := c.GetString("githubToken")
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	repos, err := h.userProfile.SearchRepositories(c.Request.Context(), token, query)
	if err != nil {
		httpErr := FromError(err)
		c.JSON(httpErr.Status, gin.H{"error": httpErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"repositories": repos})
}
