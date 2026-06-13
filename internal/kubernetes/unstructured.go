package kubernetes

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func Structurize[T Object](
	o *unstructured.Unstructured, to T,
) error {
	return runtime.DefaultUnstructuredConverter.FromUnstructured(
		o.UnstructuredContent(), to)
}

func StructurizeList[T List](
	o *unstructured.UnstructuredList, to T,
) error {
	return runtime.DefaultUnstructuredConverter.FromUnstructured(
		o.UnstructuredContent(), to)
}

func UnstructuredListToTableFunc[
	ListTypeStruct any,
	ListType interface {
		*ListTypeStruct
		List
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
