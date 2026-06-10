package redis

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/CommitUpp/backend/auth/domain/repository"
	goredis "github.com/redis/go-redis/v9"
)

type tokenCacheRepositoryImpl struct {
	client *goredis.Client
}

func NewTokenCacheRepository(client *goredis.Client) repository.TokenCacheRepository {
	return &tokenCacheRepositoryImpl{
		client: client,
	}
}

func (r *tokenCacheRepositoryImpl) GetUserID(ctx context.Context, accessToken string) (string, error) {
	startedAt := time.Now()

	userID, err := r.client.Get(ctx, cacheKey(accessToken)).Result()
	switch {
	case err == nil && userID != "":
		observeTokenCacheOperation("get", "hit", startedAt)
	case err == goredis.Nil || (err == nil && userID == ""):
		observeTokenCacheOperation("get", "miss", startedAt)
	default:
		observeTokenCacheOperation("get", "error", startedAt)
	}

	return userID, err
}

func (r *tokenCacheRepositoryImpl) SetUserID(ctx context.Context, accessToken string, userID string, ttl time.Duration) error {
	startedAt := time.Now()

	err := r.client.Set(ctx, cacheKey(accessToken), userID, ttl).Err()
	observeTokenCacheOperation("set", redisSetResult(err), startedAt)

	return err
}

func (r *tokenCacheRepositoryImpl) DeleteUserID(ctx context.Context, accessToken string) error {
	startedAt := time.Now()

	err := r.client.Del(ctx, cacheKey(accessToken)).Err()
	observeTokenCacheOperation("delete", redisSetResult(err), startedAt)

	return err
}

func cacheKey(accessToken string) string {
	hash := sha256.Sum256([]byte(accessToken))
	return "auth:access_token:" + hex.EncodeToString(hash[:])
}
