package kubernetes

import (
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
