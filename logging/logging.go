package logging

import (
	"os"
	"strings"

	"github.com/op/go-logging"
)

// New ceates a new 9 Spokes logger with a formatter, uses stdout as a backend.
func New(category, level string) *logging.Logger {
	logger := logging.MustGetLogger(category)

	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{level:.4s} %{shortfunc}() ▶ %{message}%{color:reset}`,
	)
	logging.SetBackend(logging.NewBackendFormatter(logging.NewLogBackend(os.Stdout, "", 0), format))

	level = strings.ToUpper(level)
	switch level {
	case "DEBUG":
		logging.SetLevel(logging.DEBUG, category)
		return logger
	case "ERROR":
		logging.SetLevel(logging.ERROR, category)
		return logger
	case "CRITICAL":
		logging.SetLevel(logging.CRITICAL, category)
		return logger
	default:
		logging.SetLevel(logging.INFO, category)
		return logger
	}

}
