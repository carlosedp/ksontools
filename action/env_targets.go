package action

import (
	"github.com/bryanl/woowoo/component"
	"github.com/spf13/afero"
)

// EnvTargets sets targets for an environment.
func EnvTargets(fs afero.Fs, envName string, nsNames []string) error {
	et, err := newEnvTargets(fs, envName, nsNames)
	if err != nil {
		return err
	}

	return et.Run()
}

type envTargets struct {
	envName string
	nsNames []string

	*base
}

func newEnvTargets(fs afero.Fs, envName string, nsNames []string) (*envTargets, error) {
	b, err := new(fs)
	if err != nil {
		return nil, err
	}

	et := &envTargets{
		envName: envName,
		nsNames: nsNames,
		base:    b,
	}

	return et, nil
}

func (et *envTargets) Run() error {
	for _, nsName := range et.nsNames {
		_, err := component.GetNamespace(et.app, nsName)
		if err != nil {
			return err
		}
	}

	return et.app.UpdateTargets(et.envName, et.nsNames)
}
