package action

import (
	"github.com/bryanl/woowoo/ksplugin"
	"github.com/ksonnet/ksonnet/metadata/app"
	"github.com/spf13/afero"
)

type base struct {
	app app.App
}

func new(fs afero.Fs) (*base, error) {
	pluginEnv, err := ksplugin.Read()
	if err != nil {
		return nil, err
	}

	a, err := app.Load(fs, pluginEnv.AppDir)
	if err != nil {
		return nil, err
	}

	return &base{
		app: a,
	}, nil
}
