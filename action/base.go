package action

import (
	"github.com/bryanl/woowoo/ksplugin"
	"github.com/bryanl/woowoo/ksutil"
	"github.com/spf13/afero"
)

type base struct {
	app ksutil.SuperApp
}

func new(fs afero.Fs) (*base, error) {
	pluginEnv, err := ksplugin.Read()
	if err != nil {
		return nil, err
	}

	app, err := ksutil.LoadApp(fs, pluginEnv.AppDir)
	if err != nil {
		return nil, err
	}

	return &base{
		app: app,
	}, nil
}
