package log

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	detail = logrus.New()
)

// Level type
type Level string

const (
	traceLevel = "trace"
	debugLevel = "debug"
	infoLevel  = "info"
	warnLevel  = "warn"
	errorLevel = "error"
	fatalLevel = "fatal"
	panicLevel = "panic"
)

func init() {
	// Set log level base on os env
	envLevel := strings.ToLower(os.Getenv("LOG_LEVEL"))
	SetLevel(Level(envLevel))

	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02T15:04:05.999",
	})

	detail.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
		TimestampFormat: "2006-01-02T15:04:05.999",
	})
	detail.SetReportCaller(true)
}

// SetLevel sets the standard logger level.
func SetLevel(level Level) {
	switch level {
	case traceLevel:
		logrus.SetLevel(logrus.TraceLevel)
		detail.SetLevel(logrus.TraceLevel)
	case debugLevel:
		logrus.SetLevel(logrus.DebugLevel)
		detail.SetLevel(logrus.DebugLevel)
	case infoLevel:
		logrus.SetLevel(logrus.InfoLevel)
		detail.SetLevel(logrus.InfoLevel)
	case warnLevel:
		logrus.SetLevel(logrus.WarnLevel)
		detail.SetLevel(logrus.WarnLevel)
	case errorLevel:
		logrus.SetLevel(logrus.ErrorLevel)
		detail.SetLevel(logrus.ErrorLevel)
	case fatalLevel:
		logrus.SetLevel(logrus.FatalLevel)
		detail.SetLevel(logrus.FatalLevel)
	case panicLevel:
		logrus.SetLevel(logrus.PanicLevel)
		detail.SetLevel(logrus.PanicLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
		detail.SetLevel(logrus.InfoLevel)
	}
}

var (
	// Trace logs a message at level Trace on the standard logger.
	Trace = logrus.Trace

	// Debug logs a message at level Debug on the standard logger.
	Debug = logrus.Debug

	// Print logs a message at level Info on the standard logger.
	Print = logrus.Print

	// Info logs a message at level Info on the standard logger.
	Info = logrus.Info

	// Warn logs a message at level Warn on the standard logger.
	Warn = logrus.Warn

	// Warning logs a message at level Warn on the standard logger.
	Warning = logrus.Warning

	// Error logs a message at level Error on the standard logger.
	Error = detail.Error

	// Panic logs a message at level Panic on the standard logger.
	Panic = detail.Panic

	// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
	Fatal = detail.Fatal

	// Tracef logs a message at level Trace on the standard logger.
	Tracef = logrus.Tracef

	// Debugf logs a message at level Debug on the standard logger.
	Debugf = logrus.Debugf

	// Printf logs a message at level Info on the standard logger.
	Printf = logrus.Printf

	// Infof logs a message at level Info on the standard logger.
	Infof = logrus.Infof

	// Warnf logs a message at level Warn on the standard logger.
	Warnf = logrus.Warnf

	// Warningf logs a message at level Warn on the standard logger.
	Warningf = logrus.Warningf

	// Errorf logs a message at level Error on the standard logger.
	Errorf = detail.Errorf

	// Panicf logs a message at level Panic on the standard logger.
	Panicf = detail.Panicf

	// Fatalf logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
	Fatalf = detail.Fatalf

	// Traceln logs a message at level Trace on the standard logger.
	Traceln = logrus.Traceln

	// Debugln logs a message at level Debug on the standard logger.
	Debugln = logrus.Debugln

	// Println logs a message at level Info on the standard logger.
	Println = logrus.Println

	// Infoln logs a message at level Info on the standard logger.
	Infoln = logrus.Infoln

	// Warnln logs a message at level Warn on the standard logger.
	Warnln = logrus.Warnln

	// Warningln logs a message at level Warn on the standard logger.
	Warningln = logrus.Warningln

	// Errorln logs a message at level Error on the standard logger.
	Errorln = detail.Errorln

	// Panicln logs a message at level Panic on the standard logger.
	Panicln = detail.Panicln

	// Fatalln logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
	Fatalln = detail.Fatalln
)
