package action

import (
	"strings"

	"github.com/bryanl/woowoo/component"
	"github.com/bryanl/woowoo/params"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// ParamSetGlobal sets a parameter for a namespace.
type ParamSetGlobal struct {
	nsName   string
	rawPath  string
	rawValue string

	*base
}

// NewParamSetGlobal creates an instance of ParamSetGlobal.
func NewParamSetGlobal(fs afero.Fs, nsName, path, value string) (*ParamSetGlobal, error) {
	b, err := new(fs)
	if err != nil {
		return nil, err
	}

	paramSetGlobal := &ParamSetGlobal{
		nsName:   nsName,
		rawPath:  path,
		rawValue: value,
		base:     b,
	}

	return paramSetGlobal, nil
}

// Run runs the action.
func (psg *ParamSetGlobal) Run() error {
	path := strings.Split(psg.rawPath, ".")

	value, err := params.DecodeValue(psg.rawValue)
	if err != nil {
		return errors.Wrap(err, "value is invalid")
	}

	ns, err := component.GetNamespace(psg.fs, psg.pluginEnv.AppDir, psg.nsName)
	if err != nil {
		return errors.Wrap(err, "retrieve namespace")
	}

	if err := ns.SetParam(path, value); err != nil {
		return errors.Wrap(err, "set global param")
	}

	return nil
}
