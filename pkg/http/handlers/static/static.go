package static

import (
	"fmt"
	"net/http"
	"strings"

	"abstraction.fr/config"

	"go.uber.org/zap"
)

// Handler ...
type Handler struct {
	Config    *config.Config
	Logger    *zap.Logger
	Version   string
	FSHandler http.Handler
}

// NewHandler ...
func NewHandler(conf *config.Config, logger *zap.Logger, handler http.Handler, version string) *Handler {
	h := Handler{
		Config:    conf,
		Logger:    logger,
		Version:   version,
		FSHandler: handler,
	}

	return &h
}

// ServeHTTP ...
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Cache handling
	etag := fmt.Sprintf("\"what/%s/%p\"", h.Version, h.Config)

	if match := r.Header.Get("If-None-Match"); match != "" {
		if strings.Contains(match, etag) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}
	w.Header().Add("Etag", etag)

	h.FSHandler.ServeHTTP(w, r)
}
