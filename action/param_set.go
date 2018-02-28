package action

import (
	"strings"

	"github.com/bryanl/woowoo/component"
	"github.com/bryanl/woowoo/params"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// ParamSet sets a parameter for a component.
type ParamSet struct {
	componentName string
	rawPath       string
	rawValue      string

	*base
}

// NewParamSet creates an instance of ParamSet.
func NewParamSet(fs afero.Fs, componentName, path, value string) (*ParamSet, error) {
	b, err := new(fs)
	if err != nil {
		return nil, err
	}

	paramSet := &ParamSet{
		componentName: componentName,
		rawPath:       path,
		rawValue:      value,
		base:          b,
	}

	return paramSet, nil
}

// Run runs the action.
func (ps *ParamSet) Run() error {
	path := strings.Split(ps.rawPath, ".")

	value, err := params.DecodeValue(ps.rawValue)
	if err != nil {
		return errors.Wrap(err, "value is invalid")
	}

	c, err := component.ExtractComponent(ps.fs, ps.pluginEnv.AppDir, ps.componentName)
	if err != nil {
		return errors.Wrap(err, "could not find component")
	}

	if err := c.SetParam(path, value, component.ParamOptions{}); err != nil {
		return errors.Wrap(err, "set param")
	}

	return nil

}
