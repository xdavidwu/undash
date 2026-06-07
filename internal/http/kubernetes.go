package http

import (
	"encoding/json"
	"net/http"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func MetaV1Status(code int, message string) metav1.Status {
	return metav1.Status{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Status",
		},
		Status:  metav1.StatusFailure,
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
