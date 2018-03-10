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

func ParamSetGlobal(isGlobal bool) ParamSetOpt {
	return func(paramSet *paramSet) {
		paramSet.global = isGlobal
	}
}

// ParamSetWithIndex sets the index for the set option.
func ParamSetWithIndex(index int) ParamSetOpt {
	return func(paramSet *paramSet) {
		paramSet.index = index
	}
}

// ParamSet sets a parameter for a component.
type paramSet struct {
	name     string
	rawPath  string
	rawValue string
	index    int
	global   bool

	*base
}

// NewParamSet creates an instance of ParamSet.
func newParamSet(fs afero.Fs, name, path, value string, opts ...ParamSetOpt) (*paramSet, error) {
	b, err := new(fs)
	if err != nil {
		return nil, err
	}

	ps := &paramSet{
		name:     name,
		rawPath:  path,
		rawValue: value,
		base:     b,
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

	if ps.global {
		return ps.setGlobal(path, value)
	}

	return ps.setLocal(path, value)
}

func (ps *paramSet) setGlobal(path []string, value interface{}) error {
	ns, err := component.GetNamespace(ps.app, ps.name)
	if err != nil {
		return errors.Wrap(err, "retrieve namespace")
	}

	if err := ns.SetParam(path, value); err != nil {
		return errors.Wrap(err, "set global param")
	}

	return nil
}

func (ps *paramSet) setLocal(path []string, value interface{}) error {
	c, err := component.ExtractComponent(ps.app, ps.name)
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
