package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	middlewarechi "github.com/go-chi/chi/middleware"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"

	middleware "github.com/kakkoyun/middleware/chi"
	"github.com/kakkoyun/middleware/logger"
)

func main() {
	l := logger.NewLogger("debug", logger.LogFormatLogfmt, "example")
	defer level.Info(l).Log("msg", "exiting")

	reg := prometheus.NewRegistry()
	reg.MustRegister(
		version.NewCollector("middleware_chi_example"),
		prometheus.NewGoCollector(),
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
	)

	r := chi.NewRouter()
	r.Use(middlewarechi.RequestID)
	r.Use(middlewarechi.RealIP)
	r.Use(middlewarechi.Recoverer)
	r.Use(middlewarechi.StripSlashes)
	r.Use(middleware.Logger(l))
	r.Use(middleware.Instrumentation(reg, "middleware_chi_example"))

	r.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	r.Get("/200", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	r.Get("/401", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})
	r.Group(func(r chi.Router) {
		r.Get("/api/200", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		r.Get("/api/401", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		})
	})

	level.Info(l).Log("msg", "listening")
	log.Fatal(http.ListenAndServe(":8080", r))
}
