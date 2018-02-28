package action

import (
	"github.com/bryanl/woowoo/ksplugin"
	"github.com/spf13/afero"
)

type base struct {
	fs        afero.Fs
	pluginEnv ksplugin.PluginEnv
}

func new(fs afero.Fs) (*base, error) {
	pluginEnv, err := ksplugin.Read()
	if err != nil {
		return nil, err
	}

	return &base{
		fs:        fs,
		pluginEnv: pluginEnv,
	}, nil
}
