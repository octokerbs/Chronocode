package cache

import "context"

type Cache interface {
	AddUserRepository(ctx context.Context, userID string, repo UserRepository) error
	GetUserRepositories(ctx context.Context, userID string) ([]UserRepository, error)
	RemoveUserRepository(ctx context.Context, userID string, repoID string) error
}
