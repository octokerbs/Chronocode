package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/octokerbs/chronocode-backend/internal/application"
)

type QuerierHandler struct {
	Querier *application.Querier
}

func NewQuerierHandler(querier *application.Querier) *QuerierHandler {
	return &QuerierHandler{
		Querier: querier,
	}
}

func (q *QuerierHandler) GetSubcommits(c *gin.Context) {
	repoID := c.Query("repo_id")
	if repoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Falta el parámetro repoID."})
		return
	}

	subcommits, err := q.Querier.GetSubcommitsFromRepo(c.Request.Context(), repoID)
	if err != nil {
		httpErr := FromError(err)
		c.JSON(httpErr.Status, gin.H{"message": httpErr.Message})
		return
	}

	// Eliminar la alerta si se encuentran commits.
	if len(subcommits) > 0 {
		c.Header("HX-Trigger", "analysisComplete")
	}

	c.JSON(http.StatusOK, gin.H{
		"Subcommits": subcommits,
		"RepoID":     repoID,
	})
}
