package log

import (
	"fmt"
	"io"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// Level describes the log severity level.
type Level uint8

const (
	// PanicLevel level
	PanicLevel Level = iota
	// FatalLevel level
	FatalLevel
	// ErrorLevel level
	ErrorLevel
	// WarnLevel level
	WarnLevel
	// InfoLevel level
	InfoLevel
	// DebugLevel level
	DebugLevel
	// TraceLevel level
	TraceLevel
)

// Logger is an interface that describes logging.
type Logger interface {
	SetLevel(level Level)
	SetOut(out io.Writer)

	Traceln(...interface{})
	Tracef(string, ...interface{})

	Debugln(...interface{})
	Debugf(string, ...interface{})

	Infoln(...interface{})
	Infof(string, ...interface{})

	Warnln(...interface{})
	Warnf(string, ...interface{})

	Errorln(...interface{})
	Errorf(string, ...interface{})

	Fatalln(...interface{})
	Fatalf(string, ...interface{})

	Panicln(v ...interface{})
	Panicf(format string, v ...interface{})
}

type logger struct {
	entry *logrus.Entry
}

// With attaches a key-value pair to a logger.
func (l logger) With(key string, value interface{}) Logger {
	return logger{l.entry.WithField(key, value)}
}

// WithError attaches an error to a logger.
func (l logger) WithError(err error) Logger {
	return logger{l.entry.WithError(err)}
}

// SetLevel sets the level of a logger.
func (l logger) SetLevel(level Level) {
	l.entry.Logger.Level = logrus.Level(level)
}

// SetOut sets the output destination for a logger.
func (l logger) SetOut(out io.Writer) {
	l.entry.Logger.Out = out
}

// Trace logs a message at level Debug on the standard logger.
func (l logger) Traceln(args ...interface{}) {
	l.sourced().Traceln(args...)
}

// Tracef logs a message at level Debug on the standard logger.
func (l logger) Tracef(format string, args ...interface{}) {
	l.sourced().Tracef(format, args...)
}

// Trace logs a message at level Debug on the standard logger.
func (l logger) Debugln(args ...interface{}) {
	l.sourced().Debugln(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func (l logger) Debugf(format string, args ...interface{}) {
	l.sourced().Debugf(format, args...)
}

// Info logs a message at level Info on the standard logger.
func (l logger) Infoln(args ...interface{}) {
	l.sourced().Infoln(args...)
}

// Infof logs a message at level Info on the standard logger.
func (l logger) Infof(format string, args ...interface{}) {
	l.sourced().Infof(format, args...)
}

// Warn logs a message at level Warn on the standard logger.
func (l logger) Warnln(args ...interface{}) {
	l.sourced().Warnln(args...)
}

// Warnf logs a message at level Warn on the standard logger.
func (l logger) Warnf(format string, args ...interface{}) {
	l.sourced().Warnf(format, args...)
}

// Error logs a message at level Error on the standard logger.
func (l logger) Errorln(args ...interface{}) {
	l.sourced().Errorln(args...)
}

// Errorf logs a message at level Error on the standard logger.
func (l logger) Errorf(format string, args ...interface{}) {
	l.sourced().Errorf(format, args...)
}

// Fatal logs a message at level Fatal on the standard logger.
func (l logger) Fatalln(args ...interface{}) {
	l.sourced().Fatalln(args...)
}

// Fatalf logs a message at level Fatal on the standard logger.
func (l logger) Fatalf(format string, args ...interface{}) {
	l.sourced().Fatalf(format, args...)
}

// Panicln logs a message at level Fatal on the standard logger.
func (l logger) Panicln(args ...interface{}) {
	l.sourced().Panicln(args...)
}

// Panic logs a message at level Fatal on the standard logger.
func (l logger) Panicf(format string, args ...interface{}) {
	l.sourced().Panicf(format, args...)
}

// sourced adds a source field to the logger that contains
// the file name and line where the logging happened.
func (l logger) sourced() *logrus.Entry {
	pc, file, line, ok := runtime.Caller(2)
	fn := "(unknown)"
	if !ok {
		file = "<???>"
		line = 1
	} else {
		function := runtime.FuncForPC(pc).Name()
		index := strings.Index(function, ".")
		slash := strings.LastIndex(file, "/")
		file = function[:index] + "/" + file[slash+1:]
		fn = function[index+1:]
	}
	return l.entry.WithField("PATH", fmt.Sprintf("%s() %s:%d", fn, file, line))
}

var origLogger = logrus.New()
var baseLogger = logger{entry: logrus.NewEntry(origLogger)}

// init
func init() {
	baseLogger.entry.Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000000",
		ForceColors:     true,
	})
	SetLevel(TraceLevel)
}

// New returns a new logger.
func New() Logger {
	return logger{entry: logrus.NewEntry(origLogger)}
}

// Base returns the base logger.
func Base() Logger {
	return baseLogger
}

// SetLevel sets the Level of the base logger
func SetLevel(level Level) {
	baseLogger.entry.Logger.Level = logrus.Level(level)
}

// SetOut sets the output destination base logger
func SetOut(out io.Writer) {
	baseLogger.entry.Logger.Out = out
}

// Withf attaches a key,value pair to a logger.
func Withf(key string, value interface{}) Logger {
	return baseLogger.With(key, value)
}

// WithError returns a Logger that will print an error along with the next message.
func WithError(err error) Logger {
	return logger{entry: baseLogger.sourced().WithError(err)}
}

// Traceln logs a message at level Debug on the standard logger.
func Traceln(args ...interface{}) {
	baseLogger.sourced().Traceln(args...)
}

// Tracef logs a message at level Debug on the standard logger.
func Tracef(format string, args ...interface{}) {
	baseLogger.sourced().Tracef(format, args...)
}

// Debugln logs a message at level Debug on the standard logger.
func Debugln(args ...interface{}) {
	baseLogger.sourced().Debugln(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	baseLogger.sourced().Debugf(format, args...)
}

// Infoln logs a message at level Info on the standard logger.
func Infoln(args ...interface{}) {
	baseLogger.sourced().Infoln(args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	baseLogger.sourced().Infof(format, args...)
}

// Warnln logs a message at level Warn on the standard logger.
func Warnln(args ...interface{}) {
	baseLogger.sourced().Warnln(args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	baseLogger.sourced().Warnf(format, args...)
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(args ...interface{}) {
	baseLogger.sourced().Errorln(args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	baseLogger.sourced().Errorf(format, args...)
}

// Fatalln logs a message at level Fatal on the standard logger.
func Fatalln(args ...interface{}) {
	baseLogger.sourced().Fatalln(args...)
}

// Fatalf logs a message at level Fatal on the standard logger.
func Fatalf(format string, args ...interface{}) {
	baseLogger.sourced().Fatalf(format, args...)
}

// Panicln logs a message at level Fatal on the standard logger.
func Panicln(args ...interface{}) {
	baseLogger.sourced().Panicln(args...)
}

// Panicf logs a message at level Fatal on the standard logger.
func Panicf(format string, args ...interface{}) {
	baseLogger.sourced().Panicf(format, args...)
}
