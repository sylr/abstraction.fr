package resume

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"abstraction.fr/config"

	ua "github.com/mileusna/useragent"
	"go.uber.org/zap"
)

// Handler ...
type Handler struct {
	config *config.Config
	logger *zap.Logger

	tpl  *template.Template
	html []byte
}

// NewHandler ...
func NewHandler(conf *config.Config, logger *zap.Logger, tpl *template.Template) *Handler {
	h := Handler{
		config: conf,
		logger: logger,
		tpl:    tpl,
	}

	var buf bytes.Buffer
	wr := bufio.NewWriter(&buf)
	err := h.tpl.ExecuteTemplate(wr, "resume", nil)

	if err != nil {
		logger.Error("", zap.Error(err))
	}

	err = wr.Flush()

	if err != nil {
		logger.Error("", zap.Error(err))
	}

	h.html = buf.Bytes()

	return &h
}

// ServeHTTP ...
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ua := ua.Parse(r.UserAgent())

	switch {
	case ua.Desktop, ua.Tablet, ua.Mobile:
		fallthrough
	case strings.Contains(r.UserAgent(), "W3C"), strings.Contains(r.UserAgent(), "Validator.nu/LV"):
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(h.html)
	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("Not Implemented: user-agent not recognized.\n"))
		h.logger.Debug(fmt.Sprintf("handler/cv: \"%s\" not recognized", r.UserAgent()))
	}
}
