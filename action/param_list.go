package action

import (
	"os"

	"github.com/bryanl/woowoo/component"
	"github.com/bryanl/woowoo/ksutil"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

func ParamList(fs afero.Fs, nsName string) error {
	pl, err := newParamList(fs, nsName)
	if err != nil {
		return err
	}

	return pl.run()
}

type paramList struct {
	nsName string

	*base
}

func newParamList(fs afero.Fs, nsName string) (*paramList, error) {
	b, err := new(fs)
	if err != nil {
		return nil, err
	}

	pl := &paramList{
		nsName: nsName,
		base:   b,
	}

	return pl, nil
}

func (pl *paramList) run() error {
	ns, err := component.GetNamespace(pl.fs, pl.pluginEnv.AppDir, pl.nsName)
	if err != nil {
		return errors.Wrap(err, "could not find namespace")
	}

	paramData, err := ns.Params()
	if err != nil {
		return errors.Wrap(err, "could not list parameters")
	}

	table := ksutil.NewTable(os.Stdout)

	table.SetHeader([]string{"COMPONENT", "INDEX", "KEY", "VALUE"})
	for _, data := range paramData {
		table.Append([]string{data.Component, data.Index, data.Key, data.Value})
	}

	table.Render()

	return nil
}
