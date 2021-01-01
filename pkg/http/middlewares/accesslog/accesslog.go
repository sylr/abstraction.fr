package accesslog

import (
	"fmt"
	"net/http"
	"runtime/trace"

	"abstraction.fr/config"

	log "github.com/sirupsen/logrus"
)

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
		"_package": "middlewares.accesslog",
		"_self":    fmt.Sprintf("%p", &mdw),
	})

	return &mdw
}

// Middleware ...
func (mdw *Middleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tr := trace.StartRegion(r.Context(), "AccesslogMiddleware")
		defer tr.End()

		func() {
			defer recover()

			mdw.Logger.WithFields(log.Fields{
				"_func": "Middleware",
				"host":  r.Host,
				"path":  r.URL.Path,
			}).Info()
		}()

		next.ServeHTTP(w, r)
	})
}
