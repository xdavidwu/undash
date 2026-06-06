package context

import (
	"context"
	"log/slog"
)

type loggerKey struct{}

func WithLogger(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, l)
}

func GetLogger(ctx context.Context) *slog.Logger {
	v := ctx.Value(loggerKey{})
	if logger, ok := v.(*slog.Logger); ok && logger != nil {
		return logger
	} else {
		slog.Warn("context.GetLogger called where not injected")
		return slog.Default()
	}
}
