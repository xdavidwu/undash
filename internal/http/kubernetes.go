package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	undashctx "github.com/xdavidwu/undash/internal/context"
	"github.com/xdavidwu/undash/internal/kubernetes"
)

var (
	statusReasonFromCode = map[int]metav1.StatusReason{
		http.StatusUnauthorized: metav1.StatusReasonUnauthorized,
		http.StatusForbidden:    metav1.StatusReasonForbidden,
		http.StatusNotFound:     metav1.StatusReasonNotFound,
		http.StatusConflict:     metav1.StatusReasonConflict,
		http.StatusGone:         metav1.StatusReasonGone,
		// metav1.StatusReasonInvalid is more specific than http.StatusUnprocessableEntity
		// metav1.StatusReasonTimeout is more specific than http.StatusGatewayTimeout
		http.StatusTooManyRequests: metav1.StatusReasonTooManyRequests,
		// metav1.StatusReasonBadRequest is more specific than http.StatusBadRequest
		http.StatusMethodNotAllowed:      metav1.StatusReasonMethodNotAllowed,
		http.StatusNotAcceptable:         metav1.StatusReasonNotAcceptable,
		http.StatusRequestEntityTooLarge: metav1.StatusReasonRequestEntityTooLarge,
		http.StatusUnsupportedMediaType:  metav1.StatusReasonUnsupportedMediaType,
		http.StatusInternalServerError:   metav1.StatusReasonInternalError,
		http.StatusServiceUnavailable:    metav1.StatusReasonServiceUnavailable,
	}
)

func MetaV1Status(code int, message string) metav1.Status {
	reason, ok := statusReasonFromCode[code]
	if !ok {
		if code >= 500 {
			reason = metav1.StatusReasonInternalError
		} else {
			reason = metav1.StatusReasonUnknown
		}
	}

	return metav1.Status{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Status",
		},
		Status:  metav1.StatusFailure,
		Reason:  reason,
		Message: message,
		Code:    int32(code),
	}
}

func ErrorToMetaV1Status(err error) metav1.Status {
	if apiStatus, ok := err.(apierrors.APIStatus); ok {
		return apiStatus.Status()
	}

	return MetaV1Status(http.StatusInternalServerError, err.Error())
}

func WriteErrorAsMetaV1Status(w http.ResponseWriter, err error) {
	status := ErrorToMetaV1Status(err)
	w.Header().Set("Content-Type", runtime.ContentTypeJSON)

	if status.Code != 0 {
		w.WriteHeader(int(status.Code))
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(status)
}

func RewriteObjectAsTableIfRequested(r *http.Response) error {
	if r.StatusCode != http.StatusOK {
		return nil
	}
	if !AcceptsExactMediaType(r.Request, MetaV1TableJSON) {
		return nil
	}

	b, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return err
	}

	obj := &unstructured.Unstructured{}
	if err := json.Unmarshal(b, obj); err != nil {
		return err
	}

	gvk := obj.GroupVersionKind()
	if toTable, ok := kubernetes.UnstructuredToTableFuncs[gvk]; ok {
		table, err := toTable(obj)
		if err != nil {
			return fmt.Errorf("cannot convert object to table: %w", err)
		}

		return rewriteToJSON(r, table, MetaV1TableJSON.String())
	}

	ctx := r.Request.Context()
	undashctx.GetLogger(ctx).WarnContext(ctx, "object table func not registered", "gvk", gvk.String())

	r.Body = io.NopCloser(bytes.NewBuffer(b))
	return nil
}
