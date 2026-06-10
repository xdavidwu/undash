package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"k8s.io/apimachinery/pkg/runtime"
)

func rewriteToJSON(r *http.Response, obj any, mediaType string) error {
	b, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("cannot marshal transformed response: %w", err)
	}

	r.Header.Set("Content-Length", strconv.Itoa(len(b)))
	r.Header.Set("Content-Type", mediaType)
	r.Body = io.NopCloser(bytes.NewBuffer(b))

	return nil
}

// Rewrites plain-text error responses as [metav1.Status],
// for [httputil.ReverseProxy.ModifyResponse].
// Expected to be paired with [NoExplicitCompression] - no compression is handled here.
func ErrorResponseAsMetaV1Status(r *http.Response) error {
	if r.StatusCode < 400 {
		return nil
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	r.Body.Close()

	return rewriteToJSON(
		r,
		MetaV1Status(r.StatusCode, string(body)),
		runtime.ContentTypeJSON,
	)
}

func ChainModifyResponse(m ...func(*http.Response) error) func(*http.Response) error {
	return func(r *http.Response) error {
		for _, fn := range m {
			if err := fn(r); err != nil {
				return err
			}
		}
		return nil
	}
}
