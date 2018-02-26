package component

import (
	"bufio"
	"io"

	"github.com/bryanl/woowoo/k8sutil"
	"github.com/go-yaml/yaml"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	amyaml "k8s.io/apimachinery/pkg/util/yaml"
)

// ImportYaml converts a reader containing YAML to a TypeSpec and Properties.
func ImportYaml(r io.Reader) (*TypeSpec, Properties, error) {
	// TODO: use apimachinery yaml util
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

// YAML represents a YAML component.
type YAML struct {
	source string
	fs     afero.Fs
}

var _ Component = (*YAML)(nil)

// Objects converts YAML to a slice apimachinery Unstructured objects.
func (y *YAML) Objects() ([]*unstructured.Unstructured, error) {
	// vmFactory := jsonnetutil.VMFactory{}
	// vm, err := vmFactory.VM()
	// if err != nil {
	// 	return nil, errors.Wrap(err, "create jsonnet vm")
	// }

	objects, err := y.readObject()
	if err != nil {
		return nil, errors.Wrap(err, "read object")
	}

	list, err := k8sutil.FlattenToV1(objects)
	if err != nil {
		return nil, errors.Wrap(err, "flatten objects")
	}

	return list, nil
}

func (y *YAML) readObject() ([]runtime.Object, error) {
	f, err := y.fs.Open(y.source)
	if err != nil {
		return nil, err
	}

	decoder := amyaml.NewYAMLReader(bufio.NewReader(f))
	ret := []runtime.Object{}
	for {
		bytes, err := decoder.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		if len(bytes) == 0 {
			continue
		}
		jsondata, err := amyaml.ToJSON(bytes)
		if err != nil {
			return nil, err
		}
		obj, _, err := unstructured.UnstructuredJSONScheme.Decode(jsondata, nil, nil)
		if err != nil {
			return nil, err
		}
		ret = append(ret, obj)
	}
	return ret, nil
}
