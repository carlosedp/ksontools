package component

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bryanl/woowoo/ksutil"
	"github.com/bryanl/woowoo/params"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

func nsErrorMsg(format, nsName string) string {
	s := fmt.Sprintf("namespace %q", nsName)
	if nsName == "" {
		s = "root namespace"
	}

	return fmt.Sprintf(format, s)
}

// Namespace is a component namespace.
type Namespace struct {
	path string

	app ksutil.SuperApp
}

// ExtractNamespacedComponent extracts a namespace and a component from a path.
func ExtractNamespacedComponent(app ksutil.SuperApp, path string) (Namespace, string) {
	nsPath, component := filepath.Split(path)
	ns := Namespace{path: nsPath, app: app}
	return ns, component
}

// Name returns the namespace name.
func (n *Namespace) Name() string {
	if n.path == "" {
		return "/"
	}
	return n.path
}

// GetNamespace gets a namespace by path.
func GetNamespace(app ksutil.SuperApp, nsName string) (Namespace, error) {
	parts := strings.Split(nsName, "/")
	nsDir := filepath.Join(append([]string{app.Root(), componentsRoot}, parts...)...)

	exists, err := afero.Exists(app.Fs(), nsDir)
	if err != nil {
		return Namespace{}, err
	}

	if !exists {
		return Namespace{}, errors.New(nsErrorMsg("unable to find %s", nsName))
	}

	return Namespace{path: nsName, app: app}, nil
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
	return afero.WriteFile(n.app.Fs(), n.ParamsPath(), []byte(src), 0644)
}

// Dir is the absolute directory for a namespace.
func (n *Namespace) Dir() string {
	parts := strings.Split(n.path, "/")
	path := []string{n.app.Root(), componentsRoot}
	if len(n.path) != 0 {
		path = append(path, parts...)
	}

	return filepath.Join(path...)
}

// NamespaceParameter is a namespaced paramater.
type NamespaceParameter struct {
	Component string
	Index     string
	Key       string
	Value     string
}

// ResolvedParams resolves paramsters for a namespace. It returns a JSON encoded
// string of component parameters.
func (n *Namespace) ResolvedParams() (string, error) {
	s, err := n.readParams()
	if err != nil {
		return "", err
	}

	return applyGlobals(s)
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
	b, err := afero.ReadFile(n.app.Fs(), n.ParamsPath())
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// NamespacesFromEnv returns all namespaces given an environment.
func NamespacesFromEnv(app ksutil.SuperApp, env string) ([]Namespace, error) {
	paths, err := MakePaths(app, env)
	if err != nil {
		return nil, err
	}

	prefix := app.Root() + "/components"

	seen := make(map[string]bool)
	var namespaces []Namespace
	for _, path := range paths {
		nsName := strings.TrimPrefix(path, prefix)
		if _, ok := seen[nsName]; !ok {
			seen[nsName] = true
			ns, err := GetNamespace(app, nsName)
			if err != nil {
				return nil, err
			}

			namespaces = append(namespaces, ns)
		}
	}

	return namespaces, nil
}

// Namespaces returns all component namespaces
func Namespaces(app ksutil.SuperApp) ([]Namespace, error) {
	componentRoot := filepath.Join(app.Root(), componentsRoot)

	var namespaces []Namespace

	err := afero.Walk(app.Fs(), componentRoot, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.IsDir() {
			ok, err := isComponentDir(app.Fs(), path)
			if err != nil {
				return err
			}

			if ok {
				nsPath := strings.TrimPrefix(path, componentRoot)
				nsPath = strings.TrimPrefix(nsPath, string(filepath.Separator))
				ns := Namespace{path: nsPath, app: app}
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
	nsDir := filepath.Join(append([]string{n.app.Root(), componentsRoot}, parts...)...)

	fis, err := afero.ReadDir(n.app.Fs(), nsDir)
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
			component := NewYAML(n.app, path, n.ParamsPath())
			components = append(components, component)
		}
	}

	return components, nil
}
