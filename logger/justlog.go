package log

import (
	"sync"
	"time"

	"github.com/rifflock/lfshook"
	logrus "github.com/sirupsen/logrus"

	nested "github.com/antonfisher/nested-logrus-formatter"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

// ILogger is
type ILogger interface {
	Debug(title string, description string, args ...interface{})
	Info(title string, description string, args ...interface{})
	Warn(title string, description string, args ...interface{})
	Error(title string, description string, args ...interface{})
	Fatal(title string, description string, args ...interface{})
	Panic(title string, description string, args ...interface{})
	WithFile(filename string, maxAge int)
}

var defaultLogger logrusImpl
var defaultLoggerOnce sync.Once

// LogrusImpl is
type logrusImpl struct {
	theLogger *logrus.Logger
	useFile   bool
}

// GetLog is
func GetLog() ILogger {
	defaultLoggerOnce.Do(func() {
		defaultLogger = logrusImpl{theLogger: logrus.New()}
		defaultLogger.useFile = false
		defaultLogger.theLogger.SetFormatter(&nested.Formatter{
			NoColors: true,
			HideKeys: true,
		})
	})
	return &defaultLogger
}

// WithFile is command to state the log will printing to files
// the rolling log file will put in logs/ directory
//
// filename is just a name of log file without any extension
//
// maxAge is age (in days) of the logs file before it gets purged from the file system
func (l *logrusImpl) WithFile(filename string, maxAge int) {

	if !l.useFile {

		if maxAge <= 0 {
			panic("maxAge should > 0")
		}

		path := filename + ".log"
		writer, _ := rotatelogs.New(
			"/var/log/"+path+".%Y%m%d",
			rotatelogs.WithLinkName(path),
			rotatelogs.WithMaxAge(time.Duration(maxAge*24)*time.Hour),
			rotatelogs.WithRotationTime(time.Duration(1*24)*time.Hour),
		)

		defaultLogger.theLogger.AddHook(lfshook.NewHook(
			lfshook.WriterMap{
				logrus.InfoLevel:  writer,
				logrus.WarnLevel:  writer,
				logrus.ErrorLevel: writer,
				logrus.DebugLevel: writer,
			},
			defaultLogger.theLogger.Formatter,
		))

		l.useFile = true
	}
}

// Debug is
func (l *logrusImpl) Debug(title string, description string, args ...interface{}) {
	l.theLogger.WithField("key", title).Debugf(description, args...)
}

// Info is
func (l *logrusImpl) Info(title string, description string, args ...interface{}) {
	l.theLogger.WithField("key", title).Infof(description, args...)
}

// Warn is
func (l *logrusImpl) Warn(title string, description string, args ...interface{}) {
	l.theLogger.WithField("key", title).Warnf(description, args...)
}

// Error is
func (l *logrusImpl) Error(title string, description string, args ...interface{}) {
	l.theLogger.WithField("key", title).Errorf(description, args...)
}

// Fatal is
func (l *logrusImpl) Fatal(title string, description string, args ...interface{}) {
	l.theLogger.WithField("key", title).Fatalf(description, args...)
}

// Panic is
func (l *logrusImpl) Panic(title string, description string, args ...interface{}) {
	l.theLogger.WithField("key", title).Panicf(description, args...)
}
