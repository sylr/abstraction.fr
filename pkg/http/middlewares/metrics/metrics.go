package metrics

import (
	"net/http"
	"runtime/trace"

	"abstraction.fr/config"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

var (
	httpRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "abstraction",
			Subsystem: "www",
			Name:      "http_request_total",
			Help:      "Total number of http requests received",
		},
		[]string{"host", "user_agent", "method"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestTotal)
}

// Middleware ...
type Middleware struct {
	Config *config.Config
	Logger *zap.Logger
}

// NewMiddleware ...
func NewMiddleware(conf *config.Config, logger *zap.Logger) *Middleware {
	mdw := Middleware{
		Config: conf,
		Logger: logger,
	}

	return &mdw
}

// Middleware ...
func (mdw *Middleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tr := trace.StartRegion(r.Context(), "MetricsMiddleware")
		defer tr.End()

		func() {
			defer recover()

			r.URL.Host = r.Host
			httpRequestTotal.WithLabelValues(
				r.URL.Hostname(),
				r.UserAgent(),
				r.Method,
			).Inc()
		}()

		next.ServeHTTP(w, r)
	})
}
