package component

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/bryanl/woowoo/k8sutil"
	"github.com/bryanl/woowoo/params"
	utilyaml "github.com/bryanl/woowoo/pkg/util/yaml"
	"github.com/go-yaml/yaml"
	"github.com/ksonnet/ksonnet/metadata/app"
	jsonnetutil "github.com/ksonnet/ksonnet/pkg/util/jsonnet"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	amyaml "k8s.io/apimachinery/pkg/util/yaml"
)

const (
	paramsComponentRoot = "components"
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

// YAML represents a YAML component. Since JSON is a subset of YAML, it can handle JSON as well.
type YAML struct {
	app        app.App
	nsName     string
	source     string
	paramsPath string
}

var _ Component = (*YAML)(nil)

// NewYAML creates an instance of YAML.
func NewYAML(a app.App, nsName, source, paramsPath string) *YAML {
	return &YAML{
		app:        a,
		nsName:     nsName,
		source:     source,
		paramsPath: paramsPath,
	}
}

// Name is the component name.
func (y *YAML) Name(wantsNameSpaced bool) string {
	base := filepath.Base(y.source)
	name := strings.TrimSuffix(base, filepath.Ext(base))

	if !wantsNameSpaced {
		return name
	}

	return strings.TrimPrefix(path.Join(y.nsName, name), "/")
}

// Params returns params for a component.
func (y *YAML) Params() ([]NamespaceParameter, error) {
	libPath, err := y.app.LibPath("default")
	if err != nil {
		return nil, err
	}

	k8sPath := filepath.Join(libPath, "k8s.libsonnet")
	obj, err := jsonnetutil.ImportFromFs(k8sPath, y.app.Fs())
	if err != nil {
		return nil, err
	}

	ve := NewValueExtractor(obj)

	// find all the params for this component
	// keys will look like `component-id`
	paramsData, err := y.readParams()
	if err != nil {
		return nil, err
	}

	props, err := params.ToMap("", paramsData, paramsComponentRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not find components")
	}

	re, err := regexp.Compile(fmt.Sprintf(`^%s-(\d+)$`, y.Name(false)))
	if err != nil {
		return nil, err
	}

	readers, err := utilyaml.Decode(y.app.Fs(), y.source)
	if err != nil {
		return nil, err
	}

	var params []NamespaceParameter
	for componentName, componentValue := range props {
		matches := re.FindAllStringSubmatch(componentName, 1)
		if len(matches) > 0 {
			index := matches[0][1]
			i, err := strconv.Atoi(index)
			if err != nil {
				return nil, err
			}

			ts, props, err := ImportYaml(readers[i])
			if err != nil {
				return nil, err
			}

			valueMap, err := ve.Extract(ts.GVK(), props)
			if err != nil {
				return nil, err
			}

			m, ok := componentValue.(map[string]interface{})
			if !ok {
				return nil, errors.Errorf("component value for %q was not a map", componentName)
			}

			childParams, err := y.paramValues(y.Name(false), index, valueMap, m, nil)
			if err != nil {
				return nil, err
			}

			params = append(params, childParams...)
		}
	}

	return params, nil
}

func isLeaf(path []string, key string, valueMap map[string]Values) (string, bool) {
	childPath := strings.Join(append(path, key), ".")
	for _, v := range valueMap {
		if strings.Join(v.Lookup, ".") == childPath {
			return childPath, true
		}
	}

	return "", false
}

func (y *YAML) paramValues(componentName, index string, valueMap map[string]Values, m map[string]interface{}, path []string) ([]NamespaceParameter, error) {
	var params []NamespaceParameter

	for k, v := range m {
		var s string
		switch t := v.(type) {
		default:
			if childPath, exists := isLeaf(path, k, valueMap); exists {
				s = fmt.Sprintf("%v", v)
				p := NamespaceParameter{
					Component: componentName,
					Index:     index,
					Key:       childPath,
					Value:     s,
				}
				params = append(params, p)
			}

		case map[string]interface{}:
			if childPath, exists := isLeaf(path, k, valueMap); exists {
				b, err := json.Marshal(&v)
				if err != nil {
					return nil, err
				}
				s = string(b)
				p := NamespaceParameter{
					Component: componentName,
					Index:     index,
					Key:       childPath,
					Value:     s,
				}
				params = append(params, p)
			} else {
				childPath := append(path, k)
				childParams, err := y.paramValues(componentName, index, valueMap, t, childPath)
				if err != nil {
					return nil, err
				}

				params = append(params, childParams...)
			}
		case []interface{}:
			if childPath, exists := isLeaf(path, k, valueMap); exists {
				b, err := json.Marshal(&v)
				if err != nil {
					return nil, err
				}
				s = string(b)
				p := NamespaceParameter{
					Component: componentName,
					Index:     index,
					Key:       childPath,
					Value:     s,
				}
				params = append(params, p)
			}
		}
	}

	return params, nil
}

// Objects converts YAML to a slice of apimachinery Unstructured objects. Params for a YAML
// based component are keyed like, `name-id`, where `name` is the file name sans the extension,
// and the id is the position within the file (starting at 0). Params are named this way
// because a YAML file can contain more than one object.
func (y *YAML) Objects(paramsStr, envName string) ([]*unstructured.Unstructured, error) {
	if paramsStr == "" {
		dir := filepath.Dir(y.source)
		paramsFile := filepath.Join(dir, "params.libsonnet")

		b, err := afero.ReadFile(y.app.Fs(), paramsFile)
		if err != nil {
			return nil, err
		}

		paramsStr = string(b)
	}

	return y.raw(paramsStr)
}

// SetParam set parameter for a component.
func (y *YAML) SetParam(path []string, value interface{}, options ParamOptions) error {
	entry := fmt.Sprintf("%s-%d", y.Name(false), options.Index)
	paramsData, err := y.readParams()
	if err != nil {
		return err
	}

	updatedParams, err := params.Set(path, paramsData, entry, value, paramsComponentRoot)
	if err != nil {
		return err
	}

	if err = y.writeParams(updatedParams); err != nil {
		return err
	}

	return nil
}

// DeleteParam deletes a param.
func (y *YAML) DeleteParam(path []string, options ParamOptions) error {
	entry := fmt.Sprintf("%s-%d", y.Name(false), options.Index)
	paramsData, err := y.readParams()
	if err != nil {
		return err
	}

	updatedParams, err := params.Delete(path, paramsData, entry, paramsComponentRoot)
	if err != nil {
		return err
	}

	if err = y.writeParams(updatedParams); err != nil {
		return err
	}

	return nil
}

func (y *YAML) readParams() (string, error) {
	b, err := afero.ReadFile(y.app.Fs(), y.paramsPath)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (y *YAML) writeParams(src string) error {
	return afero.WriteFile(y.app.Fs(), y.paramsPath, []byte(src), 0644)
}

func (y *YAML) raw(paramsStr string) ([]*unstructured.Unstructured, error) {
	objects, err := y.readObject(paramsStr)
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

	exists, err := afero.Exists(y.app.Fs(), paramsFile)
	if err != nil || !exists {
		return false, nil
	}

	paramsObj, err := jsonnetutil.ImportFromFs(paramsFile, y.app.Fs())
	if err != nil {
		return false, errors.Wrap(err, "import params")
	}

	componentPath := []string{
		paramsComponentRoot,
		fmt.Sprintf("%s-0", y.Name(false)),
	}
	_, err = jsonnetutil.FindObject(paramsObj, componentPath)
	if err != nil {
		return false, nil
	}

	return true, nil
}

func (y *YAML) readObject(paramsStr string) ([]runtime.Object, error) {
	f, err := y.app.Fs().Open(y.source)
	if err != nil {
		return nil, err
	}

	base := strings.TrimSuffix(filepath.Base(y.source), filepath.Ext(y.source))

	decoder := amyaml.NewYAMLReader(bufio.NewReader(f))
	ret := []runtime.Object{}
	i := 0
	for {
		componentName := fmt.Sprintf("%s-%d", base, i)
		i++
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

		patched, err := patchJSON(string(jsondata), paramsStr, componentName)
		if err != nil {
			return nil, err
		}

		jsondata = []byte(patched)

		obj, _, err := unstructured.UnstructuredJSONScheme.Decode(jsondata, nil, nil)
		if err != nil {
			return nil, err
		}
		ret = append(ret, obj)
	}
	return ret, nil
}

// Summarize generates a summary for a YAML component. For each manifest, it will
// return a slice of summaries of resources described.
func (y *YAML) Summarize() ([]Summary, error) {
	var summaries []Summary

	readers, err := utilyaml.Decode(y.app.Fs(), y.source)
	if err != nil {
		return nil, err
	}

	for i, r := range readers {
		ts, props, err := ImportYaml(r)
		if err != nil {
			return nil, err
		}

		name, err := props.Name()
		if err != nil {
			return nil, err
		}

		summary := Summary{
			ComponentName: y.Name(false),
			IndexStr:      strconv.Itoa(i),
			Type:          y.ext(),
			APIVersion:    ts.apiVersion,
			Kind:          ts.kind,
			Name:          name,
		}
		summaries = append(summaries, summary)
	}

	return summaries, nil
}

func (y *YAML) ext() string {
	return strings.TrimPrefix(filepath.Ext(y.source), ".")

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
