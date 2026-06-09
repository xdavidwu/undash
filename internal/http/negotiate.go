package http

import (
	"mime"
	"net/http"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type MediaType struct {
	TypeSubType string
	Parameters  map[string]string
}

func (m *MediaType) String() string {
	return mime.FormatMediaType(m.TypeSubType, m.Parameters)
}

func AsGVK(gvk schema.GroupVersionKind) map[string]string {
	return map[string]string{
		"as": gvk.Kind,
		"v":  gvk.Version,
		"g":  gvk.Group,
	}
}

var (
	MetaV1TableJSON = MediaType{
		TypeSubType: runtime.ContentTypeJSON,
		Parameters:  AsGVK(metav1.SchemeGroupVersion.WithKind("Table")),
	}
)

// use of as= make most impl not really useful since what kubernetes want here is a little bit special:
// we can't feed application/json;as=Table;v=v1;g=meta.k8s.io to client wanting application/json

// Find if a request explicitly accepts [mediaType] containing all its [MediaType.Params]
func AcceptsExactMediaType(r *http.Request, target MediaType) bool {
	for _, accept := range r.Header.Values("Accept") {
	ForMediaType:
		for mediaType := range strings.SplitSeq(accept, ",") {
			typeSubType, params, err := mime.ParseMediaType(strings.TrimSpace(mediaType))
			if err != nil {
				continue
			}

			if typeSubType == target.TypeSubType {
				for k, v := range target.Parameters {
					if params[k] != v {
						continue ForMediaType
					}
				}
				return true
			}
		}
	}
	return false
}
