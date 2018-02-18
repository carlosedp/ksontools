package yaml2jsonnet

import (
	"fmt"
	"strings"

	"github.com/google/go-jsonnet/ast"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/astext"
	"github.com/pkg/errors"
)

var (
	// ErrNotFound is a not found error.
	ErrNotFound = errors.New("not found")
)

// FindNode finds a node by name in a parent node.
func FindNode(node ast.Node, name string) (*astext.Object, error) {
	root, ok := node.(*astext.Object)
	if !ok {
		return nil, errors.New("node is not an object")
	}

	for _, of := range root.Fields {
		if of.Id == nil {
			continue
		}

		id := string(*of.Id)
		if id == name {
			if of.Expr2 == nil {
				return nil, errors.New("child object was nil")
			}

			child, ok := of.Expr2.(*astext.Object)
			if !ok {
				return nil, errors.New("child was not an Object")
			}

			return child, nil
		}
	}

	return nil, errors.Errorf("could not find %s", name)
}

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

// Search searches for nodes given a path.
func (n *Node) Search(path ...string) (SearchResult, error) {
	return searchObj(n.obj, path...)
}

func searchObj(obj *astext.Object, path ...string) (SearchResult, error) {
	om, err := objMembers(obj)
	if err != nil {
		return SearchResult{}, err
	}

	if len(path) == 0 {
		return SearchResult{
			Fields:    om.fields,
			Functions: om.functions,
			Types:     om.types,
		}, nil
	}

	cur, err := FindNode(obj, path[0])
	if err != nil {
		path = append([]string{"mixin"}, path...)
		cur, err = FindNode(obj, path[0])
		if err != nil {
			// is there a function which matches this?
			fn, ferr := om.findFunction(path[1])
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

type objMember struct {
	fields    []string
	functions []string
	types     []string
}

func (om *objMember) findFunction(name string) (string, error) {
	var hasSetter, hasSetterMixin, hasType bool

	name2 := strings.Title(name)

	for _, id := range om.functions {
		if fn := fmt.Sprintf("with%s", name2); fn == id && stringInSlice(fn, om.functions) {
			hasSetter = true
		}
		if fn := fmt.Sprintf("with%sMixin", name2); fn == id && stringInSlice(fn, om.functions) {
			hasSetterMixin = true
		}
		if t := fmt.Sprintf("%sType", name); t == id && stringInSlice(t, om.types) {
			hasType = true
		}
	}

	if hasSetter && hasSetterMixin && hasType {
		return fmt.Sprintf("with%s", name2), nil
	} else if hasSetter && hasSetterMixin {
		return fmt.Sprintf("with%s", name2), nil
	} else if hasType {
		return "", errors.New("what to do with mixins")
	} else if hasSetter {
		return fmt.Sprintf("with%s", name2), nil
	}

	return "", errors.Errorf("could not find function %s", name)
}

func objMembers(obj *astext.Object) (objMember, error) {
	if obj == nil {
		return objMember{}, errors.New("object is nil")
	}

	var om objMember

	for _, of := range obj.Fields {
		if of.Id == nil {
			continue
		}

		id := string(*of.Id)

		if of.Method != nil && !strings.HasPrefix(id, "__") && !strings.HasPrefix(id, "mixin") {
			om.functions = append(om.functions, id)
			continue
		}

		if _, ok := of.Expr2.(*astext.Object); ok && !strings.HasPrefix(id, "__") {
			om.fields = append(om.fields, id)
			continue
		}

		if strings.HasSuffix(id, "Type") {
			om.types = append(om.types, id)
			continue
		}
	}

	return om, nil
}

var (
	ignoredProps = []string{"mixin", "kind", "new", "mixinInstance"}
)

func (n *Node) FindFunction(p, name string) (string, error) {
	var hasSetter, hasSetterMixin, hasType bool

	name2 := strings.Title(name)

	for _, f := range n.obj.Fields {
		id := string(*f.Id)
		if fmt.Sprintf("with%s", name2) == id {
			hasSetter = true
		}
		if fmt.Sprintf("with%sMixin", name2) == id {
			hasSetterMixin = true
		}
		if fmt.Sprintf("%sType", name) == id {
			hasType = true
		}
	}

	if hasSetter && hasSetterMixin && hasType {
		return fmt.Sprintf("with%s", name2), nil
	} else if hasSetter && hasSetterMixin {
		return fmt.Sprintf("with%s", name2), nil
	} else if hasType {
		return "", errors.New("what to do with mixins")
	}

	return fmt.Sprintf("with%s", name2), nil
}

func stringInSlice(s string, sl []string) bool {
	for i := range sl {
		if sl[i] == s {
			return true
		}
	}

	return false
}

type PropertyType int

const (
	PropertyTypeItem PropertyType = iota
	PropertyTypeObject
	PropertyTypeArray
)

type Property struct {
	Type PropertyType
	Node *Node
}
