package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Counts API HTTP requests by method, route path, and HTTP status code.
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_http_requests_total",
			Help: "Total number of HTTP requests handled by the API.",
		},
		[]string{"method", "path", "status"},
	)

	// Measures API HTTP request latency by method, route path, and HTTP status code.
	httpRequestDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "api_http_request_duration_seconds",
			Help:    "HTTP request latency in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	// Tracks the current number of API HTTP requests being processed.
	httpInFlightRequests = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "api_http_in_flight_requests",
			Help: "Current number of HTTP requests being handled by the API.",
		},
	)
)

func MetricsMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Path() == "/metrics" {
				return next(c)
			}

			startedAt := time.Now()
			httpInFlightRequests.Inc()
			defer httpInFlightRequests.Dec()

			err := next(c)
			status := responseStatus(c, err)
			path := c.Path()
			if path == "" {
				path = c.Request().URL.Path
			}

			labels := []string{
				c.Request().Method,
				path,
				strconv.Itoa(status),
			}

			httpRequestsTotal.WithLabelValues(labels...).Inc()
			httpRequestDurationSeconds.WithLabelValues(labels...).Observe(time.Since(startedAt).Seconds())

			return err
		}
	}
}

func responseStatus(c echo.Context, err error) int {
	if err == nil {
		status := c.Response().Status
		if status == 0 {
			return http.StatusOK
		}

		return status
	}

	httpErr, ok := err.(*echo.HTTPError)
	if ok {
		return httpErr.Code
	}

	return http.StatusInternalServerError
}
