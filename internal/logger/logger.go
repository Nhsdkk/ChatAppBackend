package logger

import (
	"fmt"
	"io"
	"log"
)

type ILogger interface {
	CreateInfoMessage(msg string) ILogMessage
	CreateInfoMessageF(format string, args ...interface{}) ILogMessage
	CreateErrorMessage(err error) ILogMessage
	CreateWarningMessage(msg string) ILogMessage
	CreateWarningMessageF(format string, args ...interface{}) ILogMessage
	CreateDebugMessage(msg string) ILogMessage
	CreateDebugMessageF(format string, args ...interface{}) ILogMessage
}

type Logger struct {
	debug   *log.Logger
	info    *log.Logger
	warning *log.Logger
	error   *log.Logger
}

func (l Logger) CreateInfoMessage(msg string) ILogMessage {
	return &LogMessage{
		message:  msg,
		logger:   l.info,
		modifier: None,
	}
}

func (l Logger) CreateInfoMessageF(format string, args ...interface{}) ILogMessage {
	return &LogMessage{
		message:  fmt.Sprintf(format, args...),
		logger:   l.info,
		modifier: None,
	}
}

func (l Logger) CreateErrorMessage(err error) ILogMessage {
	return &LogMessage{
		message:  err.Error(),
		logger:   l.error,
		modifier: None,
	}
}

func (l Logger) CreateWarningMessage(msg string) ILogMessage {
	return &LogMessage{
		message:  msg,
		logger:   l.warning,
		modifier: None,
	}
}

func (l Logger) CreateWarningMessageF(format string, args ...interface{}) ILogMessage {
	return &LogMessage{
		message:  fmt.Sprintf(format, args...),
		logger:   l.warning,
		modifier: None,
	}
}

func (l Logger) CreateDebugMessage(msg string) ILogMessage {
	return &LogMessage{
		message:  msg,
		logger:   l.debug,
		modifier: None,
	}
}

func (l Logger) CreateDebugMessageF(format string, args ...interface{}) ILogMessage {
	return &LogMessage{
		message:  fmt.Sprintf(format, args...),
		logger:   l.debug,
		modifier: None,
	}
}

func CreateLogger(writer io.Writer) ILogger {
	logger := &Logger{}
	flags := log.Ldate | log.Ltime | log.Llongfile | log.LUTC
	logger.debug = log.New(writer, "DEBUG: ", flags)
	logger.info = log.New(writer, "INFO: ", flags)
	logger.warning = log.New(writer, "WARNING: ", flags)
	logger.error = log.New(writer, "ERROR: ", flags)
	return logger
}
