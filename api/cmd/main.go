package main

import (
	"log"
	"os"

	"github.com/CommitUpp/backend/api/application/usecase/user"
	"github.com/CommitUpp/backend/api/infrastructure/grpc"
	"github.com/CommitUpp/backend/api/interfaces/grpc/pb"
	"github.com/CommitUpp/backend/api/interfaces/handler"
	"github.com/CommitUpp/backend/api/interfaces/router"

	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "success",
		})
	})

	e.Use(handler.MetricsMiddleware())
	// Exposes API HTTP request count, latency, and in-flight request metrics for Prometheus.
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	log.Println("Goサーバーがポート 8080 で起動しました。")
	log.Fatal(e.Start(":8080"))
}
