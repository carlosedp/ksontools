package node

import (
	"github.com/google/go-jsonnet/ast"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/astext"
	"github.com/pkg/errors"
)

// Find finds a node by name in a parent node.
func Find(node ast.Node, name string) (*astext.Object, error) {
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
