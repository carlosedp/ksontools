package component

import (
	"path/filepath"
	"sort"

	"github.com/bryanl/woowoo/ksutil"

	"github.com/ksonnet/ksonnet/metadata/app"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type componentPathLocator struct {
	app     ksutil.SuperApp
	envSpec *app.EnvironmentSpec
}

func newComponentPathLocator(app ksutil.SuperApp, envName string) (*componentPathLocator, error) {
	if app == nil {
		return nil, errors.New("app is nil")
	}

	env, err := app.Environment(envName)
	if err != nil {
		return nil, err
	}

	return &componentPathLocator{
		app:     app,
		envSpec: env,
	}, nil
}

func (cpl *componentPathLocator) Locate() ([]string, error) {
	targets := cpl.envSpec.Targets
	rootPath := cpl.app.Root()

	if len(targets) == 0 {
		return []string{filepath.Join(rootPath, componentsRoot)}, nil
	}

	var paths []string

	for _, target := range targets {
		childPath := filepath.Join(rootPath, componentsRoot, target)
		exists, err := afero.DirExists(cpl.app.Fs(), childPath)
		if err != nil {
			return nil, err
		}

		if !exists {
			return nil, errors.Errorf("target %q is not valid", target)
		}

		paths = append(paths, childPath)
	}

	sort.Strings(paths)

	return paths, nil
}

// isComponent reports if a file is a component. Components have a `jsonnet` extension.
func isComponent(path string) bool {
	for _, s := range []string{".jsonnet", ".yaml", "json"} {
		if s == filepath.Ext(path) {
			return true
		}
	}
	return false
}
