package yaml2jsonnet

import (
	"strings"

	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/ksonnet"
	"github.com/pkg/errors"
)

var (
	// TODO: might need something in ksonnet lib to look this up
	groupMappings = map[string][]string{
		"apiextensions.k8s.io":      []string{"apiextensions"},
		"rbac.authorization.k8s.io": []string{"rbac"},
	}
)

// GVK is a group, version, kind descriptor.
type GVK struct {
	GroupPath []string
	Version   string
	Kind      string
}

func (gvk *GVK) Group() []string {
	g, ok := groupMappings[gvk.GroupPath[0]]
	if !ok {
		return gvk.GroupPath
	}

	return g
}

func (gvk *GVK) Path() []string {
	return append(gvk.Group(), gvk.Version, gvk.Kind)
}

// TypeSpec describes an object's type.
type TypeSpec map[string]string

// Validate validates the TypeSpec.
func (ts TypeSpec) Validate() error {
	if ts["kind"] == "" || ts["apiVersion"] == "" {
		return errors.Errorf("document doesn't describe a kubernetes object: %#v", ts)
	}

	return nil
}

// GVK returns the GVK descriptor for the TypeSpec.
func (ts TypeSpec) GVK() (GVK, error) {
	group, err := ts.Group()
	if err != nil {
		return GVK{}, err
	}

	version, err := ts.Version()
	if err != nil {
		return GVK{}, err
	}

	kind, err := ts.Kind()
	if err != nil {
		return GVK{}, err
	}

	return GVK{GroupPath: group, Version: version, Kind: kind}, nil
}

// Group is the group as defined by the TypeSpec.
func (ts TypeSpec) Group() ([]string, error) {
	if err := ts.Validate(); err != nil {
		return nil, err
	}

	parts := strings.Split(ts["apiVersion"], "/")
	if len(parts) == 1 {
		return []string{"core"}, nil
	}

	return []string{parts[0]}, nil
}

// Version is the version as defined by the TypeSpec.
func (ts TypeSpec) Version() (string, error) {
	if err := ts.Validate(); err != nil {
		return "", err
	}

	parts := strings.Split(ts["apiVersion"], "/")
	if len(parts) == 1 {
		return parts[0], nil
	}

	return parts[1], nil
}

// Kind is the kind as specified by the TypeSpec.
func (ts TypeSpec) Kind() (string, error) {
	if err := ts.Validate(); err != nil {
		return "", err
	}

	return ksonnet.FormatKind(ts["kind"]), nil
}
