package identity

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/octokerbs/chronocode-backend/internal/application/identity"
	"github.com/octokerbs/chronocode-backend/pkg/log"
)

type AuthHandler struct {
	authService *identity.AuthService
	logger      log.Logger
}

func NewAuthHandler(authService *identity.AuthService, logger log.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

func (h *AuthHandler) GithubLogin(c *gin.Context) {
	url := h.authService.GetLoginURL()
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) GithubCallback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	accessToken, err := h.authService.HandleCallback(c.Request.Context(), state, code)
	if err != nil {
		h.logger.Error("GitHub authentication failed: %v", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Lógica de HTTP: Establecer la cookie
	c.SetCookie("github_access_token", accessToken, int(time.Hour*24*7/time.Second), "/", "localhost", false, true)
	h.logger.Info("Successfully logged in with GitHub.")

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (h *AuthHandler) RenderHomePage(c *gin.Context) {
	_, err := c.Cookie("github_access_token")
	loggedIn := err == nil
	c.HTML(http.StatusOK, "index.html", gin.H{
		"IsLoggedIn": loggedIn,
	})
}
