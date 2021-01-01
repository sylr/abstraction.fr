package headers

import (
	"fmt"
	"net/http"
	"runtime/trace"

	"abstraction.fr/config"
	log "github.com/sirupsen/logrus"
)

// ServerMiddleware ...
type ServerMiddleware struct {
	Config  *config.Config
	Logger  *log.Entry
	Version string
}

// NewServerMiddleware ...
func NewServerMiddleware(conf *config.Config, logger *log.Entry, version string) *ServerMiddleware {
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
