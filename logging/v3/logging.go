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
	level    Level
	provider provider
}

type provider interface {
	do(int, string)
}

func init() {
	New("debug", "", "")
}

var (
	logger Logger
)

func New(level string, syslogd, env string) {

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
		logger = Logger{
			provider: NewSyslog(syslogd, env),
			level:    l,
		}
		return
	}

	logger = Logger{
		provider: NewOplogging("", level),
		level:    l,
	}
	return
}

func Debug(e string) {
	if logger.level > LevelDebug {
		return
	}
	logger.provider.do(LevelDebug, e)
}

func Debugf(e string, args ...interface{}) {
	if logger.level > LevelDebug {
		return
	}
	logger.provider.do(LevelDebug, fmt.Sprintf(e, args...))
}

func Info(e string) {
	if logger.level > LevelInfo {
		return
	}
	logger.provider.do(LevelInfo, e)
}

func Infof(e string, args ...interface{}) {
	if logger.level > LevelInfo {
		return
	}
	logger.provider.do(LevelInfo, fmt.Sprintf(e, args...))
}

func Error(e string) {
	if logger.level > LevelError {
		return
	}
	logger.provider.do(LevelError, e)
}

func Errorf(e string, args ...interface{}) {
	if logger.level > LevelError {
		return
	}
	logger.provider.do(LevelError, fmt.Sprintf(e, args...))
}

func Warning(e string) {
	if logger.level > LevelWarning {
		return
	}
	logger.provider.do(LevelWarning, e)
}

func Warningf(e string, args ...interface{}) {
	if logger.level > LevelWarning {
		return
	}
	logger.provider.do(LevelWarning, fmt.Sprintf(e, args...))
}

func Fatal(e string) {
	logger.provider.do(LevelFatal, e)
}

func Fatalf(e string, args ...interface{}) {
	logger.provider.do(LevelFatal, fmt.Sprintf(e, args...))
}
