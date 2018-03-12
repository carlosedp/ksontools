package action

import (
	"github.com/bryanl/woowoo/env"
	"github.com/bryanl/woowoo/pkg/client"
	"github.com/spf13/afero"
)

// Apply applies an environment.
func Apply(fs afero.Fs, env string, options client.ApplyOptions) error {
	s, err := newApply(fs, env, options)
	if err != nil {
		return err
	}

	return s.Run()
}

// Apply is a apply Action
type apply struct {
	env        string
	components []string
	options    client.ApplyOptions

	*base
}

// NewApply creates an instance of Apply.
func newApply(fs afero.Fs, env string, options client.ApplyOptions) (*apply, error) {
	b, err := new(fs)
	if err != nil {
		return nil, err
	}

	s := &apply{
		env:     env,
		options: options,
		base:    b,
	}

	return s, nil
}

// Run runs the action.
func (s *apply) Run() error {
	return env.Apply(s.app, s.env, s.components, s.options)
}
