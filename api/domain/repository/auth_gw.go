package repository

import "context"

type AuthGateway interface {
	VerifyToken(ctx context.Context, accessToken string) (string, error)
}
