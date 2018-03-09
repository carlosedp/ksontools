package component

import (
	"path/filepath"
	"strings"

	"github.com/ksonnet/ksonnet/metadata/app"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// ParamOptions is options for parameters.
type ParamOptions struct {
	Index int
}

// Component is a ksonnet Component interface.
type Component interface {
	Name() string
	Objects() ([]*unstructured.Unstructured, error)
	SetParam(path []string, value interface{}, options ParamOptions) error
	DeleteParam(path []string, options ParamOptions) error
	Params() ([]NamespaceParameter, error)
}

const (
	// componentsDir is the name of the directory which houses components.
	componentsRoot = "components"
	// paramsFile is the params file for a component namespace.
	paramsFile = "params.libsonnet"
)

// Path returns returns the file system path for a component.
func Path(fs afero.Fs, root, name string) (string, error) {
	ns, localName := ExtractNamespacedComponent(fs, root, name)

	fis, err := afero.ReadDir(fs, ns.Dir())
	if err != nil {
		return "", err
	}

	var fileName string
	files := make(map[string]bool)

	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}

		base := strings.TrimSuffix(fi.Name(), filepath.Ext(fi.Name()))
		if _, ok := files[base]; ok {
			return "", errors.Errorf("Found multiple component files with component name %q", name)
		}
		files[base] = true

		if base == localName {
			fileName = fi.Name()
		}
	}

	if fileName == "" {
		return "", errors.Errorf("No component name %q found", name)
	}

	return filepath.Join(ns.Dir(), fileName), nil
}

// ExtractComponent extracts a component from a path.
func ExtractComponent(fs afero.Fs, root, path string) (Component, error) {
	ns, componentName := ExtractNamespacedComponent(fs, root, path)
	members, err := ns.Components()
	if err != nil {
		return nil, err
	}

	for _, member := range members {
		if componentName == member.Name() {
			return member, nil
		}
	}

	return nil, errors.Errorf("unable to find component %q", componentName)
}

func isComponentDir(fs afero.Fs, path string) (bool, error) {
	files, err := afero.ReadDir(fs, path)
	if err != nil {
		return false, errors.Wrapf(err, "read files in %s", path)
	}

	for _, file := range files {
		if file.Name() == paramsFile {
			return true, nil
		}
	}

	return false, nil
}

// AppSpecer is implemented by any value that has a AppSpec method. The AppSpec method is
// used to retrieve a ksonnet AppSpec.
type AppSpecer interface {
	AppSpec() (*app.Spec, error)
}

// MakePathsByNamespace creates a map of component paths categorized by namespace.
func MakePathsByNamespace(fs afero.Fs, appSpecer AppSpecer, root, env string) (map[Namespace][]string, error) {
	paths, err := MakePaths(fs, appSpecer, root, env)
	if err != nil {
		return nil, err
	}

	m := make(map[Namespace][]string)

	for i := range paths {
		prefix := root + "/components/"
		if strings.HasSuffix(root, "/") {
			prefix = root + "components/"
		}
		path := strings.TrimPrefix(paths[i], prefix)
		ns, _ := ExtractNamespacedComponent(fs, root, path)
		if _, ok := m[ns]; !ok {
			m[ns] = make([]string, 0)
		}

		m[ns] = append(m[ns], paths[i])
	}

	return m, nil
}

// MakePaths creates a slice of component paths
func MakePaths(fs afero.Fs, appSpecer AppSpecer, root, env string) ([]string, error) {
	cpl, err := newComponentPathLocator(fs, appSpecer, env)
	if err != nil {
		return nil, errors.Wrap(err, "create component path locator")
	}

	return cpl.Locate(root)
}
