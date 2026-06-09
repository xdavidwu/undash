package kubernetes

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/runtime"
	apiscorev1 "k8s.io/kubernetes/pkg/apis/core/v1"
	"k8s.io/kubernetes/pkg/printers"
	internalprinters "k8s.io/kubernetes/pkg/printers/internalversion"
)

var (
	printer = printers.NewTableGenerator()
)

func init() {
	internalprinters.AddHandlers(printer)
}

func ToTableFunc[
	APIType runtime.Object,
	InternalTypeStruct any,
	InternalType interface {
		*InternalTypeStruct
		runtime.Object
	},
](
	converter func(in APIType, out InternalType, s conversion.Scope) error,
) func(APIType) (*metav1.Table, error) {
	return func(obj APIType) (*metav1.Table, error) {
		var internalObj InternalType = new(InternalTypeStruct)
		if err := converter(obj, internalObj, nil); err != nil {
			return nil, fmt.Errorf("cannot convert object to kubernetes internal version: %w", err)
		}

		table, err := printer.GenerateTable(internalObj, printers.GenerateOptions{
			Wide: true,
		})
		if err != nil {
			return nil, err
		}

		table.SetGroupVersionKind(metav1.SchemeGroupVersion.WithKind("Table"))
		for i := range table.Rows {
			row := &table.Rows[i]

			// TODO full object support (convert back to api type?)
			metav1Obj, err := meta.Accessor(row.Object.Object)
			if err != nil {
				return nil, fmt.Errorf("cannot interpret internal version as metav1.Object: %w", err)
			}
			pom := meta.AsPartialObjectMetadata(metav1Obj)
			pom.SetGroupVersionKind(metav1.SchemeGroupVersion.WithKind("PartialObjectMetadata"))
			row.Object.Object = pom
		}
		return table, nil
	}
}

var (
	V1ServiceListToTable = ToTableFunc(apiscorev1.Convert_v1_ServiceList_To_core_ServiceList)
	// TODO more, perhaps via code gen
)
