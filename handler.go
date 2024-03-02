package errlog

import (
	"context"
	"log/slog"
)

// Handler is a slog.Handler that adds error and stack trace information to log records.
type Handler struct {
	base slog.Handler
	opts HandlerOptions
}

var _ slog.Handler = (*Handler)(nil)

// HandlerOptions contains options for the Handler.
type HandlerOptions struct {
	// SuppressStackTrace suppresses the stack trace from being added to log records.
	SuppressStackTrace bool
	// OverrideSource overrides the source location of the log record with the source location of the error.
	OverrideSource bool
	// StackTraceFormatter is a function that formats the stack trace.
	StackTraceFormatter func(stack []uintptr) string
}

// NewHandler returns a new Handler that wraps the given base slog.handler.
func NewHandler(base slog.Handler, opts *HandlerOptions) *Handler {
	if opts == nil {
		opts = &HandlerOptions{}
	}

	return &Handler{
		base: base,
		opts: *opts,
	}
}

// Enabled is a thin wrapper around the base handler's Enabled method.
func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.base.Enabled(ctx, level)
}

// WithAttrs is a thin wrapper around the base handler's WithAttrs method.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{base: h.base.WithAttrs(attrs)}
}

// WithGroup is a thin wrapper around the base handler's WithGroup method.
func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{base: h.base.WithGroup(name)}
}

// Handle adds error and stack trace information to the log record.
func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	r.Attrs(func(a slog.Attr) bool {
		const stopIter = false

		if a.Key != ErrorKey {
			return true
		}

		err, ok := a.Value.Any().(error)
		if !ok {
			return stopIter
		}

		a.Value = slog.StringValue(err.Error())

		stack := make([]uintptr, 0)
		if str, ok := err.(StackTracer); ok {
			stack = str.Stack()
		}

		if len(stack) == 0 {
			return stopIter
		}

		if h.opts.OverrideSource {
			r.PC = stack[0]
		}

		if !h.opts.SuppressStackTrace {
			if h.opts.StackTraceFormatter != nil {
				r.AddAttrs(slog.String(StackTraceKey, h.opts.StackTraceFormatter(stack)))
				return stopIter
			}
			r.AddAttrs(slog.String(StackTraceKey, formatStack(stack)))
		}
		return stopIter
	})
	return h.base.Handle(ctx, r)
}
