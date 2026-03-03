package model

type SubcommitJSON struct {
	ID          int64    `json:"id"`
	CreatedAt   string   `json:"createdAt"`
	Title       string   `json:"title"`
	Idea        string   `json:"idea"`
	Description string   `json:"description"`
	CommitSHA   string   `json:"commitSha"`
	Type        string   `json:"type"`
	Epic        string   `json:"epic"`
	Files       []string `json:"files"`
}
