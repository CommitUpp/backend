package main

import (
	"log"
	"os"

	"github.com/CommitUpp/backend/api/application/usecase/user"
	"github.com/CommitUpp/backend/api/infrastructure/grpc"
	"github.com/CommitUpp/backend/api/interfaces/grpc/pb"
	"github.com/CommitUpp/backend/api/interfaces/handler"
	"github.com/CommitUpp/backend/api/interfaces/router"

	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	authAddr := os.Getenv("AUTH_SERVICE_ADDR")
	if authAddr == "" {
		authAddr = "backend-auth:50051"
	}

	conn, err := ggrpc.NewClient(authAddr, ggrpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	pbClient := pb.NewAuthServiceClient(conn)
	authGateway := grpc.NewAuthGateway(pbClient)
	authUsecase := user.NewAuthUsecase(authGateway)
	authHandler := handler.NewAuthHandler(authUsecase)

	routerConfig := router.RouterConfig{
		AuthHandler: authHandler,
	}

	e := router.NewRouter(routerConfig)

	log.Println("Goサーバーがポート 8080 で起動しました。")
	log.Fatal(e.Start(":8080"))
}
