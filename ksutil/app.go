package ksutil

import (
	"github.com/ksonnet/ksonnet/metadata/app"
	"github.com/spf13/afero"
)

// SuperApp is a super set of ksonnet's app. One day it'll be integrated
// back into ksonnet.
type SuperApp interface {
	app.App

	Fs() afero.Fs
	Root() string
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
