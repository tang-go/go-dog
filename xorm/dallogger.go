package xorm

import (
	"github.com/tang-go/go-dog/log"
	"xorm.io/core"
)

// Logger is the default implment of core.ILogger
type Logger struct {
	level   core.LogLevel
	showSQL bool
}

//NewLogger 初始化日志
func NewLogger(level core.LogLevel, showSQL bool) *Logger {
	l := new(Logger)
	l.level = level
	l.showSQL = showSQL
	return l
}

// Error implement core.ILogger
func (s *Logger) Error(v ...interface{}) {
	log.Errorln(v...)
	return
}

// Errorf implement core.ILogger
func (s *Logger) Errorf(format string, v ...interface{}) {
	log.Errorf(format, v...)
	return
}

// Debug implement core.ILogger
func (s *Logger) Debug(v ...interface{}) {
	log.Debugln(v...)
	return
}

// Debugf implement core.ILogger
func (s *Logger) Debugf(format string, v ...interface{}) {
	if !s.showSQL {
		return
	}
	if s.level > core.LOG_DEBUG {
		return
	}
	log.Debugf(format, v...)
	return
}

// Info implement core.ILogger
func (s *Logger) Info(v ...interface{}) {
	log.Debugln(v...)
	return
}

// Infof implement core.ILogger
func (s *Logger) Infof(format string, v ...interface{}) {
	log.Infof(format, v...)
	return
}

// Warn implement core.ILogger
func (s *Logger) Warn(v ...interface{}) {
	log.Warnln(v...)
	return
}

// Warnf implement core.ILogger
func (s *Logger) Warnf(format string, v ...interface{}) {
	log.Warnf(format, v...)
	return
}

// Level implement core.ILogger
func (s *Logger) Level() core.LogLevel {
	return s.level
}

// SetLevel implement core.ILogger
func (s *Logger) SetLevel(l core.LogLevel) {
	s.level = l
	return
}

// ShowSQL implement core.ILogger
func (s *Logger) ShowSQL(show ...bool) {
	if len(show) == 0 {
		s.showSQL = true
		return
	}
	s.showSQL = show[0]
}

// IsShowSQL implement core.ILogger
func (s *Logger) IsShowSQL() bool {
	return s.showSQL
}
