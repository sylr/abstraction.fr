package main

import (
	"fmt"

	"go.uber.org/zap"
)

// Logger ...
type Logger struct {
	logger *zap.Logger
}

// Tracef ...
func (l *Logger) Tracef(format string, vals ...interface{}) {
	l.logger.Debug(fmt.Sprintf(format, vals...))
}

// Debugf ...
func (l *Logger) Debugf(format string, vals ...interface{}) {
	l.logger.Debug(fmt.Sprintf(format, vals...))
}

// Infof ...
func (l *Logger) Infof(format string, vals ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, vals...))
}

// Warnf ...
func (l *Logger) Warnf(format string, vals ...interface{}) {
	l.logger.Warn(fmt.Sprintf(format, vals...))
}

// Errorf ...
func (l *Logger) Errorf(format string, vals ...interface{}) {
	l.logger.Error(fmt.Sprintf(format, vals...))
}

// Fatalf ...
func (l *Logger) Fatalf(format string, vals ...interface{}) {
	l.logger.Fatal(fmt.Sprintf(format, vals...))
}
