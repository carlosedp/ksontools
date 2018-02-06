package yaml2jsonnet

import (
	"fmt"
	"strings"

	"github.com/google/go-jsonnet/ast"
	"github.com/pkg/errors"
)

var (
	// ErrNotFound is a not found error.
	ErrNotFound = errors.New("not found")
)

func FindNode(node ast.Node, name string) (*ast.Object, error) {
	root, ok := node.(*ast.Object)
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

			child, ok := of.Expr2.(*ast.Object)
			if !ok {
				return nil, errors.New("child was not an Object")
			}

			return child, nil
		}
	}

	return nil, errors.Errorf("could not find %s", name)
}

type Node struct {
	name    string
	obj     *ast.Object
	IsMixin bool
}

func NewNode(name string, obj *ast.Object) *Node {
	return &Node{
		name: name,
		obj:  obj,
	}
}

type SearchResult struct {
	Fields    []string
	Functions []string
	Types     []string

	Setter string
	Value  interface{}
}

func (n *Node) Search(path ...string) (SearchResult, []string, error) {
	return searchObj(n.obj, path...)
}

func searchObj(obj *ast.Object, path ...string) (SearchResult, []string, error) {
	om, err := objMembers(obj)
	if err != nil {
		return SearchResult{}, nil, err
	}

	if len(path) == 0 {
		return SearchResult{
			Fields:    om.fields,
			Functions: om.functions,
			Types:     om.types,
		}, nil, nil
	}

	cur, err := FindNode(obj, path[0])
	if err != nil {
		path = append([]string{"mixin"}, path...)
		cur, err = FindNode(obj, path[0])
		if err != nil {
			// is there a function which matches this?
			fn, err := om.findFunction(path[1])
			if err != nil {
				return SearchResult{}, nil, errors.New("node not found")
			}

			return SearchResult{Setter: fn}, []string{path[1]}, nil
		}
	}

	sr, mp, err := searchObj(cur, path[1:]...)
	if err != nil {
		return SearchResult{}, nil, err
	}

	return sr, append([]string{path[0]}, mp...), nil
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

func objMembers(obj *ast.Object) (objMember, error) {
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

		if _, ok := of.Expr2.(*ast.Object); ok && !strings.HasPrefix(id, "__") {
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

func (n *Node) Properties() ([]Property, error) {
	if n.obj == nil {
		return nil, errors.New("object is nil")
	}

	var props []Property
	for _, of := range n.obj.Fields {
		if of.Id == nil {
			return nil, errors.New("property has nil identifier")
		}

		id := string(*of.Id)

		if stringInSlice(id, ignoredProps) {
			continue
		}

		if strings.HasSuffix(id, "Mixin") {
			continue
		}

		// prop := Property{
		// 	Node: NewNode(id, of.Expr2),
		// }

		// props = append(props, prop)
	}

	return props, nil
}

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
