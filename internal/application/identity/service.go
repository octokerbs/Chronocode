package identity

import (
	"context"
	"errors"
	"fmt"

	"github.com/octokerbs/chronocode-backend/internal/infrastructure/identity/githubauth"
)

const OAuthStateString = "secure-random-string-for-security"

type AuthService struct {
	githubProvider *githubauth.GitHubAuthenticationProvider
}

func NewAuthService(provider *githubauth.GitHubAuthenticationProvider) *AuthService {
	return &AuthService{
		githubProvider: provider,
	}
}

func (s *AuthService) GetLoginURL() string {
	return s.githubProvider.GetAuthURL(OAuthStateString)
}

func (s *AuthService) HandleCallback(ctx context.Context, state string, code string) (accessToken string, err error) {
	if state != OAuthStateString {
		return "", errors.New("oauth state mismatch")
	}

	if code == "" {
		return "", errors.New("authorization code not found")
	}

	token, err := s.githubProvider.ExchangeCode(ctx, code)
	if err != nil {
		return "", fmt.Errorf("failed to exchange code for token: %w", err)
	}

	return token.AccessToken, nil
}
