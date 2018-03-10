package action

import (
	"os"

	"github.com/go-yaml/yaml"
	"github.com/spf13/afero"
)

// EnvDescribe describes an environment by printing its configuration.
func EnvDescribe(fs afero.Fs, envName string) error {
	ed, err := newEnvDescribe(fs, envName)
	if err != nil {
		return err
	}

	return ed.Run()
}

type envDescribe struct {
	envName string

	*base
}

func newEnvDescribe(fs afero.Fs, envName string) (*envDescribe, error) {
	b, err := new(fs)
	if err != nil {
		return nil, err
	}

	ed := &envDescribe{
		envName: envName,
		base:    b,
	}

	return ed, nil
}

func (ed *envDescribe) Run() error {
	env, err := ed.app.Environment(ed.envName)
	if err != nil {
		return err
	}

	env.Name = ed.envName

	return yaml.NewEncoder(os.Stdout).Encode(env)
}
