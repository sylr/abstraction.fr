package main

import (
	log "github.com/sirupsen/logrus"
)

// Logger ...
type Logger struct {
	logger *log.Logger
}

// Tracef ...
func (l *Logger) Tracef(format string, vals ...interface{}) {
	l.logger.Tracef("go-libqd/config: "+format, vals...)
}

// Debugf ...
func (l *Logger) Debugf(format string, vals ...interface{}) {
	l.logger.Debugf("go-libqd/config: "+format, vals...)
}

// Infof ...
func (l *Logger) Infof(format string, vals ...interface{}) {
	l.logger.Infof("go-libqd/config: "+format, vals...)
}

// Warnf ...
func (l *Logger) Warnf(format string, vals ...interface{}) {
	l.logger.Warnf("go-libqd/config: "+format, vals...)
}

// Errorf ...
func (l *Logger) Errorf(format string, vals ...interface{}) {
	l.logger.Errorf("go-libqd/config: "+format, vals...)
}

// Fatalf ...
func (l *Logger) Fatalf(format string, vals ...interface{}) {
	l.logger.Fatalf("go-libqd/config: "+format, vals...)
}
