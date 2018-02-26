package component

import (
	"io/ioutil"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestYAML_Objects(t *testing.T) {
	fs := afero.NewMemMapFs()
	b, err := ioutil.ReadFile("testdata/certificate-crd.yaml")
	require.NoError(t, err)

	err = afero.WriteFile(fs, "/certificate-crd.yaml", b, 0644)
	require.NoError(t, err)

	y := YAML{
		fs:     fs,
		source: "/certificate-crd.yaml",
	}

	list, err := y.Objects()
	require.NoError(t, err)

	expected := []*unstructured.Unstructured{
		{
			Object: map[string]interface{}{
				"apiVersion": "apiextensions.k8s.io/v1beta1",
				"kind":       "CustomResourceDefinition",
				"metadata": map[string]interface{}{
					"labels": map[string]interface{}{
						"app":      "cert-manager",
						"chart":    "cert-manager-0.2.2",
						"heritage": "Tiller",
						"release":  "cert-manager",
					},
					"name": "certificates.certmanager.k8s.io",
				},
				"spec": map[string]interface{}{
					"version": "v1alpha1",
					"group":   "certmanager.k8s.io",
					"names": map[string]interface{}{
						"kind":   "Certificate",
						"plural": "certificates",
					},
					"scope": "Namespaced",
				},
			},
		},
	}

	require.Equal(t, expected, list)
}
