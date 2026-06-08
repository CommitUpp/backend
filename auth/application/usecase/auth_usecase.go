package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/CommitUpp/backend/auth/domain/repository"
)

var ErrUnauthorized = errors.New("unauthorized")

type AuthUsecase interface {
	VerifyToken(ctx context.Context, accessToken string) (string, error)
}

type authUsecaseImpl struct {
	tokenCache    repository.TokenCacheRepository
	tokenVerifier repository.TokenVerifierRepository
}

func NewAuthUsecase(
	tc repository.TokenCacheRepository,
	tv repository.TokenVerifierRepository,
) AuthUsecase {
	return &authUsecaseImpl{
		tokenCache:    tc,
		tokenVerifier: tv,
	}
}

func (u *authUsecaseImpl) VerifyToken(ctx context.Context, accessToken string) (string, error) {
	if accessToken == "" {
		return "", ErrUnauthorized
	}

	userID, err := u.tokenCache.GetUserID(ctx, accessToken)
	if err == nil && userID != "" {
		return userID, nil
	}

	verifiedToken, err := u.tokenVerifier.Verify(ctx, accessToken)
	if err != nil || verifiedToken == nil || verifiedToken.UserID == "" {
		return "", ErrUnauthorized
	}

	ttl := time.Duration(verifiedToken.TTL) * time.Second
	if ttl > 0 {
		_ = u.tokenCache.SetUserID(ctx, accessToken, verifiedToken.UserID, ttl)
	}

	return verifiedToken.UserID, nil
}
