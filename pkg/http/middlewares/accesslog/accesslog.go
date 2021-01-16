package accesslog

import (
	"net/http"
	"runtime/trace"

	"abstraction.fr/config"

	"go.uber.org/zap"
)

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
		tr := trace.StartRegion(r.Context(), "AccesslogMiddleware")
		defer tr.End()

		func() {
			defer recover()

			mdw.Logger.Info(
				"",
				zap.String("_package", "middlewares.accesslog"),
				zap.String("_func", "Middleware"),
				zap.String("host", r.Host),
				zap.String("path", r.URL.Path),
			)
		}()

		next.ServeHTTP(w, r)
	})
}
