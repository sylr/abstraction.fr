package headers

import (
	"fmt"
	"net/http"
	"runtime/trace"

	"abstraction.fr/config"

	"go.uber.org/zap"
)

// ServerMiddleware ...
type ServerMiddleware struct {
	Config  *config.Config
	Logger  *zap.Logger
	Version string
}

// NewServerMiddleware ...
func NewServerMiddleware(conf *config.Config, logger *zap.Logger, version string) *ServerMiddleware {
	mdw := ServerMiddleware{
		Config:  conf,
		Logger:  logger,
		Version: version,
	}

	return &mdw
}

// Middleware ...
func (mdw *ServerMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tr := trace.StartRegion(r.Context(), "HeadersServerMiddleware")
		defer tr.End()

		serverString := fmt.Sprintf("abstraction.fr/%s", mdw.Version)
		w.Header().Add("Server", serverString)

		next.ServeHTTP(w, r)
	})
}
