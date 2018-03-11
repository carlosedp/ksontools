package action

import (
	"github.com/bryanl/woowoo/env"
	"github.com/spf13/afero"
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
	return env.Show(s.app, s.env, s.components)
}
