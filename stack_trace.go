package errlog

import (
	"runtime"
	"strconv"
	"strings"
)

// StackTracer is an interface that represents an error that can provide a stack trace.
type StackTracer interface {
	Stack() []uintptr
}

func formatStack(stack []uintptr) string {
	frames := runtime.CallersFrames(stack)

	var sb strings.Builder

	for {
		frame, more := frames.Next()

		sb.WriteString(frame.Function)
		sb.WriteString("\n\t")
		sb.WriteString(frame.File)
		sb.WriteString(":")
		sb.WriteString(strconv.Itoa(frame.Line))
		sb.WriteString("\n")

		if !more {
			break
		}
	}

	return sb.String()
}

type withStackTraceError struct {
	err   error
	stack []uintptr
}

var (
	_ error       = (*withStackTraceError)(nil)
	_ StackTracer = (*withStackTraceError)(nil)
)

func wrapError(err error, additionalSkip int) error {
	const skipCount = 2
	const depth = 16

	stack := make([]uintptr, depth)
	frameCount := runtime.Callers(skipCount+additionalSkip, stack)

	return &withStackTraceError{
		err:   err,
		stack: stack[:frameCount],
	}
}

func (e withStackTraceError) Error() string {
	return e.err.Error()
}

func (e withStackTraceError) Unwrap() error {
	return e.err
}

func (e withStackTraceError) Stack() []uintptr {
	return e.stack
}
