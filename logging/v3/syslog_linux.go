//go:build linux

package logging

import (
	"fmt"
	"log"
	"log/syslog"
	"os"
)

type LoggerSyslog struct {
	writer *syslog.Writer
}

func NewSyslog(host, tag string) *LoggerSyslog {
	w, err := syslog.Dial("udp", host, syslog.LOG_EMERG|syslog.LOG_KERN, tag)
	if err != nil {
		log.Fatal("failed to dial syslog")
	}

	return &LoggerSyslog{
		writer: w,
	}

}
func (p *LoggerSyslog) do(l int, m string) {
	logContent := fmtLogContent(1, m)
	switch l {
	case LevelDebug:
		p.writer.Debug(fmt.Sprintf(ColorDebug, logContent))
	case LevelError:
		p.writer.Err(fmt.Sprintf(ColorError, logContent))
	case LevelWarning:
		p.writer.Warning(fmt.Sprintf(ColorWarning, logContent))
	case LevelFatal:
		p.writer.Emerg(fmt.Sprintf(ColorFatal, logContent))
		os.Exit(1)
	default:
		p.writer.Info(fmt.Sprintf(ColorInfo, logContent))
	}
}
