package log

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/cloudnativedevelopment/cnd/pkg/config"
	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
	"k8s.io/klog"
)

const (
	red    = "\033[1;31m%s\033[0m\n"
	yellow = "\033[1;33m%s\033[0m\n"
	green  = "\033[1;32m%s\033[0m\n"
)

type logger struct {
	out  *logrus.Logger
	file *logrus.Entry
}

var log = &logger{
	out: logrus.New(),
}

// Init configures the logger for the package to use.
func Init(level logrus.Level, actionID string) {
	log.out.SetOutput(os.Stdout)
	log.out.SetLevel(level)

	fileLogger := logrus.New()
	fileLogger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	logPath := path.Join(config.GetCNDHome(), fmt.Sprintf("%s%s", config.GetBinaryName(), ".log"))
	rolling := getRollingLog(logPath)
	fileLogger.SetOutput(rolling)
	fileLogger.SetLevel(logrus.DebugLevel)
	log.file = fileLogger.WithFields(logrus.Fields{"action": actionID})

	klog.InitFlags(nil)
	flag.Set("logtostderr", "false")
	flag.Set("alsologtostderr", "false")
	flag.Parse()
	klog.SetOutput(log.file.Writer())
}

func getRollingLog(path string) io.Writer {
	return &lumberjack.Logger{
		Filename:   path,
		MaxSize:    1, // megabytes
		MaxBackups: 10,
		MaxAge:     28, //days
		Compress:   true,
	}
}

// SetLevel sets the level of the main logger
func SetLevel(level string) {
	l, err := logrus.ParseLevel(level)

	if err == nil {
		log.out.SetLevel(l)
	}
}

// Debug writes a debug-level log
func Debug(args ...interface{}) {
	log.out.Debug(args...)
	if log.file != nil {
		log.file.Debug(args...)
	}

}

// Debugf writes a debug-level log with a format
func Debugf(format string, args ...interface{}) {
	log.out.Debugf(format, args...)
	if log.file != nil {
		log.file.Debugf(format, args...)
	}
}

// Info writes a info-level log
func Info(args ...interface{}) {
	log.out.Info(args...)
	if log.file != nil {
		log.file.Info(args...)
	}
}

// Infof writes a info-level log with a format
func Infof(format string, args ...interface{}) {
	log.out.Infof(format, args...)
	if log.file != nil {
		log.file.Infof(format, args...)
	}

}

// Error writes a error-level log
func Error(args ...interface{}) {
	log.out.Error(args...)
	if log.file != nil {
		log.file.Error(args...)
	}
}

// Errorf writes a error-level log with a format
func Errorf(format string, args ...interface{}) {
	log.out.Errorf(format, args...)
	if log.file != nil {
		log.file.Errorf(format, args...)
	}
}

// Red writes a line in red
func Red(format string, args ...interface{}) {
	fmt.Printf(red, fmt.Sprintf(format, args...))
	log.file.Errorf(format, args...)
}

// Yellow writes a line in yellow
func Yellow(format string, args ...interface{}) {
	fmt.Printf(yellow, fmt.Sprintf(format, args...))
	log.file.Warnf(format, args...)
}

// Green writes a line in green
func Green(format string, args ...interface{}) {
	fmt.Printf(green, fmt.Sprintf(format, args...))
	log.file.Infof(format, args...)
}
