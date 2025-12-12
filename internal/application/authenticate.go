package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/octokerbs/chronocode-backend/internal/domain/auth"
)

const OAuthStateString = "secure-random-string-for-security"

type Auth struct {
	auth auth.Auth
}

func NewAuth(auth auth.Auth) *Auth {
	return &Auth{
		auth,
	}
}

func (a *Auth) GetLoginURL() string {
	return a.auth.GetAuthURL(OAuthStateString)
}

func (a *Auth) HandleCallback(ctx context.Context, state string, code string) (string, error) {
	if state != OAuthStateString {
		return "", errors.New("oauth state mismatch")
	}

	if code == "" {
		return "", errors.New("authorization code not found")
	}

	token, err := a.auth.ExchangeCode(ctx, code)
	if err != nil {
		return "", fmt.Errorf("failed to exchange code for token: %w", err)
	}

	return token.AccessToken, nil
}
