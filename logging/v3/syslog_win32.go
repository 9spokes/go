//go:build windows

package logging

func NewSyslog(host, tag string) *LoggerLocal {
	return NewOplogging(host, tag)

}
