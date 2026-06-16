package redis

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Counts Redis token cache get/set operations by operation type and result.
	tokenCacheOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_redis_token_cache_operations_total",
			Help: "Total number of Redis token cache operations.",
		},
		[]string{"operation", "result"},
	)

	// Measures Redis token cache get/set latency by operation type and result.
	tokenCacheOperationDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "auth_redis_token_cache_operation_duration_seconds",
			Help:    "Redis token cache operation latency in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "result"},
	)
)

func observeTokenCacheOperation(operation string, result string, startedAt time.Time) {
	tokenCacheOperationsTotal.WithLabelValues(operation, result).Inc()
	tokenCacheOperationDurationSeconds.WithLabelValues(operation, result).Observe(time.Since(startedAt).Seconds())
}

func redisSetResult(err error) string {
	if err != nil {
		return "error"
	}

	return "success"
}
