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

// Error type for status code check from [CheckStatusCode].
// [Response] is passed as-is without reading body, until [StatusCodeError.Error],
// which drains the body into unexported fields.
type StatusCodeError struct {
	*http.Response
	body    []byte
	bodyErr error
}

func (s *StatusCodeError) Error() string {
	if s.body == nil && s.bodyErr == nil {
		s.body, s.bodyErr = io.ReadAll(s.Response.Body)
		s.Response.Body.Close()
	}

	if s.bodyErr != nil {
		return fmt.Sprintf("server returned status %d, cannot read response: %s", s.Response.StatusCode, s.bodyErr)
	}

	return fmt.Sprintf("server returned status %d: %s", s.Response.StatusCode, string(s.body))
}

type checkStatusCodeRawClient struct {
	RawClient
}

func (c *checkStatusCodeRawClient) Do(req *http.Request) (*http.Response, error) {
	r, err := c.RawClient.Do(req)
	if err != nil {
		return nil, err
	}

	if r.StatusCode >= 400 {
		return nil, &StatusCodeError{Response: r}
	}

	return r, nil
}

func CheckStatusCode(c RawClient) RawClient {
	return &checkStatusCodeRawClient{c}
}

type Client struct {
	raw RawClient
}

func NewClient(raw RawClient) *Client {
	return &Client{raw: raw}
}

func NewDefaultClient() *Client {
	return NewClient(CheckStatusCode(&http.Client{
		Transport: RequestLog(http.DefaultTransport),
	}))
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
