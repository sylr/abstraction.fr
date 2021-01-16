package config

import (
	"github.com/prometheus/client_golang/prometheus"
)

//go:generate deepcopy-gen --input-dirs . --output-package . --output-file-base config_deepcopy --go-header-file /dev/null

var (
	// Version ...
	Version *string
)

var (
	wwwConfigReload = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "abstraction",
			Subsystem: "www",
			Name:      "config_reloads_total",
			Help:      "Total number of config reloads",
		},
		[]string{},
	)
)

func init() {
	prometheus.MustRegister(wwwConfigReload)
}

// Config ...
// +k8s:deepcopy-gen:interfaces=sylr.dev/libqd/config.Config
// +k8s:deepcopy-gen=true
type Config struct {
	File             string          `                         short:"f" long:"config"  description:"Yaml config"`
	Verbose          []bool          `yaml:"verbose"           short:"v" long:"verbose" description:"Show verbose debug information"`
	Version          bool            `                                   long:"version" description:"Show version"`
	ListeningAddress string          `yaml:"listening_address" short:"a" long:"address" description:"Listening address" default:"127.0.0.1"`
	ListeningPort    int             `yaml:"listening_port"    short:"p" long:"port"    description:"Listening port" default:"8080"`
	UnixSocket       string          `yaml:"unix_socket"       short:"u" long:"unix"    description:"Listening unix socket"`
	Internal         *InternalConfig `yaml:"internal"`
	TemplatesDir     string          `yaml:"templates"`
	StaticDir        string          `yaml:"static"`
	GoModules        []GoModule      `yaml:"go-modules"`
	Reloads          int
}

// InternalConfig ...
// +k8s:deepcopy-gen=true
type InternalConfig struct {
	TrustedSubnets []string `yaml:"trusted_subnets"`
}

// GoModule ...
// +k8s:deepcopy-gen=true
type GoModule struct {
	Name   string `yaml:"name"`
	Import string `yaml:"go-import"`
	Source string `yaml:"go-source"`
}

// ConfigFile ...
func (c *Config) ConfigFile() string {
	return c.File
}
