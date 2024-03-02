package errlog_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"testing"

	"github.com/ichizero/errlog"
)

func newWrappedError(t *testing.T) error {
	t.Helper()
	return errlog.WrapError(errors.New("test error"))
}

func newLogger(t *testing.T, sOpts *slog.HandlerOptions, eOpts *errlog.HandlerOptions) (*slog.Logger, *bytes.Buffer) {
	t.Helper()
	var buf bytes.Buffer

	h := slog.NewJSONHandler(&buf, sOpts)
	return slog.New(errlog.NewHandler(h, eOpts)), &buf
}

func unmarshalLog(t *testing.T, buf *bytes.Buffer) map[string]any {
	t.Helper()

	m := map[string]any{}
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatal(err)
	}
	return m
}

func assertTrue(t *testing.T, value bool) {
	t.Helper()

	if !value {
		t.Error("Should be true")
	}
}
