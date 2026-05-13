package orm

import (
	"context"
	"io"
	"os"
)

var (
	DebugMode   bool      = false
	DebugWriter io.Writer = os.Stdout
)

type contextKeyType int

const (
	ctxKeySkipHooks contextKeyType = iota
	ctxKeySkipTimestamps
	ctxKeyDebug
	ctxKeyDebugWriter
)

func WithDebug(ctx context.Context, debug bool) context.Context {
	return context.WithValue(ctx, ctxKeyDebug, debug)
}

func IsDebug(ctx context.Context) bool {
	if v, ok := ctx.Value(ctxKeyDebug).(bool); ok {
		return v
	}
	return DebugMode
}

func WithDebugWriter(ctx context.Context, w io.Writer) context.Context {
	return context.WithValue(ctx, ctxKeyDebugWriter, w)
}

func DebugWriterFrom(ctx context.Context) io.Writer {
	if w, ok := ctx.Value(ctxKeyDebugWriter).(io.Writer); ok {
		return w
	}
	return DebugWriter
}
