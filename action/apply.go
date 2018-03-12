package action

import (
	"github.com/bryanl/woowoo/env"
	"github.com/bryanl/woowoo/pkg/client"
	"github.com/spf13/afero"
)

type ApplyOptions struct {
	Create bool
	SkipGc bool
	GcTag  string
	DryRun bool
	Client *client.Config
}

// Apply applies an environment.
func Apply(fs afero.Fs, env string, options ApplyOptions) error {
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
	options    ApplyOptions

	*base
}

// NewApply creates an instance of Apply.
func newApply(fs afero.Fs, env string, options ApplyOptions) (*apply, error) {
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
	options := env.ApplyOptions{
		Create: s.options.Create,
		SkipGc: s.options.SkipGc,
		GcTag:  s.options.GcTag,
		DryRun: s.options.DryRun,
		Client: s.options.Client,
	}

	return env.Apply(s.app, s.env, s.components, options)
}
