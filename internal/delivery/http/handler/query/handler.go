package query

import (
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/octokerbs/chronocode-backend/internal/application/query"
	"github.com/octokerbs/chronocode-backend/internal/domain/analysis"
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

func (q *QuerierHandler) GetSubcommits() gin.HandlerFunc {
	// repoID := c.Query("repo_id")

	// subcommits, err := q.Querier.GetSubcommitsFromRepo(c.Request.Context(), repoID)
	// if err != nil {
	// 	httpErr := handler.FromError(err)
	// 	c.JSON(httpErr.Status, gin.H{"message": httpErr.Message})
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{"subcommits": subcommits})
	return func(c *gin.Context) {

		// Datos simulados para demostración:
		subcommits := []analysis.Subcommit{
			{
				ID:        1,
				CreatedAt: ptrTime(time.Now().Add(-2 * 24 * time.Hour)),
				Title:     "Feature: Implementar OAuth2 con GitHub",
				Idea:      "Permitir a los usuarios autenticarse con sus cuentas de GitHub para acceder a funcionalidades avanzadas.",
				CommitSHA: "abcdef123",
				Type:      "feature",
				Epic:      "Autenticación",
				Files:     []string{"main.go", "handlers/auth.go"},
			},
			{
				ID:        2,
				CreatedAt: ptrTime(time.Now().Add(-1 * 24 * time.Hour)),
				Title:     "Fix: Errores de validación en el formulario de análisis",
				Idea:      "Corregir errores de validación y añadir mensajes de error más claros para el usuario.",
				CommitSHA: "fedcba456",
				Type:      "fix",
				Epic:      "Interfaz de Usuario",
				Files:     []string{"templates/index.html"},
			},
			{
				ID:        3,
				CreatedAt: ptrTime(time.Now()),
				Title:     "Refactor: Mejorar manejo de tokens en backend",
				Idea:      "Optimizar el almacenamiento y la recuperación del token de GitHub para mayor seguridad.",
				CommitSHA: "987654zyx",
				Type:      "refactor",
				Epic:      "Seguridad",
				Files:     []string{"http/server.go", "middlewares/auth.go"},
			},
		}

		sort.Slice(subcommits, func(i, j int) bool {
			if subcommits[i].CreatedAt == nil {
				return false
			}
			if subcommits[j].CreatedAt == nil {
				return true
			}
			return subcommits[i].CreatedAt.After(*subcommits[j].CreatedAt)
		})

		c.HTML(http.StatusOK, "subcommits_timeline.html", gin.H{
			"Subcommits": subcommits,
		})
	}
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
