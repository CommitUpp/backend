package handler

import (
	"context"

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

func (h *authGRPCHandler) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) {

}
