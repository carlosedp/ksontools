package yaml2jsonnet

import (
	"strings"

	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/ksonnet"
	"github.com/pkg/errors"
)

type GVK struct {
	Group   string
	Version string
	Kind    string
}

type TypeSpec map[string]string

func (ts TypeSpec) Validate() error {
	if ts["kind"] == "" || ts["apiVersion"] == "" {
		return errors.Errorf("document doesn't describe a resource: %#v", ts)
	}

	return nil
}

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

	return GVK{Group: group, Version: version, Kind: kind}, nil
}

func (ts TypeSpec) Group() (string, error) {
	if err := ts.Validate(); err != nil {
		return "", err
	}

	parts := strings.Split(ts["apiVersion"], "/")
	if len(parts) == 1 {
		return "core", nil
	}

	return parts[0], nil
}

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

func (ts TypeSpec) Kind() (string, error) {
	if err := ts.Validate(); err != nil {
		return "", err
	}

	return ksonnet.FormatKind(ts["kind"]), nil
}
