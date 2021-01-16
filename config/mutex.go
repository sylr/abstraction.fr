package config

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
	qdconfig "sylr.dev/libqd/config"
)

// NewMutex ...
func NewMutex(manager *qdconfig.Manager, l *log.Logger) Mutex {
	return Mutex{
		&sync.RWMutex{},
		manager,
		l,
	}
}

// Mutex ...
type Mutex struct {
	*sync.RWMutex
	manager *qdconfig.Manager
	l       *log.Logger
}

// ConfigValidator ...
func (cm *Mutex) ConfigValidator(currentConfig qdconfig.Config, newConfig qdconfig.Config) []error {
	var currentConf *Config
	var newConf *Config
	var ok bool

	// currentConfig is nil the first time the validator is called
	if currentConfig != nil {
		currentConf, ok = currentConfig.(*Config)

		if !ok {
			return []error{fmt.Errorf("Can not cast currentConfig to (*Config)")}
		}
	}

	newConf, ok = newConfig.(*Config)

	if !ok {
		return []error{fmt.Errorf("Can not cast newConfig to (*Config)")}
	}

	// ---------------------------------------------------------------------
	// Here begins the actual validation of the values of newConfig
	// ---------------------------------------------------------------------
	var errs []error

	if currentConfig == nil {
		if newConf.ListeningPort < 0 || newConf.ListeningPort > 65535 {
			errs = append(errs, fmt.Errorf("ListeningPort `%d` is not valid", newConf.ListeningPort))
		}
	} else {
		if newConf.ListeningPort != currentConf.ListeningPort {
			errs = append(errs, fmt.Errorf("ListeningPort `%d` can not be changed to `%d`", currentConf.ListeningPort, newConf.ListeningPort))
		}
		if newConf.UnixSocket != currentConf.UnixSocket {
			errs = append(errs, fmt.Errorf("UnixSocket `%s` can not be changed to `%s`", currentConf.UnixSocket, newConf.UnixSocket))
		}
	}

	if len(newConf.Verbose) > 5 {
		errs = append(errs, fmt.Errorf("Too many verbose flags"))
	}

	return errs
}

// ConfigApplier ...
func (cm *Mutex) ConfigApplier(currentConfig qdconfig.Config, newConfig qdconfig.Config) error {
	var currentConf *Config
	var newConf *Config
	var ok bool

	// currentConfig is nil the first time the validator is called
	if currentConfig != nil {
		currentConf, ok = currentConfig.(*Config)

		if !ok {
			return fmt.Errorf("Can not cast currentConfig to (*Config)")
		}
	}

	newConf, ok = newConfig.(*Config)

	if !ok {
		return fmt.Errorf("Can not cast newConfig to (*Config)")
	}

	switch len(newConf.Verbose) {
	case 1:
		log.SetLevel(log.ErrorLevel)
	case 2:
		log.SetLevel(log.WarnLevel)
	case 3:
		log.SetLevel(log.InfoLevel)
	case 4:
		log.SetLevel(log.DebugLevel)
	case 5:
		log.SetLevel(log.TraceLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	if currentConf != nil {
		wwwConfigReload.WithLabelValues().Inc()
	}

	return nil
}
