package errlog

import (
	"log/slog"
)

const (
	ErrorKey      = "error"
	StackTraceKey = "stack_trace"
)

// Err returns an attribute that contains the given error.
// If the error does not implement the StackTracer interface, it will be wrapped with the stack trace.
func Err(err error) slog.Attr {
	if _, ok := err.(StackTracer); !ok {
		err = wrapError(err, 1)
	}
	return slog.Any(ErrorKey, err)
}

// WrapError wraps the given error with a stack trace.
func WrapError(err error) error {
	return wrapError(err, 1)
}
