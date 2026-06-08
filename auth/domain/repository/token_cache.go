package repository

import (
	"context"
	"time"
)

type TokenCacheRepository interface {
	GetUserID(ctx context.Context, accessToken string) (string, error)
	SetUserID(ctx context.Context, accessToken string, userID string, ttl time.Duration) error
}
