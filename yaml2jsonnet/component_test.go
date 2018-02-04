package yaml2jsonnet

import (
	"io/ioutil"
	"testing"

	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/ast"
	"github.com/stretchr/testify/require"
)

func TestComponent_Generate(t *testing.T) {
	c := NewComponent()

	err := c.AddParam("option", "value")
	require.NoError(t, err)
	err = c.AddParam("int", 9)
	require.NoError(t, err)

	c.AddDeclaration(Declaration{Name: "a", Value: NewDeclarationString("a")})
	c.AddDeclaration(Declaration{Name: "b", Value: NewDeclarationString("b")})

	n := &ast.Object{}
	got, err := c.Generate(n)
	require.NoError(t, err)

	b, err := ioutil.ReadFile("testdata/declarations.libsonnet")
	require.NoError(t, err)

	expected := string(b)

	require.Equal(t, expected, got)
}
