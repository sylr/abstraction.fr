package metrics

import (
	"net/http"
	"runtime/trace"

	"abstraction.fr/config"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
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
	Logger *log.Entry
}

// NewMiddleware ...
func NewMiddleware(conf *config.Config, logger *log.Entry) *Middleware {
	mdw := Middleware{
		Config: conf,
	}

	mdw.Logger = logger.WithFields(log.Fields{
		"_package": "middlewares.metrics",
	})

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
