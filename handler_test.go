package errlog_test

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"testing"

	"github.com/ichizero/errlog"
)

func TestHandler_Handle_OverrideSource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		overrideSource bool
		wantFileSuffix string
	}{
		{
			name:           "should override source with stack trace",
			overrideSource: true,
			wantFileSuffix: "errlog/testhelper_test.go",
		},
		{
			name:           "should not override source with stack trace",
			overrideSource: false,
			wantFileSuffix: "errlog/handler_test.go",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			logger, buf := newLogger(t,
				&slog.HandlerOptions{AddSource: true},
				&errlog.HandlerOptions{OverrideSource: tt.overrideSource})

			errWithStack := newWrappedError(t)
			logger.ErrorContext(context.Background(), "test", slog.Any(errlog.ErrorKey, errWithStack))

			got := unmarshalLog(t, buf)

			source, ok := got[slog.SourceKey]
			assertTrue(t, ok)

			sourceMap, ok := source.(map[string]any)
			assertTrue(t, ok)

			file, ok := sourceMap["file"]
			assertTrue(t, ok)

			fileStr, ok := file.(string)
			assertTrue(t, ok)

			if !strings.HasSuffix(fileStr, tt.wantFileSuffix) {
				t.Errorf("file is incorrect: wantFileSuffix=%s, got=%s", tt.wantFileSuffix, file)
			}
		})
	}
}

func TestHandler_Handle_StackTrace(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		err                error
		suppressStackTrace bool
		wantStackTrace     bool
	}{
		{
			name:               "should not output stack trace when err has no stack trace",
			err:                errors.New("no stack trace"),
			suppressStackTrace: false,
			wantStackTrace:     false,
		},
		{
			name:               "should not output stack trace when err has no stack trace and suppress option is true",
			err:                errors.New("no stack trace"),
			suppressStackTrace: true,
			wantStackTrace:     false,
		},
		{
			name:               "should output stack trace when err has stack trace",
			err:                newWrappedError(t),
			suppressStackTrace: false,
			wantStackTrace:     true,
		},
		{
			name:               "should not output stack trace when err has stack trace but suppress option is true",
			err:                newWrappedError(t),
			suppressStackTrace: true,
			wantStackTrace:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			logger, buf := newLogger(t,
				&slog.HandlerOptions{AddSource: true},
				&errlog.HandlerOptions{SuppressStackTrace: tt.suppressStackTrace})

			logger.ErrorContext(context.Background(), "test", slog.Any(errlog.ErrorKey, tt.err))

			got := unmarshalLog(t, buf)

			stackTrace, ok := got[errlog.StackTraceKey]
			if tt.wantStackTrace != ok {
				t.Errorf("stack trace is incorrect: wantStackTrace=%t, got=%s", tt.wantStackTrace, stackTrace)
			}
			if !tt.wantStackTrace {
				return
			}

			stackTraceStr, ok := stackTrace.(string)
			assertTrue(t, ok)

			if !strings.Contains(stackTraceStr, "TestHandler_Handle_StackTrace") {
				t.Errorf("stack trace is incorrect: got=%s", stackTrace)
			}
		})
	}
}
