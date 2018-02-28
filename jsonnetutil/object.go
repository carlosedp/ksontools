package jsonnetutil

import (
	"fmt"
	"regexp"

	"github.com/google/go-jsonnet/ast"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/astext"
	nm "github.com/ksonnet/ksonnet-lib/ksonnet-gen/nodemaker"
	"github.com/pkg/errors"
)

// Set sets an object key at path to a value.
func Set(object *astext.Object, path []string, value ast.Node) error {
	if len(path) == 0 {
		return errors.New("path was empty")
	}

	curObj := object

	for i, k := range path {
		field, err := findField(curObj, k)
		if err != nil {
			switch err.(type) {
			default:
				return err
			case *unknownField:
				field, err = createFieldWithName(k)
				if err != nil {
					return err
				}
				curObj.Fields = append(curObj.Fields, *field)
			}
		}

		if i == len(path)-1 {
			field, _ = findField(curObj, k)
			if canUpdateObject(field.Expr2, value) {
				return errors.New("can't set object to non object")
			}
			field.Expr2 = value
			return nil
		}

		if field.Expr2 == nil {
			curObj = &astext.Object{}
			field.Expr2 = curObj
		} else if obj, ok := field.Expr2.(*astext.Object); ok {
			curObj = obj
		} else {
			return errors.Errorf("child is not an object at %q", k)
		}
	}

	return nil
}

var (
	// TODO: move this to nodemaker after 0.9 release
	reFieldStr = regexp.MustCompile(`^[A-Za-z]+[A-Za-z0-9\-]*$`)
	reField    = regexp.MustCompile(`^[A-Za-z]+[A-Za-z0-9]*$`)
)

func canUpdateObject(node1, node2 ast.Node) bool {
	return isNodeObject(node1) && !isNodeObject(node2)
}

func isNodeObject(node ast.Node) bool {
	_, ok := node.(*astext.Object)
	return ok
}

func createFieldWithName(name string) (*astext.ObjectField, error) {
	of := astext.ObjectField{ObjectField: ast.ObjectField{Hide: ast.ObjectFieldInherit}}
	if reField.MatchString(name) {
		id := ast.Identifier(name)
		of.Kind = ast.ObjectFieldID
		of.Id = &id
	} else if reFieldStr.MatchString(name) {
		of.Expr1 = nm.NewStringDouble(name).Node()
		of.Kind = ast.ObjectFieldStr
	} else {
		return nil, errors.Errorf("invalid field name %q", name)
	}

	return &of, nil
}

type unknownField struct {
	name string
}

func (e *unknownField) Error() string {
	return fmt.Sprintf("unable to find field %q", e.name)
}

func findField(object *astext.Object, id string) (*astext.ObjectField, error) {
	for i := range object.Fields {
		fieldID, err := FieldID(object.Fields[i])
		if err != nil {
			return nil, err
		}

		if id == fieldID {
			return &object.Fields[i], nil
		}
	}

	return nil, &unknownField{name: id}
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

	if field.Id == nil {
		return "", errors.New("field does not have an ID")
	}

	return string(*field.Id), nil
}
