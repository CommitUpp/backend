package user

import (
	"context"
	"github.com/CommitUpp/backend/api/domain/repository"
)

type AuthUsecase interface {
	LoginCallback(ctx context.Context, accessToken string) (string, error)
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
	userID, err := u.authGateway.VerifyToken(ctx, accessToken)
	if err != nil {
		return "", err
	}

	return userID, nil
}
