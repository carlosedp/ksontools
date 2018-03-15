package component

import (
	"path/filepath"
	"strings"

	"github.com/ksonnet/ksonnet/metadata/app"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type libPather interface {
	LibPath(envName string) (string, error)
}

// ParamOptions is options for parameters.
type ParamOptions struct {
	Index int
}

// Summary summarizes items found in components.
type Summary struct {
	ComponentName string
	IndexStr      string
	Index         int
	Type          string
	APIVersion    string
	Kind          string
	Name          string
}

// GVK converts a summary to a group - version - kind.
func (s *Summary) typeSpec() (*TypeSpec, error) {
	return NewTypeSpec(s.APIVersion, s.Kind)
}

// Component is a ksonnet Component interface.
type Component interface {
	// Name is the component name.
	Name(wantsNamedSpaced bool) string
	// Objects converts the component to a set of objects.
	Objects(paramsStr, envName string) ([]*unstructured.Unstructured, error)
	// SetParams sets a component paramaters.
	SetParam(path []string, value interface{}, options ParamOptions) error
	// DeleteParam deletes a component parameter.
	DeleteParam(path []string, options ParamOptions) error
	// Params returns a list of all parameters for a component.
	Params() ([]NamespaceParameter, error)
	// Summarize returns a summary of the component.
	Summarize() ([]Summary, error)
}

const (
	// componentsDir is the name of the directory which houses components.
	componentsRoot = "components"
	// paramsFile is the params file for a component namespace.
	paramsFile = "params.libsonnet"
)

// Path returns returns the file system path for a component.
func Path(a app.App, name string) (string, error) {
	ns, localName := ExtractNamespacedComponent(a, name)

	fis, err := afero.ReadDir(a.Fs(), ns.Dir())
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
func ExtractComponent(a app.App, path string) (Component, error) {
	ns, componentName := ExtractNamespacedComponent(a, path)
	members, err := ns.Components()
	if err != nil {
		return nil, err
	}

	for _, member := range members {
		if componentName == member.Name(false) {
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

// MakePathsByNamespace creates a map of component paths categorized by namespace.
func MakePathsByNamespace(a app.App, env string) (map[Namespace][]string, error) {
	paths, err := MakePaths(a, env)
	if err != nil {
		return nil, err
	}

	m := make(map[Namespace][]string)

	for i := range paths {
		prefix := a.Root() + "/components/"
		if strings.HasSuffix(a.Root(), "/") {
			prefix = a.Root() + "components/"
		}
		path := strings.TrimPrefix(paths[i], prefix)
		ns, _ := ExtractNamespacedComponent(a, path)
		if _, ok := m[ns]; !ok {
			m[ns] = make([]string, 0)
		}

		m[ns] = append(m[ns], paths[i])
	}

	return m, nil
}

// MakePaths creates a slice of component paths
func MakePaths(a app.App, env string) ([]string, error) {
	cpl, err := newComponentPathLocator(a, env)
	if err != nil {
		return nil, errors.Wrap(err, "create component path locator")
	}

	return cpl.Locate()
}
