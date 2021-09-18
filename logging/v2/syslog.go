package logging

import (
	"fmt"
	"log"
	"log/syslog"
	"os"
	"runtime"
	"strings"
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
func (p *LoggerSyslog) log(l int, m string) {

	pc, _, _, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
	lastSlash := strings.LastIndexByte(funcName, '/')
	if lastSlash < 0 {
		lastSlash = 0
	}
	lastDot := strings.LastIndexByte(funcName[lastSlash:], '.') + lastSlash

	m = fmt.Sprintf("[%s.%s()] ▶ %s", funcName[:lastDot], funcName[lastDot+1:], m)

	switch l {
	case LevelDebug:
		p.writer.Debug(fmt.Sprintf(ColorDebug, m))
	case LevelError:
		p.writer.Err(fmt.Sprintf(ColorError, m))
	case LevelWarning:
		p.writer.Warning(fmt.Sprintf(ColorWarning, m))
	case LevelFatal:
		p.writer.Emerg(fmt.Sprintf(ColorFatal, m))
		os.Exit(1)
	default:
		p.writer.Info(fmt.Sprintf(ColorInfo, m))
	}
}
