package action

import (
	"github.com/bryanl/woowoo/k8sutil"
	"github.com/bryanl/woowoo/pipeline"
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
	p := pipeline.New(s.app, s.env)

	objects, err := p.Objects(s.components)
	if err != nil {
		return err
	}

	// TODO: create better semantics around apply
	c := k8sutil.ApplyCmd{
		Env:          s.env,
		Create:       s.options.Create,
		GcTag:        s.options.GcTag,
		SkipGc:       s.options.SkipGc,
		DryRun:       s.options.DryRun,
		ClientConfig: s.options.Client,
	}

	return c.Run(objects, "")
}
