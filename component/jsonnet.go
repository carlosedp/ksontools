package component

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/bryanl/woowoo/ksutil"
	"github.com/bryanl/woowoo/params"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// Jsonnet is a component base on jsonnet.
type Jsonnet struct {
	app        ksutil.SuperApp
	source     string
	paramsPath string
}

var _ Component = (*Jsonnet)(nil)

// NewJsonnet creates an instance of Jsonnet.
func NewJsonnet(app ksutil.SuperApp, source, paramsPath string) *Jsonnet {
	return &Jsonnet{
		app:        app,
		source:     source,
		paramsPath: paramsPath,
	}
}

// Name is the name of this component.
func (j *Jsonnet) Name() string {
	base := filepath.Base(j.source)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

func (j *Jsonnet) Objects(paramsStr string) ([]*unstructured.Unstructured, error) {
	return nil, errors.Errorf("not here yet")
}

// SetParam set parameter for a component.
func (j *Jsonnet) SetParam(path []string, value interface{}, options ParamOptions) error {
	paramsData, err := j.readParams()
	if err != nil {
		return err
	}

	updatedParams, err := params.Set(path, paramsData, j.Name(), value, paramsComponentRoot)
	if err != nil {
		return err
	}

	if err = j.writeParams(updatedParams); err != nil {
		return err
	}

	return nil
}

// DeleteParam deletes a param.
func (j *Jsonnet) DeleteParam(path []string, options ParamOptions) error {
	paramsData, err := j.readParams()
	if err != nil {
		return err
	}

	updatedParams, err := params.Delete(path, paramsData, j.Name(), paramsComponentRoot)
	if err != nil {
		return err
	}

	if err = j.writeParams(updatedParams); err != nil {
		return err
	}

	return nil
}

// Params returns params for a component.
func (j *Jsonnet) Params() ([]NamespaceParameter, error) {
	paramsData, err := j.readParams()
	if err != nil {
		return nil, err
	}

	props, err := params.ToMap(j.Name(), paramsData, paramsComponentRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not find components")
	}

	var params []NamespaceParameter
	for k, v := range props {
		vStr, err := j.paramValue(v)
		if err != nil {
			return nil, err
		}
		np := NamespaceParameter{
			Component: j.Name(),
			Key:       k,
			Index:     "0",
			Value:     vStr,
		}

		params = append(params, np)
	}

	sort.Slice(params, func(i, j int) bool {
		return params[i].Key < params[j].Key
	})

	return params, nil
}

func (j *Jsonnet) paramValue(v interface{}) (string, error) {
	switch v.(type) {
	default:
		s := fmt.Sprintf("%v", v)
		return s, nil
	case string:
		s := fmt.Sprintf("%v", v)
		return strconv.Quote(s), nil
	case map[string]interface{}, []interface{}:
		b, err := json.Marshal(&v)
		if err != nil {
			return "", err
		}

		return string(b), nil
	}
}

// Summarize creates a summary for the component.
func (j *Jsonnet) Summarize() ([]Summary, error) {
	return []Summary{
		{
			ComponentName: j.Name(),
			IndexStr:      "0",
			Type:          "jsonnet",
		},
	}, nil
}

func (j *Jsonnet) readParams() (string, error) {
	b, err := afero.ReadFile(j.app.Fs(), j.paramsPath)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (j *Jsonnet) writeParams(src string) error {
	return afero.WriteFile(j.app.Fs(), j.paramsPath, []byte(src), 0644)
}
