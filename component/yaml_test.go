package component

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestYAML_Objects_no_params(t *testing.T) {
	fs := afero.NewMemMapFs()
	stageFile(t, fs, "certificate-crd.yaml", "/certificate-crd.yaml")

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

func TestYAML_Objects_params_exist_with_no_entry(t *testing.T) {
	fs := afero.NewMemMapFs()

	stageFile(t, fs, "certificate-crd.yaml", "/certificate-crd.yaml")
	stageFile(t, fs, "params-no-entry.libsonnet", "/params.libsonnet")

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

func TestYAML_Objects_params_exist_with_entry(t *testing.T) {
	fs := afero.NewMemMapFs()

	stageFile(t, fs, "certificate-crd.yaml", "/certificate-crd.yaml")
	stageFile(t, fs, "params-with-entry.libsonnet", "/params.libsonnet")

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
					"version": "v2",
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

func Test_mapToPaths(t *testing.T) {
	m := map[string]interface{}{
		"metadata": map[string]interface{}{
			"name": "name",
			"labels": map[string]interface{}{
				"label1": "label1",
			},
		},
	}

	lookup := map[string]bool{
		// "metadata":        true,
		"metadata.name":   true,
		"metadata.labels": true,
	}

	got := mapToPaths(m, lookup, nil)

	expected := []paramPath{
		{path: []string{"metadata", "labels"}, value: map[string]interface{}{"label1": "label1"}},
		{path: []string{"metadata", "name"}, value: "name"},
	}

	require.Equal(t, expected, got)
}

func Test_mergeMaps(t *testing.T) {
	m1 := map[string]interface{}{
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
			"version": "v1",
			"group":   "certmanager.k8s.io",
			"names": map[string]interface{}{
				"kind":   "Certificate",
				"plural": "certificates",
			},
			"scope": "Namespaced",
		},
	}

	m2 := map[string]interface{}{
		"spec": map[string]interface{}{
			"version": "v2",
		},
	}

	expected := map[string]interface{}{
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
			"version": "v2",
			"group":   "certmanager.k8s.io",
			"names": map[string]interface{}{
				"kind":   "Certificate",
				"plural": "certificates",
			},
			"scope": "Namespaced",
		},
	}

	got, err := mergeMaps(m1, m2, nil)
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func stageFile(t *testing.T, fs afero.Fs, src, dest string) {
	in := filepath.Join("testdata", src)

	b, err := ioutil.ReadFile(in)
	require.NoError(t, err)

	dir := filepath.Dir(dest)
	err = fs.MkdirAll(dir, 0755)
	require.NoError(t, err)

	err = afero.WriteFile(fs, dest, b, 0644)
	require.NoError(t, err)
}
