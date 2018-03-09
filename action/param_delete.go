package action

import (
	"strings"

	"github.com/bryanl/woowoo/component"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// ParamDelete deletes a parameter from a component.
func ParamDelete(fs afero.Fs, componentName, path string, opts ...ParamDeleteOpt) error {
	pd, err := newParamDelete(fs, componentName, path, opts...)
	if err != nil {
		return err
	}

	return pd.Run()
}

// ParamDeleteOpt is an option for configuration ParamDelete.
type ParamDeleteOpt func(*paramDelete)

// ParamDeleteWithIndex sets the index for the delete option.
func ParamDeleteWithIndex(index int) ParamDeleteOpt {
	return func(pd *paramDelete) {
		pd.index = index
	}
}

// ParamDelete deletes a parameter from a component.
type paramDelete struct {
	componentName string
	rawPath       string
	index         int

	*base
}

// newParamDelete creates an instance of ParamDelete.
func newParamDelete(fs afero.Fs, componentName, path string, opts ...ParamDeleteOpt) (*paramDelete, error) {
	b, err := new(fs)
	if err != nil {
		return nil, err
	}

	pd := &paramDelete{
		componentName: componentName,
		rawPath:       path,
		base:          b,
	}

	for _, opt := range opts {
		opt(pd)
	}

	return pd, nil
}

// Run runs the action.
func (pd *paramDelete) Run() error {
	path := strings.Split(pd.rawPath, ".")

	c, err := component.ExtractComponent(pd.app, pd.componentName)
	if err != nil {
		return errors.Wrap(err, "could not find component")
	}

	options := component.ParamOptions{
		Index: pd.index,
	}
	if err := c.DeleteParam(path, options); err != nil {
		return errors.Wrap(err, "delete param")
	}

	return nil

}
