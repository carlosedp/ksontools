package action

import (
	"fmt"
	"sort"

	"github.com/bryanl/woowoo/component"
	"github.com/spf13/afero"
)

// ComponentList create a list of components in a namespace.
func ComponentList(fs afero.Fs, namespace string) error {
	cl, err := newComponentList(fs, namespace)
	if err != nil {
		return err
	}

	return cl.run()
}

type componentList struct {
	nsName string

	*base
}

func newComponentList(fs afero.Fs, namespace string) (*componentList, error) {
	b, err := new(fs)
	if err != nil {
		return nil, err
	}

	cl := &componentList{
		nsName: namespace,
		base:   b,
	}

	return cl, nil
}

func (cl *componentList) run() error {
	ns, err := component.GetNamespace(cl.fs, cl.pluginEnv.AppDir, cl.nsName)
	if err != nil {
		return err
	}

	components, err := ns.Components()
	if err != nil {
		return err
	}

	var list []string
	for _, c := range components {
		list = append(list, c.Name())
	}

	sort.Strings(list)

	fmt.Println("COMPONENT")
	fmt.Println("=========")

	for _, name := range list {
		fmt.Println(name)
	}

	return nil
}
