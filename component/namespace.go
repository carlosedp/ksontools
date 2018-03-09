package component

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bryanl/woowoo/params"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// Namespace is a component namespace.
type Namespace struct {
	path string

	root string
	fs   afero.Fs
}

// ExtractNamespacedComponent extracts a namespace and a component from a path.
func ExtractNamespacedComponent(fs afero.Fs, root, path string) (Namespace, string) {
	nsPath, component := filepath.Split(path)
	ns := Namespace{path: nsPath, root: root, fs: fs}
	return ns, component
}

// Name returns the namespace name.
func (n *Namespace) Name() string {
	return n.path
}

// GetNamespace gets a namespace by path.
func GetNamespace(fs afero.Fs, root, nsName string) (Namespace, error) {
	parts := strings.Split(nsName, "/")
	nsDir := filepath.Join(append([]string{root, componentsRoot}, parts...)...)

	exists, err := afero.Exists(fs, nsDir)
	if err != nil {
		return Namespace{}, err
	}

	if !exists {
		return Namespace{}, errors.Errorf("unable to find namespace %q", nsName)
	}

	return Namespace{path: nsName, root: root, fs: fs}, nil
}

// ParamsPath generates the path to params.libsonnet for a namespace.
func (n *Namespace) ParamsPath() string {
	return filepath.Join(n.Dir(), paramsFile)
}

// SetParam sets params for a namespace.
func (n *Namespace) SetParam(path []string, value interface{}) error {
	paramsData, err := n.readParams()
	if err != nil {
		return err
	}

	props, err := params.ToMap("", paramsData, "global")
	if err != nil {
		return err
	}

	// TODO: this is duplicated in YAML.SetParam
	changes := make(map[string]interface{})
	cur := changes

	for i, k := range path {
		if i == len(path)-1 {
			cur[k] = value
		} else {
			if _, ok := cur[k]; !ok {
				m := make(map[string]interface{})
				cur[k] = m
				cur = m
			}
		}
	}

	if err = mergeMaps(props, changes, nil); err != nil {
		return err
	}

	updatedParams, err := params.Update([]string{"global"}, paramsData, changes)
	if err != nil {
		return err
	}

	if err = n.writeParams(updatedParams); err != nil {
		return err
	}

	return nil
}

func (n *Namespace) writeParams(src string) error {
	return afero.WriteFile(n.fs, n.ParamsPath(), []byte(src), 0644)
}

// Dir is the absolute directory for a namespace.
func (n *Namespace) Dir() string {
	parts := strings.Split(n.path, "/")
	path := []string{n.root, componentsRoot}
	if len(n.path) != 0 {
		path = append(path, parts...)
	}

	return filepath.Join(path...)
}

// NamespaceParameter is a namespaced paramater.
type NamespaceParameter struct {
	Component string
	Key       string
	Value     string
}

// Params returns the params for a namespace.
func (n *Namespace) Params() ([]NamespaceParameter, error) {
	components, err := n.Components()
	if err != nil {
		return nil, err
	}

	var nsps []NamespaceParameter
	for _, c := range components {
		params, err := c.Params()
		if err != nil {
			return nil, err
		}

		for _, p := range params {
			nsps = append(nsps, p)
		}
	}

	return nsps, nil
}

func (n *Namespace) readParams() (string, error) {
	b, err := afero.ReadFile(n.fs, n.ParamsPath())
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// NamespacesFromEnv returns all namespaces given an environment.
func NamespacesFromEnv(fs afero.Fs, appSpecer AppSpecer, root, env string) ([]Namespace, error) {
	paths, err := MakePaths(fs, appSpecer, root, env)
	if err != nil {
		return nil, err
	}

	var namespaces []Namespace
	seen := make(map[string]bool)
	for i := range paths {
		prefix := root + "/components/"
		if strings.HasSuffix(root, "/") {
			prefix = root + "components/"
		}

		path := strings.TrimPrefix(paths[i], prefix)
		ns, _ := ExtractNamespacedComponent(fs, root, path)
		if _, ok := seen[ns.Name()]; ok {
			continue
		}
		seen[ns.Name()] = true
		namespaces = append(namespaces, ns)
	}

	return namespaces, nil
}

// Namespaces returns all component namespaces
func Namespaces(fs afero.Fs, root string) ([]Namespace, error) {
	componentRoot := filepath.Join(root, componentsRoot)

	var namespaces []Namespace

	err := afero.Walk(fs, componentRoot, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.IsDir() {
			ok, err := isComponentDir(fs, path)
			if err != nil {
				return err
			}

			if ok {
				nsPath := strings.TrimPrefix(path, componentRoot)
				nsPath = strings.TrimPrefix(nsPath, string(filepath.Separator))
				ns := Namespace{path: nsPath, fs: fs, root: root}
				namespaces = append(namespaces, ns)
			}
		}

		return nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "walk component path")
	}

	sort.Slice(namespaces, func(i, j int) bool {
		return namespaces[i].Name() < namespaces[j].Name()
	})

	return namespaces, nil
}

// Components returns the components in a namespace.
func (n *Namespace) Components() ([]Component, error) {
	parts := strings.Split(n.path, "/")
	nsDir := filepath.Join(append([]string{n.root, componentsRoot}, parts...)...)

	fis, err := afero.ReadDir(n.fs, nsDir)
	if err != nil {
		return nil, err
	}

	var components []Component
	for _, fi := range fis {

		ext := filepath.Ext(fi.Name())
		path := filepath.Join(nsDir, fi.Name())

		switch ext {
		// TODO: these should be constants
		case ".yaml":
			component := NewYAML(n.fs, path, n.ParamsPath())
			components = append(components, component)
		}
	}

	return components, nil
}
