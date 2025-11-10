package githubauth

import (
	"context"

	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

type GitHubAuth struct {
	config *oauth2.Config
}

func NewGitHubAuthenticationProvider(clientID, clientSecret, redirectURL string) *GitHubAuth {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"repo", "read:org", "user:email"},
		Endpoint:     githuboauth.Endpoint,
	}
	return &GitHubAuth{config: config}
}

func (a *GitHubAuth) GetAuthURL(state string) string {
	return a.config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

func (a *GitHubAuth) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return a.config.Exchange(ctx, code)
}
