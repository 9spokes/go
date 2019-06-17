package logging

import (
	"os"
	"strings"

	"github.com/op/go-logging"
)

// Logger is a logging object wrapper
type Logger *logging.Logger

// New ceates a new 9 Spokes logger with a formatter, uses stdout as a backend.
func New(category, level string) Logger {
	logger := logging.MustGetLogger(category)

	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{level:.5s} %{shortfunc}() ▶ %{message}%{color:reset}`,
	)
	logging.SetBackend(logging.NewBackendFormatter(logging.NewLogBackend(os.Stdout, "", 0), format))

	level = strings.ToUpper(level)
	if level == "DEBUG" {
		logging.SetLevel(logging.DEBUG, category)
		return logger
	}
	if level == "ERROR" {
		logging.SetLevel(logging.ERROR, category)
		return logger
	}
	if level == "CRITICAL" {
		logging.SetLevel(logging.CRITICAL, category)
		return logger
	}
	logging.SetLevel(logging.INFO, category)

	return logger
}
