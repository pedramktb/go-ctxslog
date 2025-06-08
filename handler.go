package ctxslog

import (
	"context"
	"log/slog"
)

type attrHandler struct {
	next slog.Handler
}

func newAttrHandler(next slog.Handler) *attrHandler {
	return &attrHandler{next}
}

func (h *attrHandler) Enabled(ctx context.Context, lvl slog.Level) bool {
	return h.next.Enabled(ctx, lvl)
}

func (h *attrHandler) Handle(ctx context.Context, r slog.Record) error {
	rc := r.Clone()
	rc.AddAttrs(attrs(ctx, r)...)
	return h.next.Handle(ctx, rc)
}

func (h *attrHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &attrHandler{next: h.next.WithAttrs(attrs)}
}

func (h *attrHandler) WithGroup(name string) slog.Handler {
	return &attrHandler{next: h.next.WithGroup(name)}
}
