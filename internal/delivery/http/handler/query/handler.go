package query

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/octokerbs/chronocode-backend/internal/application/query"
	"github.com/octokerbs/chronocode-backend/internal/delivery/http/handler"
	"github.com/octokerbs/chronocode-backend/pkg/log"
)

type QuerierHandler struct {
	Querier *query.QuerierService
	logger  log.Logger
}

func NewQuerierHandler(querier *query.QuerierService, logger log.Logger) *QuerierHandler {
	return &QuerierHandler{
		Querier: querier,
		logger:  logger,
	}
}

func (q *QuerierHandler) GetSubcommits(c *gin.Context) {
	repoID := c.Query("repo_id")
	if repoID == "" {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"Message": "Falta el parámetro repoID."})
		return
	}

	subcommits, err := q.Querier.GetSubcommitsFromRepo(c.Request.Context(), repoID)
	if err != nil {
		httpErr := handler.FromError(err)
		c.HTML(httpErr.Status, "error.html", gin.H{"message": httpErr.Message})
		return
	}

	// Eliminar la alerta si se encuentran commits.
	if len(subcommits) > 0 {
		c.Header("HX-Trigger", "analysisComplete")
	}

	c.HTML(http.StatusOK, "subcommits_timeline.html", gin.H{
		"Subcommits": subcommits,
		"RepoID":     repoID,
	})
}
