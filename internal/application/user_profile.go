package application

import (
	"context"
	"fmt"

	"github.com/octokerbs/chronocode-backend/internal/domain/analysis"
	"github.com/octokerbs/chronocode-backend/internal/domain/cache"
	"github.com/octokerbs/chronocode-backend/internal/domain/codehost"
	"github.com/octokerbs/chronocode-backend/internal/domain/user"
)

type UserProfile struct {
	CodeHostFactory codehost.CodeHostFactory
	Cache           cache.Cache
}

func (up *UserProfile) GetProfile(ctx context.Context, accessToken string) (*user.GitHubProfile, error) {
	codeHost := up.CodeHostFactory.Create(ctx, accessToken)
	return codeHost.FetchAuthenticatedUser(ctx)
}

func (up *UserProfile) GetUserRepositories(ctx context.Context, accessToken string) ([]cache.UserRepository, error) {
	profile, err := up.GetProfile(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	userID := fmt.Sprintf("%d", profile.ID)
	return up.Cache.GetUserRepositories(ctx, userID)
}

func (up *UserProfile) SearchRepositories(ctx context.Context, accessToken string, query string) ([]*analysis.Repository, error) {
	codeHost := up.CodeHostFactory.Create(ctx, accessToken)
	return codeHost.SearchUserRepositories(ctx, query)
}
