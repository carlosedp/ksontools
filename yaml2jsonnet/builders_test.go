package yaml2jsonnet

import (
	"fmt"
	"testing"

	"github.com/bryanl/woowoo/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_buildConstructors(t *testing.T) {
	crd := "apiextensions.v1beta1.customResourceDefinition"
	in := map[string]component.Values{
		crd + ".mixin.metadata.labels": component.Values{
			Setter: crd + ".mixin.metadata.withLabels",
			Value: map[string]interface{}{
				"app":      "cert-manager",
				"chart":    "cert-manager-0.2.2",
				"heritage": "Tiller",
				"release":  "cert-manager",
			},
		},
		crd + ".mixin.metadata.name": component.Values{
			Setter: crd + ".mixin.metadata.withName",
			Value:  "certificates.certmanager.k8s.io",
		},
		crd + ".mixin.spec.group": component.Values{
			Setter: crd + ".mixin.spec.withGroup",
			Value:  "certmanager.k8s.io",
		},
		crd + ".mixin.spec.names.kind": component.Values{
			Setter: crd + ".mixin.spec.names.withKind",
			Value:  "Certificate",
		},
		crd + ".mixin.spec.names.plural": component.Values{
			Setter: crd + ".mixin.spec.names.withPlural",
			Value:  "certificates",
		},
		crd + ".mixin.spec.scope": component.Values{
			Setter: crd + ".mixin.spec.withScope",
			Value:  "Namespaced",
		},
		crd + ".mixin.spec.version": component.Values{
			Setter: crd + ".mixin.spec.withVersion",
			Value:  "v1alpha1",
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
