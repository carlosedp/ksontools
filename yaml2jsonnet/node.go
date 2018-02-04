package yaml2jsonnet

import (
	"fmt"
	"strings"

	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/ast"
	"github.com/pkg/errors"
)

var (
	// ErrNotFound is a not found error.
	ErrNotFound = errors.New("not found")
)

func FindType(gvk GVK, node ast.Node) (*ast.Object, error) {
	group, err := FindNode(node, gvk.Group)
	if err != nil {
		return nil, errors.Wrap(err, "find group")
	}

	version, err := FindNode(group, gvk.Version)
	if err != nil {
		return nil, errors.Wrap(err, "find version")
	}

	kind, err := FindNode(version, gvk.Kind)
	if err != nil {
		return nil, errors.Wrap(err, "find kind")
	}

	return kind, nil
}

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

func (n *Node) Property(name string) (*Node, error) {
	mixins, err := n.Mixins()
	if err != nil {
		return nil, err
	}

	for _, mixin := range mixins {
		if name == mixin.name {
			fmt.Println(name, "is a mixin")

			_, err := mixin.Properties()
			if err != nil {
				return nil, errors.Wrap(err, "looking inside mixin")
			}

			return &mixin, nil
		}
	}

	return nil, errors.Errorf("property %s not found", name)
}

func (n *Node) Mixins() ([]Node, error) {
	if n.obj == nil {
		return nil, errors.New("object is nil")
	}

	mixinObj, err := FindNode(n.obj, "mixin")
	if err != nil {
		return nil, errors.New("object does not have a mixin field")
	}

	var mixins []Node
	for _, of := range mixinObj.Fields {
		if of.Id == nil {
			return nil, errors.New("mixin field has nil identifier")
		}

		id := string(*of.Id)

		// mixins are objects.
		obj, ok := of.Expr2.(*ast.Object)
		if ok {
			n := NewNode(id, obj)
			n.IsMixin = true
			mixins = append(mixins, *n)
		}
	}

	return mixins, nil
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
	} else {
		return fmt.Sprintf("with%s", name2), nil
	}

	return "", nil
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
