package logger

import (
	"github.com/sirupsen/logrus"
	"json-rpc-node-proxy/pkg/env"
	"os"
)

type ILogger interface {
	Log(message string)
	Warn(message string)
	Error(message string)
	Trace(message string)
	Debug(message string)
	Fatal(message string)
	Panic(message string)
}

type Logger struct {
	internalLogger *logrus.Logger
}

func NewLogger(environment env.Environment) *Logger {
	logger := logrus.New()

	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	if environment == env.EnvTest {
		logger.SetLevel(logrus.TraceLevel)
	} else {
		logger.SetLevel(logrus.DebugLevel)
	}

	return &Logger{
		internalLogger: logger,
	}
}

func (l *Logger) Log(message string) {
	l.internalLogger.Info(message)
}

func (l *Logger) Warn(message string) {
	l.internalLogger.Warn(message)
}

func (l *Logger) Error(message string) {
	l.internalLogger.Error(message)
}

func (l *Logger) Trace(message string) {
	l.internalLogger.Trace(message)
}

func (l *Logger) Debug(message string) {
	l.internalLogger.Debug(message)
}

func (l *Logger) Fatal(message string) {
	l.internalLogger.Fatal(message)
}

func (l *Logger) Panic(message string) {
	l.internalLogger.Panic(message)
}
