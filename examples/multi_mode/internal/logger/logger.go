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

func (l *Logger) LogLimited(msg string) {
	message := "["+l.permPrefix+"]["+l.tempPrefix+"] `LIMITED` :: "+msg
	fmt.Println(message)
}

func (l *Logger) LogVerbose(msg string) {
	message := "["+l.permPrefix+"]["+l.tempPrefix+"] `VERBOSE` :: "+msg
	fmt.Println(message)
}