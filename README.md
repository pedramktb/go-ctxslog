# ctxslog

ctxslog is a package that provides slog logger that can be used with context.Context for storing attributes. The logger itself can also be stored and retrieved from context.Context.

## Behavior

ctxslog uses a `slog.Handler` middleware that reads attributes from context and adds them to the log record.

### Constructors

`ctxslog.NewContext` creates a new context with logger stored in it. The logger can be retrieved from context using `ctxslog.NewContext` function. The logger can also be created using `ctxslog.New` function which directly returns the logger without manipulating context. The logger can also be created using `ctxslog.FromContext` function which directly returns the logger without a need for context.

### Attributes

Attributes can be added to the context using `ctxslog.WithAttrs` function. After this point when using context aware slog functions like `logger.InfoContext` or `logger.ErrorContext`, attributes will be added similar to how they can be added directly as parameters. `ctxslog.WithAttrs` also supports adding attributes using a function that returns an array of `slog.Attr`. This is useful when you want to add dynamic attributes based on the context or the `slog.Record`. It can also be used to reference values that are stored in context using custom context keys.

## Example
```go
package main

import (
    "github.com/pedramktb/go-ctxslog"
    "context"
    "log/slog"
    "os"
    "time"

    // other imports...
)

func main() {
    ctx := ctxslog.NewContext(context.Background(), slog.NewTextHandler(os.Stdout, nil))

    // ...

    attrs := []slog.Attr{
        slog.String("db", db.name),
        slog.String("ip", db.ip),
    }

    ctx = ctxslog.WithAttrs(ctx, attrs...)

    ctx = ctxslog.WithAttrs(ctx, func(ctx context.Context, r slog.Record) []slog.Attr {
        if r.Level == slog.LevelDebug {
            return []slog.Attr{
                slog.String("time_took_to_compute_this_log_attr", time.Now().Sub(r.Time).String()),
            }
        } else {
            return []slog.Attr{
                slog.String(ctx.Value("stage").(string)),
            }
        }
    })

    // ...

    logger := ctxslog.FromContext(ctx)

    // this will log something like: db=products ip=127.0.0.1 time_took_to_compute_this_log_attr=14ns
    logger.DebugContext(ctx, "database started")

    // this will log something like: db=products ip=127.0.0.1 stage=dev
    logger.InfoContext(ctx, "database started")

    // ...
}
```
