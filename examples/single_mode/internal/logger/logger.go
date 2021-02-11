package logger

import "fmt"

type Logger struct {
	permPrefix string
	tempPrefix string
}

func (l *Logger) SetMainPrefix(prefix string) {
	l.permPrefix = prefix
}

func (l *Logger) SetTempPrefix(prefix string) {
	l.tempPrefix = prefix
}

func (l *Logger) Log(msg string, level string) {
	message := "["+l.permPrefix+"]["+l.tempPrefix+"] :: "+msg
	fmt.Println(message)
}