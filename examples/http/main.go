package main

import (
	"log"
	"net/http"

	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"

	middleware "github.com/kakkoyun/middleware/http"
	"github.com/kakkoyun/middleware/logger"
)

func main() {
	l := logger.NewLogger("debug", logger.LogFormatLogfmt, "example")
	defer level.Info(l).Log("msg", "exiting")

	reg := prometheus.NewRegistry()
	reg.MustRegister(
		version.NewCollector("middleware_http_example"),
		prometheus.NewGoCollector(),
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
	)

	ins := middleware.NewInstrumentation(reg, "middleware_http_example")

	mux := http.NewServeMux()
	mux.Handle("/200", ins.NewHandler("ok", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))
	mux.Handle("/401", ins.NewHandler("unauthorized", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})))
	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	srv := &http.Server{Addr: ":8080", Handler: mux}
	level.Info(l).Log("msg", "listening")
	log.Fatal(srv.ListenAndServe())
}
