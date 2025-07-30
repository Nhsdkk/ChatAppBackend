package logger

import "log"

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
	message  string
	logger   *log.Logger
	modifier LogModifier
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

func (l LogMessage) Log() {
	switch l.modifier {
	case None:
		l.logger.Println(l.message)
	case Fatal:
		l.logger.Fatalln(l.message)
	case Panic:
		l.logger.Panicln(l.message)
	}
}
