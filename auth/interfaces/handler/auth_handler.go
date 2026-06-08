package handler

import (
	"context"
	"errors"

	"github.com/CommitUpp/backend/auth/application/usecase"
	"github.com/CommitUpp/backend/auth/interfaces/grpc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authGRPCHandler struct {
	pb.UnimplementedAuthServiceServer
	authUsecase usecase.AuthUsecase
}

func NewAuthGRPCHandler(au usecase.AuthUsecase) pb.AuthServiceServer {
	return &authGRPCHandler{
		authUsecase: au,
	}
}

func (h *authGRPCHandler) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	userID, err := h.authUsecase.VerifyToken(ctx, req.GetAccessToken())
	if err != nil {
		if errors.Is(err, usecase.ErrUnauthorized) {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.VerifyTokenResponse{
		UserId:  userID,
		Message: "success",
	}, nil
}
