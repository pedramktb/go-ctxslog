package ctxslog

import (
	"context"
	"log/slog"
)

func New(handler slog.Handler) *slog.Logger {
	return slog.New(newAttrHandler(handler))
}

type loggerCtxKey struct{}

type attrsCtxKey struct{}

func NewContext(ctx context.Context, handler slog.Handler) context.Context {
	return context.WithValue(ctx, loggerCtxKey{}, slog.New(newAttrHandler(handler)))
}

func FromContext(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(loggerCtxKey{}).(*slog.Logger)
	if !ok {
		return slog.Default()
	}
	return logger
}

func WithAttrs[T slog.Attr | func(context.Context, slog.Record) []slog.Attr](ctx context.Context, attrs ...T) context.Context {
	attrsCtx, ok := ctx.Value(attrsCtxKey{}).([]any)
	if !ok {
		attrsCtx = make([]any, 0)
	}
	for i := range attrs {
		attrsCtx = append(attrsCtx, attrs[i])
	}
	return context.WithValue(ctx, attrsCtxKey{}, attrsCtx)
}

func attrs(ctx context.Context, r slog.Record) []slog.Attr {
	attrsCtx, ok := ctx.Value(attrsCtxKey{}).([]any)
	if !ok {
		return nil // unreachable
	}
	var attrs []slog.Attr
	for i := range attrsCtx {
		if attr, ok := attrsCtx[i].(slog.Attr); ok {
			attrs = append(attrs, attr)
			continue
		}
		if fn, ok := attrsCtx[i].(func(context.Context, slog.Record) []slog.Attr); ok {
			attrs = append(attrs, fn(ctx, r)...)
		}
		// unreachable as WithAttrs uses generic types
	}
	return attrs
}
