package component

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/ksonnet/ksonnet/metadata/app"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type componentPathLocator struct {
	fs      afero.Fs
	envSpec *app.EnvironmentSpec
}

func newComponentPathLocator(fs afero.Fs, appSpecer AppSpecer, env string) (*componentPathLocator, error) {
	if appSpecer == nil {
		return nil, errors.New("appSpecer is nil")
	}

	if fs == nil {
		return nil, errors.New("fs is nil")
	}

	appSpec, err := appSpecer.AppSpec()
	if err != nil {
		return nil, errors.Wrap(err, "lookup application spec")
	}

	envSpec, ok := appSpec.GetEnvironmentSpec(env)
	if !ok {
		return nil, errors.Errorf("can't find %s environment", env)
	}

	return &componentPathLocator{
		fs:      fs,
		envSpec: envSpec,
	}, nil
}

func (cpl *componentPathLocator) Locate(root string) ([]string, error) {
	if len(cpl.envSpec.Targets) == 0 {
		return cpl.defaultPaths(root)
	}

	var paths []string

	for _, target := range cpl.envSpec.Targets {
		childPaths, err := cpl.expandPath(root, target)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to expand %s", target)
		}
		paths = append(paths, childPaths...)
	}

	sort.Strings(paths)

	return paths, nil
}

// expandPath take a root and a target and returns all the jsonnet components in descendant paths.
func (cpl *componentPathLocator) expandPath(root, target string) ([]string, error) {
	path := filepath.Join(root, componentsRoot, target)
	fi, err := cpl.fs.Stat(path)
	if err != nil {
		return nil, err
	}

	var paths []string

	walkFn := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !fi.IsDir() && isComponent(path) {
			paths = append(paths, path)
		}

		return nil
	}

	if fi.IsDir() {
		rootPath := filepath.Join(root, componentsRoot, fi.Name())
		if err := afero.Walk(cpl.fs, rootPath, walkFn); err != nil {
			return nil, errors.Wrapf(err, "search for components in %s", fi.Name())
		}
	} else if isComponent(fi.Name()) {
		paths = append(paths, path)
	}

	return paths, nil
}

func (cpl *componentPathLocator) defaultPaths(root string) ([]string, error) {
	var paths []string

	walkFn := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !fi.IsDir() && isComponent(path) {
			paths = append(paths, path)
		}

		return nil
	}

	componentRoot := filepath.Join(root, componentsRoot)

	if err := afero.Walk(cpl.fs, componentRoot, walkFn); err != nil {
		return nil, errors.Wrap(err, "search for components")
	}

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
