package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// TODO: Optional Buckets.

// DefaultInstrumentation TODO
type DefaultInstrumentation struct {
	RequestDuration *prometheus.HistogramVec
	RequestSize     *prometheus.SummaryVec
	RequestsTotal   *prometheus.CounterVec
	ResponseSize    *prometheus.SummaryVec
}

// New TODO
func New(reg prometheus.Registerer, name string) *DefaultInstrumentation {
	return &DefaultInstrumentation{
		RequestDuration: promauto.With(reg).NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        "http_request_duration_seconds",
				Help:        "Tracks the latencies for HTTP requests.",
				ConstLabels: prometheus.Labels{"server": name},
				Buckets:     []float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120},
			},
			[]string{"code", "handler", "method"},
		),
		RequestSize: promauto.With(reg).NewSummaryVec(
			prometheus.SummaryOpts{
				Name:        "http_request_size_bytes",
				Help:        "Tracks the size of HTTP requests.",
				ConstLabels: prometheus.Labels{"server": name},
			},
			[]string{"code", "handler", "method"},
		),
		RequestsTotal: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Name:        "http_requests_total",
				Help:        "Tracks the number of HTTP requests.",
				ConstLabels: prometheus.Labels{"server": name},
			}, []string{"code", "handler", "method"},
		),
		ResponseSize: promauto.With(reg).NewSummaryVec(
			prometheus.SummaryOpts{
				Name:        "http_response_size_bytes",
				Help:        "Tracks the size of HTTP responses.",
				ConstLabels: prometheus.Labels{"server": name},
			},
			[]string{"code", "handler", "method"},
		),
	}
}
