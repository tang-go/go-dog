package xorm

import (
	mylog "github.com/tang-go/go-dog/log"
	"xorm.io/xorm/log"
)

// Logger is a logger interface
type Logger struct {
	level   log.LogLevel
	showSQL bool
}

//Debug Debug
func (l *Logger) Debug(v ...interface{}) {
	mylog.Debugln(v...)
}

//Debugf Debugf
func (l *Logger) Debugf(format string, v ...interface{}) {
	mylog.Debugf(format, v...)
}

//Error Error
func (l *Logger) Error(v ...interface{}) {
	mylog.Errorln(v...)
}

//Errorf Errorf
func (l *Logger) Errorf(format string, v ...interface{}) {
	mylog.Errorf(format, v...)
}

//Info Info
func (l *Logger) Info(v ...interface{}) {
	mylog.Infoln(v...)
}

//Infof Infof
func (l *Logger) Infof(format string, v ...interface{}) {
	mylog.Infof(format, v...)
}

//Warn Warn
func (l *Logger) Warn(v ...interface{}) {
	mylog.Warnln(v...)
}

//Warnf Warnf
func (l *Logger) Warnf(format string, v ...interface{}) {
	mylog.Warnf(format, v...)
}

//Level Level
func (l *Logger) Level() log.LogLevel {
	return l.level
}

//SetLevel SetLevel
func (l *Logger) SetLevel(level log.LogLevel) {
	l.level = level
}

//ShowSQL ShowSQL
func (l *Logger) ShowSQL(show ...bool) {
	l.showSQL = true
	for _, s := range show {
		if s == false {
			l.showSQL = false
		}
	}
}

//IsShowSQL IsShowSQL
func (l *Logger) IsShowSQL() bool {
	return l.showSQL
}
