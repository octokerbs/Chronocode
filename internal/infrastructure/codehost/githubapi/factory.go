package githubapi

import (
	"context"

	"github.com/octokerbs/chronocode-backend/internal/domain"
)

// Creado para que la aplicacion orquestee el dominio sin saber que Code Host debe crear.

type GithubFactory struct{}

func NewGithubFactory() *GithubFactory {
	return &GithubFactory{}
}

func (f *GithubFactory) Create(ctx context.Context, accessToken string) (domain.CodeHost, error) {
	return NewGithubClient(ctx, accessToken)
}
