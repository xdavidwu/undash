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

	status := MetaV1Status(r.StatusCode, string(body))
	b, err := json.Marshal(status)
	if err != nil {
		return fmt.Errorf("cannot marshal transformed metav1.Status: %w", err)
	}

	r.Header.Set("Content-Length", strconv.Itoa(len(b)))
	r.Header.Set("Content-Type", runtime.ContentTypeJSON)
	r.Body = io.NopCloser(bytes.NewBuffer(b))

	return nil
}
