package goget

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strings"

	"abstraction.fr/config"

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
	if len(h.config.GoModules) == 0 {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	cn := fmt.Sprintf("%s%s", r.Host, r.URL.Path)
	cn = strings.TrimRight(cn, "/")

	// Reverse sort
	sort.Slice(h.config.GoModules, func(i, j int) bool {
		return h.config.GoModules[i].Name > h.config.GoModules[j].Name
	})

	for _, mod := range h.config.GoModules {
		if !strings.HasPrefix(cn, mod.Name) {
			continue
		}

		// var suffix string
		// if cn != mod.Name {
		// 	suffix = cn[len(mod.Name):]
		// 	mod.Name = mod.Name + suffix
		// }

		buf := bytes.NewBuffer(nil)
		wr := bufio.NewWriter(buf)
		err := h.tpl.ExecuteTemplate(wr, "go-get", struct{ GoModule config.GoModule }{mod})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.Write([]byte(err.Error()))
			return
		}

		wr.Flush()
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(buf.Bytes())

		return
	}

	w.WriteHeader(http.StatusNotFound)
}
