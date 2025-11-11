package auth

import (
	"context"

	"golang.org/x/oauth2"
)

type Auth interface {
	GetAuthURL(state string) string
	ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error)
}
