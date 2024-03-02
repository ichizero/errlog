package errlog_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"
	"testing/slogtest"

	"github.com/ichizero/errlog"
)

func TestRun(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	newHandler := func(*testing.T) slog.Handler {
		buf.Reset()

		h := slog.NewJSONHandler(&buf, &slog.HandlerOptions{AddSource: true})
		return errlog.NewHandler(h, &errlog.HandlerOptions{OverrideSource: true, SuppressStackTrace: false})
	}

	result := func(t *testing.T) map[string]any {
		t.Helper()

		m := map[string]any{}
		if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
			t.Fatal(err)
		}
		return m
	}

	slogtest.Run(t, newHandler, result)
}
