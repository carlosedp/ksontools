package yaml2jsonnet

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_buildConstructors(t *testing.T) {
	crd := "apiextensions.v1beta1.customResourceDefinition"
	in := map[string]documentValues{
		crd + ".mixin.metadata.labels": documentValues{
			setter: crd + ".mixin.metadata.withLabels",
			value: map[string]interface{}{
				"app":      "cert-manager",
				"chart":    "cert-manager-0.2.2",
				"heritage": "Tiller",
				"release":  "cert-manager",
			},
		},
		crd + ".mixin.metadata.name": documentValues{
			setter: crd + ".mixin.metadata.withName",
			value:  "certificates.certmanager.k8s.io",
		},
		crd + ".mixin.spec.group": documentValues{
			setter: crd + ".mixin.spec.withGroup",
			value:  "certmanager.k8s.io",
		},
		crd + ".mixin.spec.names.kind": documentValues{
			setter: crd + ".mixin.spec.names.withKind",
			value:  "Certificate",
		},
		crd + ".mixin.spec.names.plural": documentValues{
			setter: crd + ".mixin.spec.names.withPlural",
			value:  "certificates",
		},
		crd + ".mixin.spec.scope": documentValues{
			setter: crd + ".mixin.spec.withScope",
			value:  "Namespaced",
		},
		crd + ".mixin.spec.version": documentValues{
			setter: crd + ".mixin.spec.withVersion",
			value:  "v1alpha1",
		},
	}

	got, err := buildConstructors(in)
	require.NoError(t, err)

	expected := map[string][]ctorArgument{
		fmt.Sprintf("%s.mixin.metadata", crd): []ctorArgument{
			{
				setter:    "withLabels",
				paramName: "crdMetadataLabels",
				paramValue: map[string]interface{}{
					"app":      "cert-manager",
					"chart":    "cert-manager-0.2.2",
					"release":  "cert-manager",
					"heritage": "Tiller",
				},
			},
			{
				setter:     "withName",
				paramName:  "crdMetadataName",
				paramValue: "certificates.certmanager.k8s.io",
			},
		},
		fmt.Sprintf("%s.mixin.spec.names", crd): []ctorArgument{
			{
				setter:     "withKind",
				paramName:  "crdSpecNamesKind",
				paramValue: "Certificate",
			},
			{
				setter:     "withPlural",
				paramName:  "crdSpecNamesPlural",
				paramValue: "certificates",
			},
		},
		fmt.Sprintf("%s.mixin.spec", crd): []ctorArgument{
			{
				setter:     "withGroup",
				paramName:  "crdSpecGroup",
				paramValue: "certmanager.k8s.io",
			},
			{
				setter:     "withScope",
				paramName:  "crdSpecScope",
				paramValue: "Namespaced",
			},
			{
				setter:     "withVersion",
				paramName:  "crdSpecVersion",
				paramValue: "v1alpha1",
			},
		},
	}

	require.Equal(t, expected, got)
}

func Test_parseSetterNamespace(t *testing.T) {
	cases := []struct {
		name     string
		expected []string
		isErr    bool
	}{
		{
			name:     "foo.mixin.a.withFirst",
			expected: []string{"foo.mixin.a", "withFirst"},
		},
		{
			name:  "short",
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ns, setter, err := parseSetterNamespace(tc.name)
			if tc.isErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				assert.Equal(t, tc.expected[0], ns, "namespace")
				assert.Equal(t, tc.expected[1], setter, "setter")
			}
		})
	}
}
