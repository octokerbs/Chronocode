package identity

import (
	"context"
	"errors"
	"fmt"

	"github.com/octokerbs/chronocode-backend/internal/domain/identity"
)

const OAuthStateString = "secure-random-string-for-security"

type AuthService struct {
	auth identity.Auth
}

func NewAuthService(auth identity.Auth) *AuthService {
	return &AuthService{
		auth,
	}
}

func (s *AuthService) GetLoginURL() string {
	return s.auth.GetAuthURL(OAuthStateString)
}

func (s *AuthService) HandleCallback(ctx context.Context, state string, code string) (string, error) {
	if state != OAuthStateString {
		return "", errors.New("oauth state mismatch")
	}

	if code == "" {
		return "", errors.New("authorization code not found")
	}

	token, err := s.auth.ExchangeCode(ctx, code)
	if err != nil {
		return "", fmt.Errorf("failed to exchange code for token: %w", err)
	}

	return token.AccessToken, nil
}
