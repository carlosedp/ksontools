package yaml2jsonnet

import (
	"io/ioutil"
	"testing"

	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/ast"
	"github.com/stretchr/testify/require"
)

func TestComponent_Generate(t *testing.T) {
	c := NewComponent()

	ds := NewDeclarationString("a")
	c.AddDeclaration(Declaration{Name: "a", Value: ds})

	n := &ast.Object{}
	got, err := c.Declarations(n)
	require.NoError(t, err)

	b, err := ioutil.ReadFile("testdata/declarations.libsonnet")
	require.NoError(t, err)

	expected := string(b)

	require.Equal(t, expected, got)
}
