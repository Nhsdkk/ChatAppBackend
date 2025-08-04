package exceptions

import (
	"fmt"
	"runtime"
)

const TraceCallOffset = 3
const TraceBufferSize = 32

type ITrackableException interface {
	error
	GetStackTrace() []TracedFunction
}

type TrackableException struct {
	error
	stackTrace []TracedFunction
}

func (t TrackableException) GetStackTrace() []TracedFunction {
	return t.stackTrace
}

func getTrace() []TracedFunction {
	offset := TraceCallOffset

	callersBuffer := make([]uintptr, TraceBufferSize)
	callers := make([]uintptr, 0)

	for {
		n := runtime.Callers(offset, callersBuffer)
		if n == 0 {
			break
		}

		callers = append(callers, callersBuffer[:n]...)
		offset += n
	}

	callerFrames := runtime.CallersFrames(callers)
	trace := make([]TracedFunction, 0)

	for {
		frame, next := callerFrames.Next()
		if !next {
			break
		}

		trace = append(trace, TracedFunction{
			Function: frame.Function,
			Line:     frame.Line,
			File:     frame.File,
		})
	}

	return trace
}

func CreateTrackableExceptionFromStringF(format string, args ...interface{}) ITrackableException {
	return TrackableException{
		error:      fmt.Errorf(format, args...),
		stackTrace: getTrace(),
	}
}

func WrapErrorWithTrackableException(err error) ITrackableException {
	return TrackableException{
		error:      err,
		stackTrace: getTrace(),
	}
}
