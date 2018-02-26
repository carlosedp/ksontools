package yaml2jsonnet

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/bryanl/woowoo/jsonnetutil"
	"github.com/stretchr/testify/require"
)

func TestDocument_GenerateComponent(t *testing.T) {
	f, err := os.Open("testdata/certificate-crd.yaml")
	require.NoError(t, err)

	defer f.Close()

	node, err := jsonnetutil.Import("testdata/k8s.libsonnet")
	require.NoError(t, err)

	doc, err := NewDocument("certificateCrd", f, node)
	require.NoError(t, err)

	got, err := doc.GenerateComponent()
	require.NoError(t, err)

	expected, err := ioutil.ReadFile("testdata/cert-manager.jsonnet")
	require.NoError(t, err)

	require.Equal(t, string(expected), got)
}

func Test_mixinConstructorName(t *testing.T) {
	name := "apiextensions.v1beta1.customResourceDefinition.mixin.metadata"
	got := mixinConstructorName(name)
	expected := "createCustomResourceDefinitionMetadata"
	require.Equal(t, expected, got)
}

func Test_mixinObjectName(t *testing.T) {
	name := "apiextensions.v1beta1.customResourceDefinition.mixin.metadata"
	got := mixinObjectName(name)
	expected := "customResourceDefinitionMetadata"
	require.Equal(t, expected, got)
}

func Test_paramName(t *testing.T) {
	name := "apiextensions.v1beta1.customResourceDefinition.mixin.spec.group"
	got := paramName(name)
	expected := "crdSpecGroup"
	require.Equal(t, expected, got)
}
