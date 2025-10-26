package githubapi

import (
	"context"

	"github.com/octokerbs/chronocode-backend/internal/application"
)

// Creado para que la aplicacion orquestee el dominio sin saber que Code Host debe crear.

type GithubFactory struct{}

func NewGithubFactory() *GithubFactory {
	return &GithubFactory{}
}

func (f *GithubFactory) Create(ctx context.Context, accessToken string) (application.CodeHost, error) {
	return NewGithubClient(ctx, accessToken)
}
