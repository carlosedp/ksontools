package ksutil

import (
	"github.com/ksonnet/ksonnet/metadata/app"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// SuperApp is a super set of ksonnet's app. One day it'll be integrated
// back into ksonnet.
type SuperApp interface {
	app.App

	Fs() afero.Fs
	Root() string
	UpdateTargets(envName string, targets []string) error
}

// LoadApp is a wrapper around app.
func LoadApp(fs afero.Fs, root string) (SuperApp, error) {
	a, err := app.Load(fs, root)
	if err != nil {
		return nil, err
	}
	sa := &superApp{
		App:  a,
		fs:   fs,
		root: root,
	}

	return sa, nil
}

type superApp struct {
	app.App

	fs   afero.Fs
	root string
}

var _ app.App = (*superApp)(nil)

func (s *superApp) Fs() afero.Fs {
	return s.fs
}

func (s *superApp) Root() string {
	return s.root
}

// TODO: when adding this to ksonnet, App001 should return an error as it can't
// update targets
func (s *superApp) UpdateTargets(envName string, targets []string) error {
	spec, err := s.Environment(envName)
	if err != nil {
		return err
	}

	spec.Targets = targets

	return errors.Wrap(s.AddEnvironment(envName, "", spec), "update targets")
}
