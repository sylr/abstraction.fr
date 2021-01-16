package unavailable

import (
	"bytes"
	"html/template"
	"net/http"
	"runtime/trace"

	"abstraction.fr/config"
	"abstraction.fr/pkg/http/handlers/static"

	ua "github.com/mileusna/useragent"
	"go.uber.org/zap"
)

// Handler ...
type Handler struct {
	Config *config.Config
	Logger *zap.Logger
	Static *static.Handler

	template *template.Template
	html     []byte
}

// NewHandler ...
func NewHandler(conf *config.Config, logger *zap.Logger, static *static.Handler) *Handler {
	template, err := template.New("unavailable").Parse(string("unavailable"))

	if err != nil {
		logger.Panic("", zap.Error(err))
	}

	h := Handler{
		Config:   conf,
		Logger:   logger,
		Static:   static,
		template: template,
	}

	h.refreshHTML()

	return &h
}

// ServeHTTP ...
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tr := trace.StartRegion(r.Context(), "UnavailableHandler")
	defer tr.End()

	w.WriteHeader(http.StatusServiceUnavailable)

	ua := ua.Parse(r.UserAgent())

	switch {
	case ua.Desktop, ua.Tablet, ua.Mobile:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(h.html)
	default:
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("Oops, service is currently unavailable.\n"))
	}
}

func (h *Handler) refreshHTML() {
	buf := new(bytes.Buffer)
	err := h.template.Execute(buf, h.Config)

	if err != nil {
		h.Logger.Error("", zap.Error(err))
	}

	h.html = buf.Bytes()
}
