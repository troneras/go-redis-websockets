package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/troneras/gorews/logger/config"
)

var conf = config.Configure()
var log = logrus.New()

// Fields type, to be used by the caller
type Fields map[string]interface{}

func Configure() {
	conf = config.Configure()
	log.Out = conf.LogFile
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(conf.LogLevel)
}

// Info logs a message at level Info on the standard logger.
// It optionally accepts fields.
func Info(message string, fields ...Fields) {
	if len(fields) > 0 {
		log.WithFields(logrus.Fields(fields[0])).Info(message)
	} else {
		log.Info(message)
	}
}

// Warn logs a message at level Warn on the standard logger.
// It optionally accepts fields.
func Warn(message string, fields ...Fields) {
	if len(fields) > 0 {
		log.WithFields(logrus.Fields(fields[0])).Warn(message)
	} else {
		log.Warn(message)
	}
}

// Error logs a message at level Error on the standard logger.
// It optionally accepts fields.
func Error(message string, fields ...Fields) {
	if len(fields) > 0 {
		log.WithFields(logrus.Fields(fields[0])).Error(message)
	} else {
		log.Error(message)
	}
}

// Fatal logs a message at level Fatal on the standard logger.
// It optionally accepts fields.
func Fatal(message string, fields ...Fields) {
	if len(fields) > 0 {
		log.WithFields(logrus.Fields(fields[0])).Fatal(message)
	} else {
		log.Fatal(message)
	}
}

// Panic logs a message at level Panic on the standard logger.
// It optionally accepts fields.
func Panic(message string, fields ...Fields) {
	if len(fields) > 0 {
		log.WithFields(logrus.Fields(fields[0])).Panic(message)
	} else {
		log.Panic(message)
	}
}

// Debug logs a message at level Debug on the standard logger.
// It optionally accepts fields.
func Debug(message string, fields ...Fields) {
	if len(fields) > 0 {
		log.WithFields(logrus.Fields(fields[0])).Debug(message)
	} else {
		log.Debug(message)
	}
}

// Avoid using these functions directly, use the ones above instead
func Printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func Println(args ...interface{}) {
	log.Println(args...)
}
