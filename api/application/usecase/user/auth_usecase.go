package user

import (
	"context"
	"errors"

	"github.com/CommitUpp/backend/api/domain/repository"
)

var ErrUnauthorized = errors.New("unauthorized")

type AuthUsecase interface {
	LoginCallback(ctx context.Context, accessToken string) (string, error)
	VerifyToken(ctx context.Context, accessToken string) (string, error)
	Logout(ctx context.Context, accessToken string) error
}

type authUsecaseImpl struct {
	authGateway repository.AuthGateway
}

func NewAuthUsecase(ag repository.AuthGateway) AuthUsecase {
	return &authUsecaseImpl{
		authGateway: ag,
	}
}

func (u *authUsecaseImpl) LoginCallback(ctx context.Context, accessToken string) (string, error) {
	return u.VerifyToken(ctx, accessToken)
}

func (u *authUsecaseImpl) VerifyToken(ctx context.Context, accessToken string) (string, error) {
	if accessToken == "" {
		return "", ErrUnauthorized
	}

	userID, err := u.authGateway.VerifyToken(ctx, accessToken)
	if err != nil {
		return "", err
	}

	if userID == "" {
		return "", ErrUnauthorized
	}

	return userID, nil
}

func (u *authUsecaseImpl) Logout(ctx context.Context, accessToken string) error {
	if accessToken == "" {
		return ErrUnauthorized
	}

	return u.authGateway.Logout(ctx, accessToken)
}
