# errlog

[![Test](https://github.com/ichizero/errlog/actions/workflows/test.yml/badge.svg)](https://github.com/ichizero/errlog/actions/workflows/test.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/ichizero/errlog.svg)](https://pkg.go.dev/github.com/ichizero/errlog)
[![Codecov](https://codecov.io/gh/ichizero/errlog/branch/main/graph/badge.svg)](https://codecov.io/gh/ichizero/errlog)
[![Go Report Card](https://goreportcard.com/badge/github.com/ichizero/errlog)](https://goreportcard.com/report/github.com/ichizero/errlog)

`errlog` is a error logging package based on [log/slog](https://pkg.go.dev/log/slog) standard library.
It provides error logging with stack trace and source location.
It does not require any third-party package. 

## 🚀 Installation

```bash
go get github.com/ichizero/errlog
```

## 🧐 Usage

### Initialize logger
`errlog.NewHandler` wraps `slog.Handler`, so you can provide `*slog.JSONHandler`, `*slog.TextHandler`,
or any other handler.

```go
h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})
hErr := errlog.NewHandler(h, &errlog.HandlerOptions{OverrideSource: true, SuppressStackTrace: false})
slog.SetDefault(slog.New(hErr))
```

### Logging error with stack trace

#### With errlog.Err
`errlog.Err` wraps error with stack trace and returns `slog.Attr` with key `error`.

```go
err := errors.New("test error")
slog.ErrorContext(ctx, "test", errlog.Err(err))
```

#### With custom error

`errlog.NewHandler` outputs stack trace with the error that implements `errlog.StackTrace` interface,
so you can provide custom error with stack trace.

```go
type yourCustomError struct {
	err error
	stack []uintptr
}

func (e yourCustomError) Stack() []uintptr {
	return e.stack
}
```

If so, you can log stack trace without using `errlog.Err`.

```go
err := newYourCustomError("error")
slog.ErrorContext(ctx, "test", slog.Any("error", err))
```

#### Example usage

```go
package main

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/ichizero/errlog"
)

func main() {
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})
	hErr := errlog.NewHandler(h, &errlog.HandlerOptions{OverrideSource: true, SuppressStackTrace: false})
	slog.SetDefault(slog.New(hErr))

	ctx := context.Background()

	err := errors.New("test error")
	slog.ErrorContext(ctx, "test", errlog.Err(err))

	err = errlog.WrapError(err)
	slog.ErrorContext(ctx, "test", slog.Any("error", err))
}
```
