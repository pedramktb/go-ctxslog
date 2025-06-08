package ctxslog_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/pedramktb/go-ctxslog"
	"github.com/stretchr/testify/assert"
)

func Test_FromContextWithAttrs(t *testing.T) {
	reader, writer := io.Pipe()
	defer writer.Close()
	defer reader.Close()

	output := make([]byte, 0)
	go func() {
		for {
			buf := make([]byte, 1024)
			n, err := reader.Read(buf)
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrClosedPipe) {
				break
			} else if err != nil {
				t.Error(err)
			}
			output = append(output, buf[:n]...)
		}
	}()

	ctxWithLogger := ctxslog.NewContext(context.Background(), slog.NewTextHandler(writer, nil))

	attr := slog.String("ctxslog_test_attr_key", "ctxslog_test_attr_value")

	ctxWithLogger = ctxslog.WithAttrs(ctxWithLogger, attr)

	dynamicCounter := "ctxslog_test_dynamic_attr_value_zero"

	fnAttr := func(context.Context, slog.Record) []slog.Attr {
		return []slog.Attr{
			slog.String("ctxslog_test_dynamic_attr_key", dynamicCounter),
		}
	}

	ctxWithLogger = ctxslog.WithAttrs(ctxWithLogger, fnAttr)

	dynamicCounter = "ctxslog_test_dynamic_attr_value_one"

	logger := ctxslog.FromContext(ctxWithLogger)
	if logger == nil {
		t.Error("logger should not be nil")
	}

	logger.InfoContext(ctxWithLogger, "test")

	time.Sleep(time.Second)

	assert.Contains(t, string(output), "ctxslog_test_attr_key")
	assert.Contains(t, string(output), "ctxslog_test_attr_value")
	assert.Contains(t, string(output), "ctxslog_test_dynamic_attr_key")
	assert.Contains(t, string(output), dynamicCounter)
}
