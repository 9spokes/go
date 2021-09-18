package logging

import (
	"fmt"
	"strings"
)

const (
	LevelDebug = iota
	LevelInfo
	LevelWarning
	LevelError
	LevelFatal
)

const (
	ColorDebug   = "\033[2;34m%s\033[0m"
	ColorInfo    = "\033[1;97m%s\033[0m"
	ColorWarning = "\033[1;33m%s\033[0m"
	ColorError   = "\033[1;31m%s\033[0m"
	ColorFatal   = "\033[1;35m%s\033[0m"
)

type Level int

type Logger struct {
	level  Level
	logger provider
}

type provider interface {
	log(int, string)
}

// New ceates a new 9 Spokes logger with a formatter, uses stdout as a backend.
func New(level string, syslogd, env string) *Logger {

	var l Level

	switch strings.ToLower(level) {
	case "debug":
		l = LevelDebug
	case "warning":
		l = LevelWarning
	case "error":
		l = LevelError
	case "fatal":
		l = LevelFatal
	default:
		l = LevelInfo
	}

	if syslogd != "" {
		fmt.Printf("\n\n*** Live logs from this service are being sent to https://my.papertrailapp.com ***\n\n")
		return &Logger{
			logger: NewSyslog(syslogd, env),
			level:  l,
		}
	}

	return &Logger{
		logger: NewOplogging("", level),
		level:  l,
	}
}

func (l *Logger) Debug(e string) {
	if l.level > LevelDebug {
		return
	}
	l.logger.log(LevelDebug, e)
}

func (l *Logger) Debugf(e string, args ...interface{}) {
	if l.level > LevelDebug {
		return
	}
	l.logger.log(LevelDebug, fmt.Sprintf(e, args...))
}

func (l *Logger) Info(e string) {
	if l.level > LevelInfo {
		return
	}
	l.logger.log(LevelInfo, e)
}

func (l *Logger) Infof(e string, args ...interface{}) {
	if l.level > LevelInfo {
		return
	}
	l.logger.log(LevelInfo, fmt.Sprintf(e, args...))
}

func (l *Logger) Error(e string) {
	if l.level > LevelError {
		return
	}
	l.logger.log(LevelError, e)
}

func (l *Logger) Errorf(e string, args ...interface{}) {
	if l.level > LevelError {
		return
	}
	l.logger.log(LevelError, fmt.Sprintf(e, args...))
}

func (l *Logger) Warning(e string) {
	if l.level > LevelWarning {
		return
	}
	l.logger.log(LevelWarning, e)
}

func (l *Logger) Warningf(e string, args ...interface{}) {
	if l.level > LevelWarning {
		return
	}
	l.logger.log(LevelWarning, fmt.Sprintf(e, args...))
}

func (l *Logger) Fatal(e string) {
	l.logger.log(LevelFatal, e)
}

func (l *Logger) Fatalf(e string, args ...interface{}) {
	l.logger.log(LevelFatal, fmt.Sprintf(e, args...))
}
