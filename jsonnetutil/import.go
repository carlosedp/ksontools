package jsonnetutil

import (
	"io/ioutil"

	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/astext"
	"github.com/pkg/errors"
)

// Import imports jsonnet from a path.
func Import(fileName string) (*astext.Object, error) {
	if fileName == "" {
		return nil, errors.New("filename was blank")
	}

	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, errors.Wrap(err, "read lib")
	}

	return Parse(fileName, string(b))
}
