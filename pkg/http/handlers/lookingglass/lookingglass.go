package lookingglass

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"abstraction.fr/config"

	ua "github.com/mileusna/useragent"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
)

// Handler ...
type Handler struct {
	config *config.Config
	logger *log.Logger

	tpl *template.Template
}

// NewHandler ...
func NewHandler(conf *config.Config, logger *log.Logger, tpl *template.Template) *Handler {
	h := Handler{
		config: conf,
		logger: logger,
		tpl:    tpl,
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

		err := h.tpl.ExecuteTemplate(wr, "lookingglass", r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.Write([]byte(err.Error()))
			return
		}

		wr.Flush()

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(buf.Bytes())
	default:
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		h.consoleOutput(w, r)
	}
}

func (h *Handler) consoleOutput(w http.ResponseWriter, r *http.Request) {
	header := fmt.Sprintf("%s %s %s", r.Method, r.RequestURI, r.Proto)
	table := tablewriter.NewWriter(w)

	table.SetAutoFormatHeaders(false)
	table.SetAutoMergeCells(true)
	table.SetHeader([]string{"Request", header})
	table.Append([]string{"Host", r.URL.Host})
	table.Append([]string{"Scheme", r.URL.Scheme})
	table.Append([]string{"Remote Addr", r.RemoteAddr})

	for header, values := range r.Header {
		for _, value := range values {
			table.Append([]string{header, value})
		}
	}

	table.Render()
}
