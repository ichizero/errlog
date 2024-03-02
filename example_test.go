package errlog_test

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/ichizero/errlog"
)

func Example() {
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})
	hErr := errlog.NewHandler(h, &errlog.HandlerOptions{OverrideSource: true, SuppressStackTrace: false})
	slog.SetDefault(slog.New(hErr))

	ctx := context.Background()

	err := errors.New("test error")
	slog.ErrorContext(ctx, "test", errlog.Err(err))

	err = errlog.WrapError(err)
	slog.ErrorContext(ctx, "test", slog.Any("error", err))
}
