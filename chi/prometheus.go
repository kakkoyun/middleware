package chi

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"

	middlewareprom "github.com/kakkoyun/middleware/prometheus"
)

// Instrumentation returns a new Prometheus middleware for Chi.
func Instrumentation(reg prometheus.Registerer, name string) func(next http.Handler) http.Handler {
	ins := middlewareprom.New(reg, name)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)

			// https://godoc.org/github.com/go-chi/chi#Context.RoutePattern
			routePattern := chi.RouteContext(r.Context()).RoutePattern()

			ins.RequestDuration.WithLabelValues(http.StatusText(ww.Status()), routePattern, r.Method).Observe(time.Since(start).Seconds())
			ins.RequestSize.WithLabelValues(http.StatusText(ww.Status()), routePattern, r.Method).Observe(float64(computeApproximateRequestSize(r)))
			ins.RequestsTotal.WithLabelValues(http.StatusText(ww.Status()), routePattern, r.Method).Inc()
			ins.ResponseSize.WithLabelValues(http.StatusText(ww.Status()), routePattern, r.Method).Observe(float64(ww.BytesWritten()))
		})
	}
}

// https://github.com/prometheus/client_golang/blob/06b1a0a6ae29dd8b39953dc7f1954a0b2fd680be/prometheus/promhttp/instrument_server.go#L298:6
func computeApproximateRequestSize(r *http.Request) int {
	s := 0
	if r.URL != nil {
		s += len(r.URL.String())
	}

	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)

	// N.B. r.Form and r.MultipartForm are assumed to be included in r.URL.

	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}
	return s
}
