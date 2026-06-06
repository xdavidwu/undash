package http

import (
	"encoding/json"
	"net/http"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type JSONHandler[T any] func(w http.ResponseWriter, r *http.Request) (T, error)

func (j JSONHandler[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()

	res, err := j(w, r)
	if err != nil {
		var errRes metav1.Status

		if apiStatus, ok := err.(apierrors.APIStatus); ok {
			errRes = apiStatus.Status()
		} else {
			errRes = metav1.Status{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Status",
				},
				Status:  metav1.StatusFailure,
				Message: err.Error(),
			}
		}
		headers.Set("Content-Type", runtime.ContentTypeJSON)

		if errRes.Code != 0 {
			w.WriteHeader(int(errRes.Code))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		encoder := json.NewEncoder(w)
		encoder.Encode(errRes)
	}

	if headers.Get("Content-Type") == "" {
		headers.Set("Content-Type", runtime.ContentTypeJSON)
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(res)
}
