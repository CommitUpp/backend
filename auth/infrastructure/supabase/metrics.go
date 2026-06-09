package supabase

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Counts Supabase token verification requests by result and HTTP status code.
	tokenVerifyTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_supabase_token_verify_total",
			Help: "Total number of Supabase token verification requests.",
		},
		[]string{"result", "status"},
	)

	// Measures Supabase token verification latency by result and HTTP status code.
	tokenVerifyDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "auth_supabase_token_verify_duration_seconds",
			Help:    "Supabase token verification latency in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"result", "status"},
	)
)

func observeTokenVerify(result string, statusCode int, startedAt time.Time) {
	status := "none"
	if statusCode > 0 {
		status = strconv.Itoa(statusCode)
	}

	tokenVerifyTotal.WithLabelValues(result, status).Inc()
	tokenVerifyDurationSeconds.WithLabelValues(result, status).Observe(time.Since(startedAt).Seconds())
}
