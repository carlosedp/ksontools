package action

import (
	"os"
	"sort"

	"github.com/bryanl/woowoo/component"
	"github.com/bryanl/woowoo/ksutil"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// ComponentList create a list of components in a namespace.
func ComponentList(fs afero.Fs, namespace, output string) error {
	cl, err := newComponentList(fs, namespace, output)
	if err != nil {
		return err
	}

	return cl.run()
}

type componentList struct {
	nsName string
	output string

	*base
}

func newComponentList(fs afero.Fs, namespace, output string) (*componentList, error) {
	b, err := new(fs)
	if err != nil {
		return nil, err
	}

	cl := &componentList{
		nsName: namespace,
		output: output,
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

	switch cl.output {
	default:
		return errors.Errorf("invalid output option %q", cl.output)
	case "":
		cl.listComponents(components)
	case "wide":
		if err := cl.listComponentsWide(components); err != nil {
			return err
		}
	}

	return nil
}

func (cl *componentList) listComponents(components []component.Component) {
	var list []string
	for _, c := range components {
		list = append(list, c.Name())
	}

	sort.Strings(list)

	table := ksutil.NewTable(os.Stdout)
	table.SetHeader([]string{"component"})
	for _, item := range list {
		table.Append([]string{item})
	}
	table.Render()
}

func (cl *componentList) listComponentsWide(components []component.Component) error {
	var rows [][]string
	for _, c := range components {
		summaries, err := c.Summarize()
		if err != nil {
			return err
		}

		for _, summary := range summaries {
			row := []string{
				summary.ComponentName,
				summary.Type,
				summary.Index,
				summary.APIVersion,
				summary.Kind,
				summary.Name,
			}

			rows = append(rows, row)

		}
	}

	table := ksutil.NewTable(os.Stdout)
	table.SetHeader([]string{"component", "type", "index", "apiversion", "kind", "name"})
	table.AppendBulk(rows)
	table.Render()

	return nil
}
