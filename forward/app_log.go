package forward

import "log"

const (
	Debug int = iota
	Info
	Error
)

type AppLog struct {
	*log.Logger
	level int
}

func NewLog() *AppLog {
	return &AppLog{Logger: log.Default()}
}

func (l *AppLog) SetLevel(level int) {
	l.level = level
}

func (l *AppLog) Infof(format string, msg ...any) {
	l.printf(Info, "[Info] "+format, msg...)
}

func (l *AppLog) Debugf(format string, msg ...any) {
	l.printf(Debug, "[Debug] "+format, msg...)
}

func (l *AppLog) Errorf(format string, msg ...any) {
	l.printf(Error, "[Error] "+format, msg...)
}

func (l *AppLog) printf(level int, format string, msg ...any) {
	if l.level <= level {
		l.Logger.Printf(format, msg...)
	}
}
