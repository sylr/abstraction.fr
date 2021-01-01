package errorxxx

import (
	"bufio"
	"bytes"
	"html/template"
	"net/http"

	"abstraction.fr/config"

	ua "github.com/mileusna/useragent"
	log "github.com/sirupsen/logrus"
)

// Handler ...
type Handler struct {
	config *config.Config
	logger *log.Logger
	tpl    *template.Template
	data   *Data
}

// Data ...
type Data struct {
	StatusCode int
	Message    string
}

// NewHandler ...
func NewHandler(conf *config.Config, logger *log.Logger, tpl *template.Template, data *Data) *Handler {
	if data == nil {
		panic("data should not be nil")
	}

	h := Handler{
		config: conf,
		logger: logger,
		tpl:    tpl,
		data:   data,
	}

	return &h
}

// ServeHTTP ...
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ua := ua.Parse(r.UserAgent())

	switch {
	case ua.Desktop, ua.Tablet, ua.Mobile:
		buf := bytes.NewBuffer(nil)
		wr := bufio.NewWriter(buf)

		err := h.tpl.ExecuteTemplate(wr, "errorxxx", h.data)
		if err != nil {
			w.WriteHeader(h.data.StatusCode)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.Write([]byte(err.Error()))
			return
		}

		wr.Flush()

		w.WriteHeader(h.data.StatusCode)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(buf.Bytes())
	default:
		w.WriteHeader(h.data.StatusCode)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(h.data.Message))
	}
}
