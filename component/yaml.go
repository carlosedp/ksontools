package component

import (
	"io"

	"github.com/go-yaml/yaml"
	"github.com/pkg/errors"
)

// ImportYaml converts a reader containing YAML to a TypeSpec and Properties.
func ImportYaml(r io.Reader) (*TypeSpec, Properties, error) {
	var m map[interface{}]interface{}
	if err := yaml.NewDecoder(r).Decode(&m); err != nil {
		return nil, nil, errors.Wrap(err, "decode yaml")
	}

	props := Properties{}

	var kind string
	var apiVersion string

	for k, v := range m {
		switch k {
		case "apiVersion":
			apiVersion = v.(string)
		case "kind":
			kind = v.(string)
		default:
			props[k] = v
		}
	}

	ts, err := NewTypeSpec(apiVersion, kind)
	if err != nil {
		return nil, nil, err
	}

	return ts, props, nil
}
