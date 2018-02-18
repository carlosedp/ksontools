package yaml2jsonnet

import (
	"fmt"
	"sort"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
)

var (
	// TODO: need something in ksonnet lib to look this up
	groupLookup = map[string][]string{
		"apiextensions.k8s.io":      []string{"hidden", "apiextensions"},
		"rbac.authorization.k8s.io": []string{"hidden", "rbac"},
	}
)

// PropertyPath contains a property path.
type PropertyPath struct {
	Path  []string
	Value interface{}
}

// Properties are document properties
type Properties map[interface{}]interface{}

// Paths returns a list of paths in properties.
func (p Properties) Paths(gvk GVK) []PropertyPath {
	ch := make(chan PropertyPath)

	go func() {
		g, ok := groupLookup[gvk.Group[0]]
		if !ok {
			g = gvk.Group
		}

		base := append(g, gvk.Version, gvk.Kind)
		iterateMap(ch, base, p)
		close(ch)
	}()

	var out []PropertyPath
	for pr := range ch {
		out = append(out, pr)
	}

	return out
}

func iterateMap(ch chan PropertyPath, base []string, m map[interface{}]interface{}) {
	localBase := make([]string, len(base))
	copy(localBase, base)

	var keys []interface{}
	for k := range m {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		a := keys[i].(string)
		b := keys[j].(string)

		return a < b
	})

	for i := range keys {
		name := keys[i].(string)
		switch t := m[name].(type) {
		default:
			panic(fmt.Sprintf("not sure what to do with %T", t))
		case map[interface{}]interface{}:
			newBase := append(localBase, name)
			iterateMap(ch, newBase, t)
		case string, int, []interface{}:
			ch <- PropertyPath{
				Path: append(base, name),
			}
		}
	}
}

// Value returns the value at a path.
func (p Properties) Value(path []string) (interface{}, error) {
	return valueSearch(path, p)
}

func valueSearch(path []string, m map[interface{}]interface{}) (interface{}, error) {
	if len(path) > 0 && path[0] == "mixin" {
		return valueSearch(path[1:], m)
	}

	var keys []interface{}
	for k := range m {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		a := keys[i].(string)
		b := keys[j].(string)

		return a < b
	})

	for i := range keys {
		name := keys[i].(string)
		if name == path[0] {

			switch t := m[name].(type) {
			default:
				panic(fmt.Sprintf("not sure what to do with %T", t))
			case map[interface{}]interface{}:
				if len(path) == 1 {
					return t, nil
				}
				return valueSearch(path[1:], t)
			case string, int, []interface{}:
				return t, nil
			}
		}

	}

	spew.Dump(m)

	return nil, errors.Errorf("unable to find %s", strings.Join(path, "."))
}
