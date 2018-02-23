package jsonnetutil

import (
	"bytes"
	"testing"

	"github.com/google/go-jsonnet/ast"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/astext"
	nm "github.com/ksonnet/ksonnet-lib/ksonnet-gen/nodemaker"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/printer"
	"github.com/stretchr/testify/require"
)

func TestUpdateObject(t *testing.T) {
	b := nm.NewObject()
	b.Set(nm.NewKey("c"), nm.NewStringDouble("value"))

	a := nm.NewObject()
	a.Set(nm.NewKey("b"), b)

	object := nm.NewObject()
	object.Set(nm.NewKey("a"), a)

	astObject := object.Node().(*astext.Object)

	path := []string{"a", "b", "c"}
	update := nm.NewInt(9)

	err := UpdateObject(astObject, path, update.Node())
	require.NoError(t, err)

	var got bytes.Buffer
	err = printer.Fprint(&got, astObject)
	require.NoError(t, err)

	expected := "{\n  a:: {\n    b:: {\n      c:: 9,\n    },\n  },\n}"

	require.Equal(t, expected, got.String())
}

func TestFindObject(t *testing.T) {
	b := nm.NewObject()
	b.Set(nm.NewKey("c"), nm.NewStringDouble("value"))

	a := nm.NewObject()
	a.Set(nm.NewKey("b"), b)

	object := nm.NewObject()
	object.Set(nm.NewKey("a"), a)

	astObject := object.Node().(*astext.Object)

	cases := []struct {
		name     string
		path     []string
		expected ast.Node
		isErr    bool
	}{
		{
			name:     "find nested object",
			path:     []string{"a", "b", "c"},
			expected: b.Node(),
		},
		{
			name:  "invalid path",
			path:  []string{"z"},
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			node, err := FindObject(astObject, tc.path)
			if tc.isErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, node)

			}

		})
	}

}
