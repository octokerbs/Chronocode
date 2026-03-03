package utils

import (
	"time"

	"github.com/octokerbs/chronocode/internal/domain/subcommit"
	"github.com/octokerbs/chronocode/internal/ports/http/model"
)

func MapSubcommits(scs []subcommit.Subcommit) []model.SubcommitJSON {
	result := make([]model.SubcommitJSON, len(scs))
	for i, sc := range scs {
		result[i] = model.SubcommitJSON{
			ID:          sc.ID(),
			CreatedAt:   sc.CommittedAt().Format(time.RFC3339),
			Title:       sc.Title(),
			Idea:        sc.Idea(),
			Description: sc.Description(),
			CommitSHA:   sc.CommitSHA(),
			Type:        sc.ModificationType(),
			Epic:        sc.Epic(),
			Files:       sc.Files(),
		}
	}
	return result
}
