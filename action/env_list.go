package action

import (
	"os"

	"github.com/bryanl/woowoo/ksutil"
	"github.com/spf13/afero"
)

// EnvList lists available namespaces
func EnvList(fs afero.Fs) error {
	nl, err := newEnvList(fs)
	if err != nil {
		return err
	}

	return nl.Run()
}

type envList struct {
	*base
}

func newEnvList(fs afero.Fs) (*envList, error) {
	b, err := new(fs)
	if err != nil {
		return nil, err
	}

	nl := &envList{
		base: b,
	}

	return nl, nil
}

func (nl *envList) Run() error {
	environments, err := nl.app.Environments()
	if err != nil {
		return err
	}

	table := ksutil.NewTable(os.Stdout)
	table.SetHeader([]string{"name", "kubernetes-version", "namespace", "server"})

	for name, env := range environments {
		table.Append([]string{
			name,
			env.KubernetesVersion,
			env.Destination.Namespace,
			env.Destination.Server,
		})
	}

	table.Render()
	return nil
}
