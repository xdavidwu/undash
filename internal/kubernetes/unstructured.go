package kubernetes

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func Structurize[T metav1.ObjectMetaAccessor](
	o *unstructured.Unstructured, to T,
) error {
	return runtime.DefaultUnstructuredConverter.FromUnstructured(
		o.UnstructuredContent(), to)
}

func StructurizeList[T metav1.ListMetaAccessor](
	o *unstructured.UnstructuredList, to T,
) error {
	return runtime.DefaultUnstructuredConverter.FromUnstructured(
		o.UnstructuredContent(), to)
}

func UnstructuredListToTableFunc[
	ListTypeStruct any,
	ListType interface {
		*ListTypeStruct
		metav1.ListMetaAccessor
	},
](
	listToTable func(ListType) (*metav1.Table, error),
) func(*unstructured.UnstructuredList) (*metav1.Table, error) {
	return func(u *unstructured.UnstructuredList) (*metav1.Table, error) {
		var structured ListType = new(ListTypeStruct)
		err := StructurizeList(u, structured)
		if err != nil {
			return nil, fmt.Errorf("cannot convert list to api type: %w", err)
		}

		return listToTable(structured)
	}
}
