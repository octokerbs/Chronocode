package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/octokerbs/chronocode-backend/internal/application"
	"github.com/octokerbs/chronocode-backend/internal/log"
)

type AuthHandler struct {
	auth   *application.Auth
	logger log.Logger
}

func NewAuthHandler(authService *application.Auth, logger log.Logger) *AuthHandler {
	return &AuthHandler{
		auth:   authService,
		logger: logger,
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
		h.logger.Error("authentication failed: %v", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Lógica de HTTP: Establecer la cookie
	c.SetCookie("access_token", accessToken, int(time.Hour*24*7/time.Second), "/", "localhost", false, true)
	h.logger.Info("Successfully logged in with Auth.")

	c.Redirect(http.StatusTemporaryRedirect, "/")
}
