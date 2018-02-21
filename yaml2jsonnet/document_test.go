package yaml2jsonnet

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDocument_GenerateComponent(t *testing.T) {
	f, err := os.Open("testdata/certificate-crd.yaml")
	require.NoError(t, err)

	defer f.Close()

	node, err := ImportJsonnet("testdata/k8s.libsonnet")
	require.NoError(t, err)

	doc, err := NewDocument(f, node)
	require.NoError(t, err)

	// got, err := doc.GenerateComponent()
	// require.NoError(t, err)

	// expected, err := ioutil.ReadFile("testdata/cert-manager.jsonnet")
	// require.NoError(t, err)

	// require.Equal(t,)

	resolved, err := doc.resolvedPaths2()
	require.NoError(t, err)

	crd := "apiextensions.v1beta1.customResourceDefinition."

	expected := map[string]documentValues{
		crd + "mixin.metadata.labels": documentValues{
			setter: crd + "mixin.metadata.withLabels",
			value: map[interface{}]interface{}{
				"app":      "cert-manager",
				"chart":    "cert-manager-0.2.2",
				"release":  "cert-manager",
				"heritage": "Tiller",
			},
		},
		crd + "mixin.metadata.name": documentValues{
			setter: crd + "mixin.metadata.withName",
			value:  "certificates.certmanager.k8s.io",
		},
		crd + "mixin.spec.group": documentValues{
			setter: crd + "mixin.spec.withGroup",
			value:  "certmanager.k8s.io",
		},
		crd + "mixin.spec.names.kind": documentValues{
			setter: crd + "mixin.spec.names.withKind",
			value:  "Certificate",
		},
		crd + "mixin.spec.names.plural": documentValues{
			setter: crd + "mixin.spec.names.withPlural",
			value:  "certificates",
		},
		crd + "mixin.spec.scope": documentValues{
			setter: crd + "mixin.spec.withScope",
			value:  "Namespaced",
		},
		crd + "mixin.spec.version": documentValues{
			setter: crd + "mixin.spec.withVersion",
			value:  "v1alpha1",
		},
	}

	require.Equal(t, expected, resolved)
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
