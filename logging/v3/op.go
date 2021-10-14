package logging

import (
	"os"
	"strings"

	"github.com/op/go-logging"
)

type LoggerLocal struct {
	logger   *logging.Logger
	severity string
}

// New ceates a new 9 Spokes logger with a formatter, uses stdout as a backend.
func NewOplogging(category, level string) *LoggerLocal {
	logger := logging.MustGetLogger(category)

	var format = logging.MustStringFormatter(
		`%{color}%{time:2006/01/02 15:04:05.000} %{level:.4s} ▶ %{message}%{color:reset}`,
	)
	logging.SetBackend(logging.NewBackendFormatter(logging.NewLogBackend(os.Stdout, "", 0), format))

	level = strings.ToUpper(level)
	switch level {

	case "DEBUG":
		logging.SetLevel(logging.DEBUG, "")
	case "ERROR":
		logging.SetLevel(logging.ERROR, "")
	case "FATAL":
		logging.SetLevel(logging.CRITICAL, "")
	case "WARNING":
		logging.SetLevel(logging.WARNING, "")
	default:
		logging.SetLevel(logging.INFO, "")
	}
	return &LoggerLocal{
		logger:   logger,
		severity: level,
	}

}

func (p *LoggerLocal) do(l int, m string) {

	switch l {
	case LevelDebug:
		p.logger.Debug(m)
	case LevelError:
		p.logger.Error(m)
	case LevelFatal:
		p.logger.Fatal(m)
	case LevelWarning:
		p.logger.Warning(m)
	default:
		p.logger.Info(m)
	}
}
