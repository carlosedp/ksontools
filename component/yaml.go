package component

import (
	"bufio"
	"fmt"
	"io"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/bryanl/woowoo/jsonnetutil"
	"github.com/bryanl/woowoo/k8sutil"
	"github.com/bryanl/woowoo/params"
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
	source     string
	fs         afero.Fs
	paramsPath string
}

var _ Component = (*YAML)(nil)

// NewYAML creates an instance of YAML.
func NewYAML(fs afero.Fs, source, paramsPath string) *YAML {
	return &YAML{fs: fs, source: source, paramsPath: paramsPath}
}

// Name is the component name.
func (y *YAML) Name() string {
	return y.componentName()
}

// Objects converts YAML to a slice apimachinery Unstructured objects. Params for a YAML
// based component are keyed like, `name-id`, where `name` is the file name sans the extension,
// and the id is the position within the file (starting at 0). Params are named this way
// because a YAML file can contain more than one object.
func (y *YAML) Objects() ([]*unstructured.Unstructured, error) {
	isParams, err := y.hasParams()
	if err != nil {
		return nil, errors.Wrap(err, "unable to check for params")
	}

	if isParams {
		return y.applyParams()
	}

	return y.raw()
}

// SetParam set parameter for a component.
func (y *YAML) SetParam(path []string, value interface{}, options SetParamOptions) error {
	entry := fmt.Sprintf("%s-%d", y.componentName(), options.Index)
	paramsData, err := y.readParams()
	if err != nil {
		return err
	}

	props, err := params.ToMap(entry, paramsData)
	if err != nil {
		props = make(map[string]interface{})
	}

	changes := make(map[string]interface{})
	cur := changes

	for i, k := range path {
		if i == len(path)-1 {
			cur[k] = value
		} else {
			if _, ok := cur[k]; !ok {
				m := make(map[string]interface{})
				cur[k] = m
				cur = m
			}
		}
	}

	if err = mergeMaps(props, changes, nil); err != nil {
		return err
	}

	updatedParams, err := params.Update(entry, paramsData, changes)
	if err != nil {
		return err
	}

	if err = y.writeParams(updatedParams); err != nil {
		return err
	}

	return nil
}

func (y *YAML) readParams() (string, error) {
	b, err := afero.ReadFile(y.fs, y.paramsPath)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (y *YAML) writeParams(src string) error {
	return afero.WriteFile(y.fs, y.paramsPath, []byte(src), 0644)
}

func (y *YAML) applyParams() ([]*unstructured.Unstructured, error) {
	dir := filepath.Dir(y.source)
	paramsFile := filepath.Join(dir, "params.libsonnet")

	b, err := afero.ReadFile(y.fs, paramsFile)
	if err != nil {
		return nil, err
	}

	objects, err := y.raw()
	if err != nil {
		return nil, err
	}

	for i := range objects {
		cn := fmt.Sprintf("%s-%d", y.componentName(), i)
		m, err := params.ToMap(cn, string(b))
		if err != nil {
			return nil, err
		}

		err = mergeMaps(objects[i].Object, m, nil)
		if err != nil {
			return nil, err
		}
	}

	return objects, nil

	// vmFactory := jsonnetutil.VMFactory{}
	// vm, err := vmFactory.VM()
	// if err != nil {
	// 	return nil, errors.Wrap(err, "create jsonnet vm")
	// }

	// return nil, errors.New("not implemented")
}

func (y *YAML) raw() ([]*unstructured.Unstructured, error) {
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

func (y *YAML) hasParams() (bool, error) {
	dir := filepath.Dir(y.source)
	paramsFile := filepath.Join(dir, "params.libsonnet")

	exists, err := afero.Exists(y.fs, paramsFile)
	if err != nil || !exists {
		return false, nil
	}

	paramsObj, err := jsonnetutil.ImportFromFs(paramsFile, y.fs)
	if err != nil {
		return false, errors.Wrap(err, "import params")
	}

	componentPath := []string{
		"components",
		fmt.Sprintf("%s-0", y.componentName()),
	}
	_, err = jsonnetutil.FindObject(paramsObj, componentPath)
	if err != nil {
		return false, nil
	}

	return true, nil
}

func (y *YAML) componentName() string {
	base := filepath.Base(y.source)
	return strings.TrimSuffix(base, filepath.Ext(base))
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

type paramPath struct {
	path  []string
	value interface{}
}

func mapToPaths(m map[string]interface{}, lookup map[string]bool, parent []string) []paramPath {
	paths := make([]paramPath, 0)

	for k, v := range m {
		cur := append(parent, k)

		switch t := v.(type) {
		default:
			pp := paramPath{path: cur, value: v}
			paths = append(paths, pp)

		case map[string]interface{}:
			children := mapToPaths(t, lookup, cur)

			route := strings.Join(cur, ".")
			if _, ok := lookup[route]; ok {
				pp := paramPath{path: cur, value: v}
				paths = append(paths, pp)
			} else {
				paths = append(paths, children...)
			}

		}
	}

	sort.Slice(paths, func(i, j int) bool {
		a := strings.Join(paths[i].path, ".")
		b := strings.Join(paths[j].path, ".")

		return a < b
	})

	return paths
}

func mergeMaps(m1 map[string]interface{}, m2 map[string]interface{}, path []string) error {
	for k := range m2 {
		_, ok := m1[k]
		if ok {
			v1, isMap1 := m1[k].(map[string]interface{})
			v2, isMap2 := m2[k].(map[string]interface{})
			if isMap1 && isMap2 {
				err := mergeMaps(v1, v2, append(path, k))
				if err != nil {
					return err
				}
			} else if reflect.TypeOf(v1) == reflect.TypeOf(v2) {
				m1[k] = m2[k]
			} else {
				errorPath := append(path, k)
				return fmt.Errorf("not same types at %s", strings.Join(errorPath, "."))
			}
		} else {
			m1[k] = m2[k]
		}
	}

	return nil
}
