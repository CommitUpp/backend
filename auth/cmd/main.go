package main

import (
	"log"
	"net"
	"os"

	"github.com/CommitUpp/backend/auth/application/usecase"
	authredis "github.com/CommitUpp/backend/auth/infrastructure/redis"
	authsupabase "github.com/CommitUpp/backend/auth/infrastructure/supabase"
	"github.com/CommitUpp/backend/auth/interfaces/grpc/pb"
	"github.com/CommitUpp/backend/auth/interfaces/handler"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

func main() {
	e := echo.New()
	// Exposes auth gRPC, Redis token cache, and Supabase token verify metrics for Prometheus.
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	go func() {
		log.Println("auth metrics server listening on :25000")
		if err := e.Start(":25000"); err != nil {
			log.Fatalf("failed to serve auth metrics: %v", err)
		}
	}()

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "redis-cache:6379"
	}

	supabaseURL := os.Getenv("PUBLIC_SUPABASE_URL")
	if supabaseURL == "" {
		log.Fatal("PUBLIC_SUPABASE_URL is required")
	}

	supabaseAnonKey := os.Getenv("SUPABASE_ANON_KEY")
	if supabaseAnonKey == "" {
		log.Fatal("SUPABASE_ANON_KEY is required")
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	tokenCache := authredis.NewTokenCacheRepository(redisClient)
	tokenVerifier := authsupabase.NewTokenVerifierRepository(
		supabaseURL,
		supabaseAnonKey,
	)
	authUsecase := usecase.NewAuthUsecase(tokenCache, tokenVerifier)
	authHandler := handler.NewAuthGRPCHandler(authUsecase)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(handler.GRPCMetricsInterceptor()),
	)
	pb.RegisterAuthServiceServer(server, authHandler)

	log.Println("auth gRPC server listening on :50051")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
