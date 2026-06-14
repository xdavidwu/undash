package kubernetes

import (
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Object interface {
	runtime.Object
	metav1.ObjectMetaAccessor
}

type List interface {
	runtime.Object
	metav1.ListMetaAccessor
}

func IsTableUnsupported(e error) bool {
	if e == nil {
		return false
	}

	return strings.Contains(e.Error(), "no table handler registered for this type ")
}
