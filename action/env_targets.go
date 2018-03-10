package action

import (
	"github.com/spf13/afero"
)

// EnvTargets sets targets for an environment.
func EnvTargets(fs afero.Fs, envName string, components []string) error {
	et, err := newEnvTargets(fs, envName, components)
	if err != nil {
		return err
	}

	return et.Run()
}

type envTargets struct {
	envName    string
	components []string

	*base
}

func newEnvTargets(fs afero.Fs, envName string, components []string) (*envTargets, error) {
	b, err := new(fs)
	if err != nil {
		return nil, err
	}

	et := &envTargets{
		envName:    envName,
		components: components,
		base:       b,
	}

	return et, nil
}

func (et *envTargets) Run() error {
	return et.app.UpdateTargets(et.envName, et.components)
}
