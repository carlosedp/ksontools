package yaml2jsonnet

import (
	"fmt"
	"strings"

	"github.com/bryanl/woowoo/node"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/astext"
	"github.com/pkg/errors"
)

var (
	// ErrNotFound is a not found error.
	ErrNotFound = errors.New("not found")
)

// Node represents a node by name.
type Node struct {
	name    string
	obj     *astext.Object
	IsMixin bool
}

// NewNode creates an instance of Node.
func NewNode(name string, obj *astext.Object) *Node {
	return &Node{
		name: name,
		obj:  obj,
	}
}

// SearchResult is the results from a Search.
type SearchResult struct {
	Fields    []string
	Functions []string
	Types     []string

	MatchedPath []string

	Setter string
	Value  interface{}
}

// Search2 searches for a path in the node.
func (n *Node) Search2(path ...string) (*Item, error) {
	sp := searchPath{path: path}
	item, _, err := n.searchNode(n.obj, sp, make([]string, 0))
	return item, err
}

type searchPath struct {
	path []string
}

func (sp *searchPath) len() int {
	return len(sp.path)
}

func (sp *searchPath) isEmpty() bool {
	return sp.len() == 0
}

func (sp *searchPath) head() string {
	return sp.path[0]
}

func (sp *searchPath) tail() string {
	return sp.path[len(sp.path)-1]
}

func (sp *searchPath) descendant() searchPath {
	return searchPath{path: sp.path[1:]}
}

func (sp *searchPath) String() string {
	return strings.Join(sp.path, ".")
}

func (n *Node) searchNode(obj *astext.Object, sp searchPath, breadcrumbs []string) (*Item, []string, error) {
	if sp.isEmpty() {
		return nil, nil, errors.New("search path is empty")
	}

	members, err := node.FindMembers(obj)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable list object members")
	}

	if sp.len() == 1 {
		switch {
		case stringInSlice(sp.head(), members.Fields):
			path := append(breadcrumbs, sp.head())
			return &Item{Type: ItemTypeObject, Path: path}, nil, nil
		case stringInSlice("mixin", members.Fields):
			return n.findChild(obj, sp, "mixin", breadcrumbs)
		default:
			fnName, err := members.FindFunction(sp.head())
			if err != nil {
				return nil, nil, errors.Wrapf(err, "unable to find function %s", sp)
			}

			path := append(breadcrumbs, sp.head())
			name := fmt.Sprintf("%s.%s", strings.Join(breadcrumbs, "."), fnName)
			return &Item{Type: ItemTypeSetter, Name: name, Path: path}, nil, nil
		}
	}

	switch {
	case stringInSlice(sp.head(), members.Fields):
		return n.findChild(obj, sp.descendant(), sp.head(), breadcrumbs)
	case stringInSlice("mixin", members.Fields):
		return n.findChild(obj, sp, "mixin", breadcrumbs)
	}

	return nil, nil, errChildNotFound
}

var (
	errChildNotFound = errors.New("child not found")
)

func (n *Node) findChild(obj *astext.Object, sp searchPath, name string, breadcrumbs []string) (*Item, []string, error) {
	childBreadcrumbs := append(breadcrumbs, name)
	child, err := node.Find(obj, name)
	if err != nil {
		return nil, nil, err
	}

	item, path, err := n.searchNode(child, sp, childBreadcrumbs)
	if err != nil {
		if err == errChildNotFound {
			newSp := searchPath{path: append(breadcrumbs, name, sp.head())}
			return n.searchNode(n.obj, newSp, make([]string, 0))
		}

		return nil, nil, err
	}

	return item, path, nil
}

// Search searches for nodes given a path.
func (n *Node) Search(path ...string) (SearchResult, error) {
	return searchObj(n.obj, path...)
}

func searchObj(obj *astext.Object, path ...string) (SearchResult, error) {
	om, err := node.FindMembers(obj)
	if err != nil {
		return SearchResult{}, err
	}

	if len(path) == 0 {
		return SearchResult{
			Fields:    om.Fields,
			Functions: om.Functions,
			Types:     om.Types,
		}, nil
	}

	cur, err := node.Find(obj, path[0])
	if err != nil {
		path = append([]string{"mixin"}, path...)
		cur, err = node.Find(obj, path[0])
		if err != nil {
			// is there a function which matches this?
			fn, ferr := om.FindFunction(path[1])
			if ferr != nil {
				return SearchResult{}, errors.Errorf("node not found in path %s", strings.Join(path, "."))
			}

			return SearchResult{
				Setter:      fn,
				MatchedPath: []string{path[1]},
			}, nil
		}
	}

	sr, err := searchObj(cur, path[1:]...)
	if err != nil {
		return SearchResult{}, err
	}

	sr.MatchedPath = append([]string{path[0]}, sr.MatchedPath...)

	return sr, nil
}

var (
	ignoredProps = []string{"mixin", "kind", "new", "mixinInstance"}
)

func stringInSlice(s string, sl []string) bool {
	for i := range sl {
		if sl[i] == s {
			return true
		}
	}

	return false
}

// ItemType is the type of item.
type ItemType int

const (
	// ItemTypeSetter is a item that is a setter function.
	ItemTypeSetter ItemType = iota
	// ItemTypeObject is a item that is an object.
	ItemTypeObject
)

// Item identifies an item in a Node.
type Item struct {
	Type ItemType
	Name string
	Path []string
}
