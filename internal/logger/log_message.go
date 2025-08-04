package logger

import (
	"chat_app_backend/internal/exceptions"
	"fmt"
	"log"
	"strings"
)

type LogModifier int

const (
	None LogModifier = iota
	Fatal
	Panic
)

type ILogMessage interface {
	WithFatal() ILogMessage
	WithPanic() ILogMessage
	Log()
}

type LogMessage struct {
	message    string
	stackTrace *[]exceptions.TracedFunction
	logger     *log.Logger
	modifier   LogModifier
}

func (l LogMessage) WithFatal() ILogMessage {
	if l.modifier != None {
		panic("modifier already set for the message")
	}

	l.modifier = Fatal
	return l
}

func (l LogMessage) WithPanic() ILogMessage {
	if l.modifier != None {
		panic("modifier already set for the message")
	}

	l.modifier = Panic
	return l
}

func formatStackTrace(stackTrace *[]exceptions.TracedFunction) string {
	stackTraceStringSlice := make([]string, len(*stackTrace))
	for idx, function := range *stackTrace {
		stackTraceStringSlice[idx] = fmt.Sprintf("%s:%d (at %s function) ->", function.File, function.Line, function.Function)
	}
	return strings.Join(stackTraceStringSlice, "\n")
}

func (l LogMessage) Log() {
	message := fmt.Sprintf("Message: %s", l.message)

	if l.stackTrace != nil {
		message += fmt.Sprintf("\nStack trace: %s", formatStackTrace(l.stackTrace))
	}

	switch l.modifier {
	case None:
		l.logger.Println(message)
	case Fatal:
		l.logger.Fatalln(message)
	case Panic:
		l.logger.Panicln(message)
	}
}
