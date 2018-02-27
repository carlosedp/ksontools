package jsonnetutil

import (
	"strings"

	"github.com/google/go-jsonnet/ast"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/astext"
	nm "github.com/ksonnet/ksonnet-lib/ksonnet-gen/nodemaker"
	"github.com/pkg/errors"
)

// AddObject adds a node to an object graph.
func AddObject(object *astext.Object, path []string, addition ast.Node) error {
	if len(path) == 0 {
		return errors.New("search path was empty")
	}

	pathLen := len(path) - 1
	item := path[pathLen]
	search := path[:pathLen]

	root, err := FindObject(object, search)
	if err != nil {
		return errors.Wrapf(err, "find object path %s", strings.Join(search, "."))
	}

	var parent *astext.Object

	for i := range root.Fields {
		parentID := ast.Identifier(search[len(search)-1])
		if *root.Fields[i].Id == parentID {
			parent = root.Fields[i].Expr2.(*astext.Object)
		}
	}

	if parent == nil {
		return errors.New("unable to add object")
	}

	f := astext.ObjectField{
		ObjectField: ast.ObjectField{
			Kind:  ast.ObjectFieldStr,
			Hide:  ast.ObjectFieldInherit,
			Expr1: nm.NewStringDouble(item).Node(),
			Expr2: addition,
		},
	}

	parent.Fields = append(parent.Fields, f)

	return nil
}

// UpdateObject updates a location in an object with a new node.
// Returns error if it can't find the path.
func UpdateObject(object *astext.Object, path []string, update ast.Node) error {
	if len(path) == 0 {
		return errors.New("search path was empty")
	}

	item := path[len(path)-1]

	child, err := FindObject(object, path)
	if err != nil {
		return errors.New("object not found")
	}

	for i := range child.Fields {
		id, err := FieldID(child.Fields[i])
		if err != nil {
			return err
		}

		if id == item {
			child.Fields[i].Expr2 = update
			return nil
		}
	}

	return errors.Errorf("unable to find field %s", item)
}

// FindObject finds a path in an object.
func FindObject(object *astext.Object, path []string) (*astext.Object, error) {
	if len(path) == 0 {
		return nil, errors.New("search path was empty")
	}

	for i := range object.Fields {
		id, err := FieldID(object.Fields[i])
		if err != nil {
			return nil, err
		}

		if path[0] == id {
			if len(path) == 1 {
				return object, nil
			}

			child, ok := object.Fields[i].Expr2.(*astext.Object)
			if !ok {
				return nil, errors.Errorf("child is a %T. expected an object", object.Fields[i].Expr2)
			}

			return FindObject(child, path[1:])
		}
	}

	return nil, errors.New("path was not found")
}

// FieldID returns the id for an object field.
func FieldID(field astext.ObjectField) (string, error) {
	if field.Expr1 != nil {
		lf, ok := field.Expr1.(*ast.LiteralString)
		if !ok {
			return "", errors.New("field Expr1 is not a string")
		}

		return lf.Value, nil
	}

	return string(*field.Id), nil
}
