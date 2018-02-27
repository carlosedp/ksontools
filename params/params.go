package params

import (
	"bytes"

	"github.com/bryanl/woowoo/jsonnetutil"
	"github.com/google/go-jsonnet/ast"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/astext"
	nm "github.com/ksonnet/ksonnet-lib/ksonnet-gen/nodemaker"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/printer"
	"github.com/pkg/errors"
)

// Update updates a params file with the params for a component.
func Update(componentName, src string, params map[string]interface{}) (string, error) {
	obj, err := jsonnetutil.Parse("params.libsonnet", src)
	if err != nil {
		return "", errors.Wrap(err, "parse jsonnet")
	}

	paramsObject, err := nm.KVFromMap(params)
	if err != nil {
		return "", errors.Wrap(err, "convert params to object")
	}

	path := []string{"components", componentName}

	astParamsObject := paramsObject.Node().(*astext.Object)

	_, err = jsonnetutil.FindObject(astParamsObject, path)
	if err != nil {
		if err := jsonnetutil.AddObject(obj, path, paramsObject.Node()); err != nil {
			return "", errors.Wrapf(err, "update %s params", componentName)
		}
	} else {
		if err := jsonnetutil.UpdateObject(obj, path, paramsObject.Node()); err != nil {
			return "", errors.Wrapf(err, "update %s params", componentName)
		}
	}

	var buf bytes.Buffer
	if err := printer.Fprint(&buf, obj); err != nil {
		return "", errors.Wrap(err, "rebuild params")
	}

	return buf.String(), nil
}

// ToMap converts a component's params to a map.
func ToMap(componentName, src string) (map[string]interface{}, error) {
	obj, err := jsonnetutil.Parse("params.libsonnet", src)
	if err != nil {
		return nil, errors.Wrap(err, "parse jsonnet")
	}

	path := []string{"components", componentName}
	child, err := jsonnetutil.FindObject(obj, path)
	if err != nil {
		return nil, err
	}

	m, err := findValues(child)
	if err != nil {
		return nil, err
	}

	paramsMap, ok := m[componentName].(map[string]interface{})
	if !ok {
		return nil, errors.Errorf("could not find %s in components", componentName)
	}

	return paramsMap, nil
}

func findValues(obj *astext.Object) (map[string]interface{}, error) {
	m := make(map[string]interface{})

	for i := range obj.Fields {
		id, err := jsonnetutil.FieldID(obj.Fields[i])
		if err != nil {
			return nil, err
		}

		switch t := obj.Fields[i].Expr2.(type) {
		default:
			return nil, errors.Errorf("unknown value type %T", t)
		case *ast.LiteralString, *ast.LiteralBoolean, *ast.LiteralNumber:
			v, err := nodeValue(t)
			if err != nil {
				return nil, err
			}
			m[id] = v
		case *ast.Array:
			array, err := arrayValues(t)
			if err != nil {
				return nil, err
			}
			m[id] = array
		case *astext.Object:
			child, err := findValues(t)
			if err != nil {
				return nil, err
			}

			m[id] = child
		}

	}

	return m, nil
}

func nodeValue(node ast.Node) (interface{}, error) {
	switch t := node.(type) {
	default:
		return nil, errors.Errorf("unknown value type %T", t)
	case *ast.LiteralString:
		return t.Value, nil
	case *ast.LiteralBoolean:
		return t.Value, nil
	case *ast.LiteralNumber:
		return t.Value, nil
	}
}

func arrayValues(array *ast.Array) ([]interface{}, error) {
	out := make([]interface{}, 0)
	for i := range array.Elements {
		v, err := nodeValue(array.Elements[i])
		if err != nil {
			return nil, errors.Errorf("arrays can't contain at %T", array.Elements[i])
		}

		out = append(out, v)
	}

	return out, nil
}
