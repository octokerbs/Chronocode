package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/octokerbs/chronocode-backend/internal/domain/cache"
	pkg_errors "github.com/octokerbs/chronocode-backend/internal/errors"
)

const userReposTTL = 30 * 24 * time.Hour

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(redisURL string) (*RedisCache, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, pkg_errors.NewError(pkg_errors.ErrInternalFailure, fmt.Errorf("invalid redis URL: %w", err))
	}

	client := redis.NewClient(opts)

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, pkg_errors.NewError(pkg_errors.ErrInternalFailure, fmt.Errorf("redis ping failed: %w", err))
	}

	return &RedisCache{client: client}, nil
}

func (rc *RedisCache) AddUserRepository(ctx context.Context, userID string, repo cache.UserRepository) error {
	key := fmt.Sprintf("user:%s:repos", userID)

	data, err := json.Marshal(repo)
	if err != nil {
		return pkg_errors.NewError(pkg_errors.ErrInternalFailure, err)
	}

	if err := rc.client.HSet(ctx, key, repo.ID, string(data)).Err(); err != nil {
		return pkg_errors.NewError(pkg_errors.ErrInternalFailure, err)
	}

	rc.client.Expire(ctx, key, userReposTTL)
	return nil
}

func (rc *RedisCache) GetUserRepositories(ctx context.Context, userID string) ([]cache.UserRepository, error) {
	key := fmt.Sprintf("user:%s:repos", userID)

	result, err := rc.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, pkg_errors.NewError(pkg_errors.ErrInternalFailure, err)
	}

	repos := make([]cache.UserRepository, 0, len(result))
	for _, val := range result {
		var repo cache.UserRepository
		if err := json.Unmarshal([]byte(val), &repo); err != nil {
			continue
		}
		repos = append(repos, repo)
	}

	return repos, nil
}

func (rc *RedisCache) RemoveUserRepository(ctx context.Context, userID string, repoID string) error {
	key := fmt.Sprintf("user:%s:repos", userID)

	if err := rc.client.HDel(ctx, key, repoID).Err(); err != nil {
		return pkg_errors.NewError(pkg_errors.ErrInternalFailure, err)
	}

	return nil
}
