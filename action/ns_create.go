package action

import (
	"path/filepath"
	"strings"

	"github.com/bryanl/woowoo/component"
	kscomponent "github.com/ksonnet/ksonnet/component"
	"github.com/ksonnet/ksonnet/metadata/app"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// NsCreate creates a component namespace
func NsCreate(fs afero.Fs, nsName string) error {
	nc, err := newNsCreate(fs, nsName)
	if err != nil {
		return err
	}

	return nc.Run()
}

type nsCreate struct {
	nsName string

	*base
}

func newNsCreate(fs afero.Fs, nsName string) (*nsCreate, error) {
	b, err := new(fs)
	if err != nil {
		return nil, err
	}

	et := &nsCreate{
		nsName: nsName,
		base:   b,
	}

	return et, nil
}

func (nc *nsCreate) Run() error {
	_, err := component.GetNamespace(nc.app, nc.nsName)
	if err == nil {
		return errors.Errorf("namespace %q already exists", nc.nsName)
	}

	// TODO: does this belong in the component namespace? (it does)
	parts := strings.Split(nc.nsName, "/")
	dir := filepath.Join(append([]string{nc.app.Root(), "components"}, parts...)...)

	if err := nc.app.Fs().MkdirAll(dir, app.DefaultFolderPermissions); err != nil {
		return err
	}

	paramsDir := filepath.Join(dir, "params.libsonnet")
	return afero.WriteFile(nc.app.Fs(), paramsDir, kscomponent.GenParamsContent(), app.DefaultFilePermissions)
}
