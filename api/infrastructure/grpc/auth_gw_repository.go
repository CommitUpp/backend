package grpc

import (
	"context"
	"errors"

	"github.com/CommitUpp/backend/api/domain/repository"
	"github.com/CommitUpp/backend/api/interfaces/grpc/pb"
)

type authGatewayImpl struct {
	client pb.AuthServiceClient
}

func NewAuthGateway(client pb.AuthServiceClient) repository.AuthGateway {
	return &authGatewayImpl{
		client: client,
	}
}

func (g *authGatewayImpl) VerifyToken(ctx context.Context, accessToken string) (string, error) {
	req := &pb.VerifyTokenRequest{
		AccessToken: accessToken,
	}

	res, err := g.client.VerifyToken(ctx, req)
	if err != nil {
		return "", err
	}

	if res.GetUserId() == "" {
		return "", errors.New("unauthorized: returned user_id is empty")
	}

	return res.GetUserId(), nil
}
