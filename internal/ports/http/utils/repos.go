package utils

import (
	"strconv"
	"time"

	"github.com/octokerbs/chronocode/internal/domain/repo"
	"github.com/octokerbs/chronocode/internal/ports/http/model"
)

func FormatInt64(n int64) string {
	return strconv.FormatInt(n, 10)
}

func MapRepos(repos []*repo.Repo) []model.RepoJSON {
	result := make([]model.RepoJSON, len(repos))
	for i, r := range repos {
		result[i] = model.RepoJSON{
			ID:      FormatInt64(r.ID()),
			Name:    r.Name(),
			URL:     r.URL(),
			AddedAt: r.CreatedAt().Format(time.RFC3339),
		}
	}
	return result
}
