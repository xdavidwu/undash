package http

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"k8s.io/apimachinery/pkg/runtime"

	undashcontext "github.com/xdavidwu/undash/internal/context"
)

type JSONHandler[T any] func(w http.ResponseWriter, r *http.Request) (T, error)

func (j JSONHandler[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()

	res, err := j(w, r)
	if err != nil {
		WriteErrorAsMetaV1Status(w, err)
		return
	}

	if headers.Get("Content-Type") == "" {
		headers.Set("Content-Type", runtime.ContentTypeJSON)
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(res)
}

func InjectLogger(h http.Handler, base *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := rand.Text()
		logger := base.With("request.id", id)
		ctx := undashcontext.WithLogger(r.Context(), logger)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	length     int64
}

func (l *loggingResponseWriter) Write(b []byte) (int, error) {
	length, err := l.ResponseWriter.Write(b)
	l.length += int64(length)

	return length, err
}

func (l *loggingResponseWriter) WriteHeader(statusCode int) {
	l.ResponseWriter.WriteHeader(statusCode)
	l.statusCode = statusCode
}

func (l *loggingResponseWriter) StatusCode(ctx context.Context) int {
	if ctx.Err() != nil {
		return 499
	}
	if l.statusCode == 0 {
		return http.StatusOK
	}
	return l.statusCode
}

func AccessLog(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		before := time.Now()
		ctx := r.Context()
		writer := &loggingResponseWriter{ResponseWriter: w}
		logger := undashcontext.GetLogger(ctx)
		requestInfo := slog.Group("request",
			"method", r.Method,
			"url", r.URL.String(),
			"from", r.RemoteAddr,
			"at", before,
		)

		h.ServeHTTP(writer, r)

		code := writer.StatusCode(ctx)
		duration := time.Since(before)
		logger.InfoContext(ctx, "<=", requestInfo,
			"request.duration", duration,
			"response.code", code,
			"response.length", writer.length,
		)
	})
}
