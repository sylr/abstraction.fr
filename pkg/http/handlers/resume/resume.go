package resume

import (
	"bufio"
	"bytes"
	"html/template"
	"net/http"
	"strings"

	"abstraction.fr/config"

	ua "github.com/mileusna/useragent"
	log "github.com/sirupsen/logrus"
)

// Handler ...
type Handler struct {
	config *config.Config
	logger *log.Logger

	tpl  *template.Template
	html []byte
}

// NewHandler ...
func NewHandler(conf *config.Config, logger *log.Logger, tpl *template.Template) *Handler {
	h := Handler{
		config: conf,
		logger: logger,
		tpl:    tpl,
	}

	var buf bytes.Buffer
	wr := bufio.NewWriter(&buf)
	err := h.tpl.ExecuteTemplate(wr, "resume", nil)

	if err != nil {
		logger.Errorf("handler/resume: %s", err)
	}

	err = wr.Flush()

	if err != nil {
		logger.Errorf("handler/resume: %s", err)
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
		log.Debugf("handler/cv: \"%s\" not recognized", r.UserAgent())
	}
}
