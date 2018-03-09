package action

import (
	"strings"

	"github.com/bryanl/woowoo/component"
	"github.com/bryanl/woowoo/params"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// ParamSet sets a parameter for a component.
func ParamSet(fs afero.Fs, componentName, path, value string, opts ...ParamSetOpt) error {
	ps, err := newParamSet(fs, componentName, path, value, opts...)
	if err != nil {
		return err
	}

	return ps.Run()
}

// ParamSetOpt is an option for configuring ParamSet.
type ParamSetOpt func(*paramSet)

// ParamSetWithIndex sets the index for the set option.
func ParamSetWithIndex(index int) ParamSetOpt {
	return func(paramSet *paramSet) {
		paramSet.index = index
	}
}

// ParamSet sets a parameter for a component.
type paramSet struct {
	componentName string
	rawPath       string
	rawValue      string
	index         int

	*base
}

// NewParamSet creates an instance of ParamSet.
func newParamSet(fs afero.Fs, componentName, path, value string, opts ...ParamSetOpt) (*paramSet, error) {
	b, err := new(fs)
	if err != nil {
		return nil, err
	}

	ps := &paramSet{
		componentName: componentName,
		rawPath:       path,
		rawValue:      value,
		base:          b,
	}

	for _, opt := range opts {
		opt(ps)
	}

	return ps, nil
}

// Run runs the action.
func (ps *paramSet) Run() error {
	path := strings.Split(ps.rawPath, ".")

	value, err := params.DecodeValue(ps.rawValue)
	if err != nil {
		return errors.Wrap(err, "value is invalid")
	}

	c, err := component.ExtractComponent(ps.app, ps.componentName)
	if err != nil {
		return errors.Wrap(err, "could not find component")
	}

	options := component.ParamOptions{
		Index: ps.index,
	}
	if err := c.SetParam(path, value, options); err != nil {
		return errors.Wrap(err, "set param")
	}

	return nil

}
