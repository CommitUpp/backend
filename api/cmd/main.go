package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	groupusecase "github.com/CommitUpp/backend/api/application/usecase/group"
	"github.com/CommitUpp/backend/api/application/usecase/movie"
	"github.com/CommitUpp/backend/api/application/usecase/user"
	"github.com/CommitUpp/backend/api/infrastructure"
	"github.com/CommitUpp/backend/api/infrastructure/grpc"
	"github.com/CommitUpp/backend/api/infrastructure/postgres"
	"github.com/CommitUpp/backend/api/interfaces/grpc/pb"
	"github.com/CommitUpp/backend/api/interfaces/handler"
	"github.com/CommitUpp/backend/api/interfaces/router"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()

	authAddr := os.Getenv("AUTH_SERVICE_ADDR")
	if authAddr == "" {
		authAddr = "backend-auth:50051"
	}

	dbURL := strings.Trim(os.Getenv("SUPABASE_DB_URL"), `"'`)
	if dbURL == "" {
		log.Fatal("SUPABASE_DB_URL is required")
	}

	dbConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("failed to parse database config: %v", err)
	}
	dbConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	dbPool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		log.Fatalf("failed to create database pool: %v", err)
	}
	defer dbPool.Close()
	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	conn, err := ggrpc.NewClient(authAddr, ggrpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	pbClient := pb.NewAuthServiceClient(conn)
	authGateway := grpc.NewAuthGateway(pbClient)
	authUsecase := user.NewAuthUsecase(authGateway)
	groupRepository := postgres.NewGroupRepository(dbPool)
	groupUsecase := groupusecase.NewGroupUsecase(groupRepository)
	authHandler := handler.NewAuthHandler(authUsecase)
	groupHandler := handler.NewGroupHandler(groupUsecase)

	supabaseURL := strings.TrimRight(os.Getenv("PUBLIC_SUPABASE_URL"), "/")
	if supabaseURL == "" {
		log.Fatal("PUBLIC_SUPABASE_URL is required")
	}

	supabaseAnonKey := os.Getenv("SUPABASE_ANON_KEY")
	if supabaseAnonKey == "" {
		log.Fatal("SUPABASE_ANON_KEY is required")
	}

	movieStatusRepository := infrastructure.NewMovieStatusRepository(supabaseURL, supabaseAnonKey)
	movieStatusUsecase := movie.NewMovieStatusUsecase(movieStatusRepository)
	movieStatusHandler := handler.NewMovieStatusHandler(movieStatusUsecase)

	server := handler.NewServer(authHandler, movieStatusHandler, groupHandler)

	routerConfig := router.RouterConfig{
		AuthUsecase: authUsecase,
		Server:      server,
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
