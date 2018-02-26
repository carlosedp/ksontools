package k8sutil

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// FlattenToV1 expands any List-type objects into their members, and
// coerces everything to v1.Unstructured.
func FlattenToV1(objs []runtime.Object) ([]*unstructured.Unstructured, error) {
	ret := make([]*unstructured.Unstructured, 0, len(objs))
	for _, obj := range objs {
		switch o := obj.(type) {
		case *unstructured.UnstructuredList:
			for i := range o.Items {
				ret = append(ret, &o.Items[i])
			}
		case *unstructured.Unstructured:
			ret = append(ret, o)
		default:
			return nil, errors.New("Unexpected object type")
		}
	}
	return ret, nil
}
