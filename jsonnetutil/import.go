package jsonnetutil

import (
	"github.com/bryanl/woowoo/pkg/docparser"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/astext"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// Import imports jsonnet from a path.
func Import(fileName string) (*astext.Object, error) {
	fs := afero.NewOsFs()
	return ImportFromFs(fileName, fs)
}

// ImportFromFs imports jsonnet from a path on an afero filesystem.
func ImportFromFs(fileName string, fs afero.Fs) (*astext.Object, error) {
	if fileName == "" {
		return nil, errors.New("filename was blank")
	}

	b, err := afero.ReadFile(fs, fileName)
	if err != nil {
		return nil, errors.Wrap(err, "read lib")
	}

	return Parse(fileName, string(b))

}

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
