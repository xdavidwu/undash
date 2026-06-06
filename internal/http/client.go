package http

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	undashcontext "github.com/xdavidwu/undash/internal/context"
)

type loggingRoundTripper struct {
	http.RoundTripper
}

func (l *loggingRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	ctx := r.Context()
	logger := undashcontext.GetLogger(ctx)
	before := time.Now()

	res, err := l.RoundTripper.RoundTrip(r)

	duration := time.Since(before)
	requestInfo := slog.Group("request",
		"method", r.Method,
		"url", r.URL.String(),
		"at", before,
	)

	if err != nil {
		logger.ErrorContext(ctx, "=>!", requestInfo, "error", err.Error())
	} else {
		responseInfo := slog.Group(
			"response",
			"code", res.StatusCode,
			"header_duration", duration,
		)
		logger.InfoContext(ctx, "=>", requestInfo, responseInfo)
	}

	return res, err
}

func RequestLog(r http.RoundTripper) http.RoundTripper {
	return &loggingRoundTripper{r}
}

type RawClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	raw RawClient
}

func NewClient(raw RawClient) *Client {
	return &Client{raw: raw}
}

func (c *Client) Raw() RawClient {
	return c.raw
}

func (c *Client) Call(
	ctx context.Context, method string, url string, body io.Reader,
) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("cannot construct request: %w", err)
	}
	return c.raw.Do(req)
}
