package action

import (
	"os"

	"github.com/bryanl/woowoo/component"
	"github.com/bryanl/woowoo/ksutil"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// Show shows an environment.
func Show(fs afero.Fs, env string, opts ...ShowOpt) error {
	s, err := newShow(fs, env, opts...)
	if err != nil {
		return err
	}

	return s.Run()
}

// ShowOpt is an option for configuring Show.
type ShowOpt func(*show)

// ShowWithComponents selects the components to be show.
func ShowWithComponents(names ...string) ShowOpt {
	return func(s *show) {
		s.components = names
	}
}

// Show is a show Action
type show struct {
	env        string
	components []string

	*base
}

// NewShow creates an instance of Show.
func newShow(fs afero.Fs, env string, opts ...ShowOpt) (*show, error) {
	b, err := new(fs)
	if err != nil {
		return nil, err
	}

	s := &show{
		env:  env,
		base: b,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s, nil
}

// Run runs the action.
func (s *show) Run() error {
	namespaces, err := component.NamespacesFromEnv(s.app, s.env)
	if err != nil {
		return errors.Wrap(err, "find namespaces")
	}

	var objects []*unstructured.Unstructured
	for _, ns := range namespaces {
		members, err := ns.Components()
		if err != nil {
			return errors.Wrap(err, "find components")
		}

		for _, c := range members {
			if !s.isAvailable(c.Name()) {
				continue
			}

			o, err := c.Objects()
			if err != nil {
				return errors.Wrap(err, "get objects")
			}
			objects = append(objects, o...)
		}
	}

	if err := ksutil.Fprint(os.Stdout, objects, "yaml"); err != nil {
		return errors.Wrap(err, "print YAML")
	}

	return nil
}

func (s *show) isAvailable(name string) bool {
	if len(s.components) == 0 {
		return true
	}

	return stringInSlice(name, s.components)
}

func stringInSlice(s string, sl []string) bool {
	for i := range sl {
		if sl[i] == s {
			return true
		}
	}

	return false
}
