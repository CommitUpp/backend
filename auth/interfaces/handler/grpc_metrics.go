package handler

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	// Counts auth gRPC requests by RPC method and gRPC status code.
	grpcRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_grpc_requests_total",
			Help: "Total number of gRPC requests handled by the auth service.",
		},
		[]string{"method", "code"},
	)

	// Measures auth gRPC request latency by RPC method and gRPC status code.
	grpcRequestDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "auth_grpc_request_duration_seconds",
			Help:    "Auth gRPC request latency in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "code"},
	)

	// Tracks the current number of auth gRPC requests being processed.
	grpcInFlightRequests = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "auth_grpc_in_flight_requests",
			Help: "Current number of gRPC requests being handled by the auth service.",
		},
	)
)

func GRPCMetricsInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		startedAt := time.Now()
		grpcInFlightRequests.Inc()
		defer grpcInFlightRequests.Dec()

		res, err := handler(ctx, req)
		code := status.Code(err).String()

		grpcRequestsTotal.WithLabelValues(info.FullMethod, code).Inc()
		grpcRequestDurationSeconds.WithLabelValues(info.FullMethod, code).Observe(time.Since(startedAt).Seconds())

		return res, err
	}
}
