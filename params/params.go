package params

import (
	"bytes"

	"github.com/bryanl/woowoo/jsonnetutil"
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
