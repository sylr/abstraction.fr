package http

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	_ "net/http/pprof" // Registering pprof
	"os"
	"path/filepath"
	"strings"

	"abstraction.fr/config"
	"abstraction.fr/pkg/http/handlers/errorxxx"
	"abstraction.fr/pkg/http/handlers/goget"
	"abstraction.fr/pkg/http/handlers/lookingglass"
	"abstraction.fr/pkg/http/handlers/resume"
	"abstraction.fr/pkg/http/handlers/static"
	"abstraction.fr/pkg/http/middlewares/accesscontrol"
	"abstraction.fr/pkg/http/middlewares/headers"
	"abstraction.fr/pkg/http/middlewares/metrics"
	tstatic "abstraction.fr/static"
	"abstraction.fr/templates"

	"github.com/Masterminds/sprig/v3"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

const (
	staticPrefix = "/static"
)

// NewHTTPRouter returns an http.Handler based on given configuration.
// The given router should replace the one which had the previous configuration.
func NewHTTPRouter(conf *config.Config, logger *zap.Logger) http.Handler {
	router := mux.NewRouter()

	// Templates
	t := template.New("abstraction.fr").Funcs(sprig.FuncMap())

	var tpl *template.Template
	if len(conf.TemplatesDir) > 0 {
		err := filepath.Walk(conf.TemplatesDir, func(path string, info os.FileInfo, err error) error {
			if strings.Contains(path, ".html") {
				_, err = t.ParseFiles(path)
				if err != nil {
					logger.Error(fmt.Sprintf("walking %s failed", path), zap.Error(err))
				}
			}

			return err
		})

		tpl = template.Must(t, err)
	} else {
		tpl = template.Must(t.ParseFS(templates.StaticFS, "*.html"))
	}

	// Static
	var httpFS http.FileSystem
	if len(conf.StaticDir) > 0 {
		httpFS = http.Dir(conf.StaticDir)
	} else {
		httpFS = http.FS(tstatic.StaticFS)
	}

	// Handlers
	fsHandler := http.FileServer(httpFS)
	staticfsHandler := http.StripPrefix(staticPrefix, fsHandler)
	staticHandler := static.NewHandler(conf, logger, staticfsHandler, *config.Version)
	faviconHandler := static.NewHandler(conf, logger, fsHandler, *config.Version)
	resumeHandler := resume.NewHandler(conf, logger, tpl)
	// unavailableHandler := unavailable.NewHandler(conf, logger, staticHandler)
	forbiddenHandler := errorxxx.NewHandler(conf, logger, tpl, &errorxxx.Data{StatusCode: 403, Message: "/!\\ Forbidden"})
	notfoundHandler := errorxxx.NewHandler(conf, logger, tpl, &errorxxx.Data{StatusCode: 404, Message: "/!\\ Not Found"})
	lookingglassHandler := lookingglass.NewHandler(conf, logger, tpl)
	gogetHandler := goget.NewHandler(conf, logger, tpl)

	// Middlewares
	serverHeaderMiddleware := headers.NewServerMiddleware(conf, logger, *config.Version)
	probesWhitelistMiddleware := accesscontrol.NewIPWhitelistMiddleware(conf, logger, forbiddenHandler)
	metricsMiddleware := metrics.NewMiddleware(conf, logger)
	defaultMiddlewares := []mux.MiddlewareFunc{
		metricsMiddleware.Middleware,
		serverHeaderMiddleware.Middleware,
	}

	// Liveness
	subrouter := router.Path("/ping").Subrouter()
	subrouter.NewRoute().HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})
	subrouter.Use(serverHeaderMiddleware.Middleware, probesWhitelistMiddleware.Middleware)

	// Readiness
	subrouter = router.Path("/ready").Subrouter()
	subrouter.NewRoute().HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})
	subrouter.Use(serverHeaderMiddleware.Middleware, probesWhitelistMiddleware.Middleware)

	// Metrics
	subrouter = router.Path("/metrics").Subrouter()
	subrouter.NewRoute().Handler(promhttp.Handler())
	subrouter.Use(serverHeaderMiddleware.Middleware, probesWhitelistMiddleware.Middleware)

	// Profiling
	subrouter = router.PathPrefix("/debug/pprof/").Subrouter()
	subrouter.NewRoute().Handler(http.DefaultServeMux)
	subrouter.Use(serverHeaderMiddleware.Middleware, probesWhitelistMiddleware.Middleware)

	// favicon.ico
	subrouter = router.Path("/favicon.ico").Subrouter()
	subrouter.NewRoute().Handler(faviconHandler)
	subrouter.Use(serverHeaderMiddleware.Middleware)

	// Static content
	subrouter = router.PathPrefix(staticPrefix).Subrouter()
	subrouter.NewRoute().Handler(staticHandler)
	subrouter.Use(serverHeaderMiddleware.Middleware)

	// Looking Glass
	subrouter = router.Path("/lg").Subrouter()
	subrouter.Use(defaultMiddlewares...)
	subrouter.NewRoute().Handler(lookingglassHandler)

	// Go Get
	subrouter = router.PathPrefix("/").Subrouter()
	subrouter.Queries("go-get", "1").Handler(gogetHandler)

	// Index
	subrouter = router.Path("/").Subrouter()
	subrouter.Use(defaultMiddlewares...)
	subrouter.NewRoute().Handler(resumeHandler)

	// Error pages router
	subrouter = router.PathPrefix("/").Subrouter()
	subrouter.Use(defaultMiddlewares...)

	// Specific headers routing
	subrouter.Headers("X-Error-Code", "403").Handler(forbiddenHandler)
	subrouter.Headers("X-Error-Code", "404").Handler(notfoundHandler)

	// Specific endpoints routing
	subrouter.Path("/403").Handler(forbiddenHandler)
	subrouter.Path("/404").Handler(notfoundHandler)

	// Default
	subrouter.NewRoute().Handler(notfoundHandler)

	return router
}
