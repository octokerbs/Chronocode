package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/octokerbs/chronocode-backend/internal/application"
)

type AuthHandler struct {
	auth        *application.Auth
	frontendURL string
}

func NewAuthHandler(authService *application.Auth, frontendURL string) *AuthHandler {
	return &AuthHandler{
		auth:        authService,
		frontendURL: frontendURL,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	url := h.auth.GetLoginURL()
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) LoginCallback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	accessToken, err := h.auth.HandleCallback(c.Request.Context(), state, code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication failed"})
		return
	}

	c.SetCookie("access_token", accessToken, int(time.Hour*24*7/time.Second), "/", "", false, true)

	c.Redirect(http.StatusTemporaryRedirect, h.frontendURL)
}
