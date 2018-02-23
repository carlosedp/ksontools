package jsonnetutil

import (
	"github.com/bryanl/woowoo/pkg/docparser"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/astext"
	"github.com/pkg/errors"
)

// Parse converts a jsonnet snippet to AST.
func Parse(fileName, src string) (*astext.Object, error) {
	tokens, err := docparser.Lex(fileName, src)
	if err != nil {
		return nil, errors.Wrap(err, "lex lib")
	}

	node, err := docparser.Parse(tokens)
	if err != nil {
		return nil, errors.Wrap(err, "parse lib")
	}

	root, ok := node.(*astext.Object)
	if !ok {
		return nil, errors.New("root was not an object")
	}

	return root, nil
}
