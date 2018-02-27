package ksplugin

import (
	"os"

	"github.com/pkg/errors"
)

// PluginEnv is the plugin environment.
type PluginEnv struct {
	AppDir string
}

// Read reads the plugin environment.
func Read() (PluginEnv, error) {
	ksAppDir := os.Getenv("KS_APP_DIR")
	if ksAppDir == "" {
		return PluginEnv{}, errors.New("cannot find ks application directory")
	}

	return PluginEnv{
		AppDir: ksAppDir,
	}, nil
}
