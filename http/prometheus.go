package http

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	middlewareprom "github.com/kakkoyun/middleware/prometheus"
)

// InstrumentationMiddleware TODO
type InstrumentationMiddleware struct {
	ins *middlewareprom.DefaultInstrumentation
}

// NewInstrumentation provides default Instrumentation middleware.
func NewInstrumentation(reg prometheus.Registerer, name string) *InstrumentationMiddleware {
	return &InstrumentationMiddleware{middlewareprom.New(reg, name)}
}

// NewHandler wraps the given HTTP handler for instrumentation.
func (m *InstrumentationMiddleware) NewHandler(handlerName string, handler http.Handler) http.HandlerFunc {
	return promhttp.InstrumentHandlerDuration(
		m.ins.RequestDuration.MustCurryWith(prometheus.Labels{"handler": handlerName}),
		promhttp.InstrumentHandlerRequestSize(
			m.ins.RequestSize.MustCurryWith(prometheus.Labels{"handler": handlerName}),
			promhttp.InstrumentHandlerCounter(
				m.ins.RequestsTotal.MustCurryWith(prometheus.Labels{"handler": handlerName}),
				promhttp.InstrumentHandlerResponseSize(
					m.ins.ResponseSize.MustCurryWith(prometheus.Labels{"handler": handlerName}),
					handler,
				),
			),
		),
	)
}
