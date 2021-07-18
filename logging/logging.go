package logging

import (
	"fmt"
	"os"
	"strings"

	"github.com/op/go-logging"
)

// New ceates a new 9 Spokes logger with a formatter, uses stdout as a backend.
func New(category, level string) *logging.Logger {
	logger := logging.MustGetLogger(category)

	var format = logging.MustStringFormatter(
		`%{color}%{time:2006/01/02 15:04:05.000} %{level:.4s} %{shortfunc}() ▶ %{message}%{color:reset}`,
	)
	logging.SetBackend(logging.NewBackendFormatter(logging.NewLogBackend(os.Stdout, "", 0), format))

	SetLogLevel(level)
	return logger

}

// Update 9 Spokes logger level.
func SetLogLevel(level string) error {

	level = strings.ToUpper(level)
	switch level {
	case "INFO":
		logging.SetLevel(logging.INFO, "")
	case "DEBUG":
		logging.SetLevel(logging.DEBUG, "")
	case "ERROR":
		logging.SetLevel(logging.ERROR, "")
	case "CRITICAL":
		logging.SetLevel(logging.CRITICAL, "")
	default:
		return fmt.Errorf("unsupported log level %s", level)
	}
	return nil
}
