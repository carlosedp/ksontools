package action

import (
	"os"

	"github.com/bryanl/woowoo/component"
	"github.com/bryanl/woowoo/ksutil"
	"github.com/spf13/afero"
)

// NsList lists available namespaces
func NsList(fs afero.Fs) error {
	nl, err := newNsList(fs)
	if err != nil {
		return err
	}

	return nl.Run()
}

type nsList struct {
	*base
}

func newNsList(fs afero.Fs) (*nsList, error) {
	b, err := new(fs)
	if err != nil {
		return nil, err
	}

	nl := &nsList{
		base: b,
	}

	return nl, nil
}

func (nl *nsList) Run() error {
	namespaces, err := component.Namespaces(nl.app)
	if err != nil {
		return err
	}

	table := ksutil.NewTable(os.Stdout)
	table.SetHeader([]string{"namespace"})

	for _, ns := range namespaces {
		table.Append([]string{ns.Name()})
	}

	table.Render()
	return nil
}
