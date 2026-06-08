package repository

import "context"

type VerifiedToken struct {
	UserID string
	TTL    int64
}

type TokenVerifierRepository interface {
	Verify(ctx context.Context, accessToken string) (*VerifiedToken, error)
}
