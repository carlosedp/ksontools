package component

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

// Component is a ksonnet Component interface.
type Component interface {
	Objects() ([]*unstructured.Unstructured, error)
}
