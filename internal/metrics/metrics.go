package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Define a counter metric
	RequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cryptkeeper_request_count",
			Help: "Total number of requests",
		},
		[]string{"path"},
	)

	// Define a histogram metric
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cryptkeeper_request_duration_seconds",
			Help:    "Duration of requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)
)

func init() {
	// Register the metrics with Prometheus's default registry
	prometheus.MustRegister(RequestCount)
	prometheus.MustRegister(RequestDuration)
}
