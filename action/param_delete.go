package action

import (
	"strings"

	"github.com/bryanl/woowoo/component"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// ParamDeleteOpt is an option for configuration ParamDelete.
type ParamDeleteOpt func(*ParamDelete)

// ParamDeleteWithIndex sets the index for the delete option.
func ParamDeleteWithIndex(index int) ParamDeleteOpt {
	return func(paramDelete *ParamDelete) {
		paramDelete.index = index
	}
}

// ParamDelete deletes a parameter from a component.
type ParamDelete struct {
	componentName string
	rawPath       string
	index         int

	*base
}

// NewParamDelete creates an instance of ParamDelete.
func NewParamDelete(fs afero.Fs, componentName, path string, opts ...ParamDeleteOpt) (*ParamDelete, error) {
	b, err := new(fs)
	if err != nil {
		return nil, err
	}

	paramDelete := &ParamDelete{
		componentName: componentName,
		rawPath:       path,
		base:          b,
	}

	for _, opt := range opts {
		opt(paramDelete)
	}

	return paramDelete, nil
}

// Run runs the action.
func (ps *ParamDelete) Run() error {
	path := strings.Split(ps.rawPath, ".")

	c, err := component.ExtractComponent(ps.fs, ps.pluginEnv.AppDir, ps.componentName)
	if err != nil {
		return errors.Wrap(err, "could not find component")
	}

	options := component.ParamOptions{
		Index: ps.index,
	}
	if err := c.DeleteParam(path, options); err != nil {
		return errors.Wrap(err, "delete param")
	}

	return nil

}
