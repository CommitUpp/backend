package redis

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/CommitUpp/backend/auth/domain/repository"
	"github.com/redis/go-redis/v9"
)

type tokenCacheRepositoryImpl struct {
	client *redis.Client
}

func NewTokenCacheRepository(client *redis.Client) repository.TokenCacheRepository {
	return &tokenCacheRepositoryImpl{
		client: client,
	}
}

func (r *tokenCacheRepositoryImpl) GetUserID(ctx context.Context, accessToken string) (string, error) {
	return r.client.Get(ctx, cacheKey(accessToken)).Result()
}

func (r *tokenCacheRepositoryImpl) SetUserID(ctx context.Context, accessToken string, userID string, ttl time.Duration) error {
	return r.client.Set(ctx, cacheKey(accessToken), userID, ttl).Err()
}

func cacheKey(accessToken string) string {
	hash := sha256.Sum256([]byte(accessToken))
	return "auth:access_token:" + hex.EncodeToString(hash[:])
}
