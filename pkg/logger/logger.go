package logger

import "fmt"

type ILogger interface {
	Log(message string)
}

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Log(message string) {
	fmt.Println(message)
}
