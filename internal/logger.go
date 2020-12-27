package internal

import (
	"io"
	"time"
)

type LogLevel int

const (
	logLevelInfo  LogLevel = 0
	logLevelWarn  LogLevel = 1
	logLevelError LogLevel = 2
	logLevelOff   LogLevel = 3
)

type Logger struct {
	infoLog  io.Writer
	errorLog io.Writer
	level    LogLevel
}

func (log Logger) Error(message string) {
	if log.level < 3 {
		_, _ = log.errorLog.Write([]byte("[" + timeNow() + "] ERROR: " + message + "\n"))
	}
}

func (log Logger) Warning(message string) {
	if log.level < 2 {
		_, _ = log.infoLog.Write([]byte("[" + timeNow() + "] WARN: " + message + "\n"))
	}
}

func (log Logger) Info(message string) {
	if log.level < 1 {
		_, _ = log.infoLog.Write([]byte("[" + timeNow() + "] INFO: " + message + "\n"))
	}
}

func (log Logger) LogError(err error) {
	log.Error(err.Error())
}

func timeNow() string {
	return time.Now().Format(time.RFC3339Nano)
}
