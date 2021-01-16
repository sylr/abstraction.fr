package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"abstraction.fr/config"
	www "abstraction.fr/pkg/http"
	"abstraction.fr/pkg/http/handlers/safewrapper"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	qdconfig "sylr.dev/libqd/config"
)

var (
	version   = "v0.0.0"
	goVersion = runtime.Version()
)

var (
	wwwBuildInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "abstraction",
			Subsystem: "fr",
			Name:      "build_info",
			Help:      "abstraction.fr build info",
		},
		[]string{"version"},
	)
)

func init() {
	// Set & Register build info metric
	wwwBuildInfo.WithLabelValues(version).Set(1)
	prometheus.MustRegister(wwwBuildInfo)
}

func main() {
	// looping for --version in args
	for _, val := range os.Args {
		if val == "--version" {
			fmt.Printf("abstraction.fr version %s\n", version)
			os.Exit(0)
		} else if val == "--" {
			break
		}
	}

	// Logger
	atomlevel := zap.NewAtomicLevel()
	atomlevel.SetLevel(zap.DebugLevel)
	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.EncoderConfig.ConsoleSeparator = " "
	zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := zapConfig.Build()
	logger.Info(fmt.Sprintf("abstraction.fr version %s", version))
	qdlogger := Logger{logger.WithOptions(zap.AddCallerSkip(1))}

	// Version
	config.Version = &version

	// Config default values
	conf := &config.Config{}

	ctx := context.Background()

	// Configuration
	configManager := qdconfig.GetManager(&qdlogger)

	// mutex to prevent data races around conf
	mu := config.NewMutex(configManager, logger, &atomlevel)

	// Add a validator/applier functions
	configManager.AddValidators(nil, mu.ConfigValidator)
	configManager.AddAppliers(nil, mu.ConfigApplier)

	// Make the config
	err := configManager.MakeConfig(ctx, nil, conf)

	if err != nil {
		logger.Fatal("", zap.Error(err))
	}

	// Add templates and static dirs to config watcher
	// It allows to create a new router with updated templates
	for _, path := range []string{conf.TemplatesDir, conf.StaticDir} {
		if len(path) == 0 {
			continue
		}

		filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				err = configManager.GetWatcher(nil).Add(path)
				if err != nil {
					logger.Fatal("", zap.Error(err))
				}
			}

			return nil
		})
	}

	confChan := configManager.NewConfigChan(nil)

	// HTTP router
	router := www.NewHTTPRouter(conf, logger)
	safehandler := safewrapper.New(router)

	// HTTP Server
	server := http.Server{
		Handler:      safehandler,
		Addr:         fmt.Sprintf("%s:%d", conf.ListeningAddress, conf.ListeningPort),
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  30 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	go server.ListenAndServe()

	if len(conf.UnixSocket) > 0 {
		os.Remove(conf.UnixSocket)
		unixListener, err := net.Listen("unix", conf.UnixSocket)
		defer os.Remove(conf.UnixSocket)
		if err != nil {
			panic(err)
		}

		go server.Serve(unixListener)
	}

	// Replace router when new conf is sent through the config chan
	for {
		select {
		case newConf := <-confChan:
			newRouter := www.NewHTTPRouter(newConf.(*config.Config), logger)
			safehandler.SwapHandler(newRouter)
		}
	}
}
