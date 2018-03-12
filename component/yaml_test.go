package component

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestYAML_Params(t *testing.T) {
	app, fs := appMock("/")

	stageFile(t, fs, "params-mixed.libsonnet", "/components/params.libsonnet")
	stageFile(t, fs, "deployment.yaml", "/components/deployment.yaml")
	stageFile(t, fs, "k8s.libsonnet", "/lib/v1.8.7/k8s.libsonnet")

	y := NewYAML(app, "/components/deployment.yaml", "/components/params.libsonnet")
	params, err := y.Params()
	require.NoError(t, err)

	require.Len(t, params, 1)

	param := params[0]
	expected := NamespaceParameter{
		Component: "deployment",
		Index:     "0",
		Key:       "metadata.labels",
		Value:     `{"label1":"label1","label2":"label2"}`,
	}
	require.Equal(t, expected, param)
}

func TestYAML_Params_literal(t *testing.T) {
	app, fs := appMock("/")

	stageFile(t, fs, "params-mixed.libsonnet", "/params.libsonnet")
	stageFile(t, fs, "rbac.yaml", "/rbac.yaml")
	stageFile(t, fs, "k8s.libsonnet", "/lib/v1.8.7/k8s.libsonnet")

	y := NewYAML(app, "/rbac.yaml", "/params.libsonnet")
	params, err := y.Params()
	require.NoError(t, err)

	require.Len(t, params, 1)

	param := params[0]
	expected := NamespaceParameter{
		Component: "rbac",
		Index:     "1",
		Key:       "metadata.name",
		Value:     "cert-manager2",
	}
	require.Equal(t, expected, param)
}

func TestYAML_Objects_no_params(t *testing.T) {
	app, fs := appMock("/")

	stageFile(t, fs, "certificate-crd.yaml", "/certificate-crd.yaml")
	stageFile(t, fs, "params-no-entry.libsonnet", "/params.libsonnet")

	y := NewYAML(app, "/certificate-crd.yaml", "/params.libsonnet")

	list, err := y.Objects("")
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
func TestYAML_Objects_no_params_with_json(t *testing.T) {
	app, fs := appMock("/")

	stageFile(t, fs, "certificate-crd.json", "/certificate-crd.json")
	stageFile(t, fs, "params-no-entry.libsonnet", "/params.libsonnet")

	y := NewYAML(app, "/certificate-crd.json", "/params.libsonnet")

	list, err := y.Objects("")
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
	app, fs := appMock("/")

	stageFile(t, fs, "certificate-crd.yaml", "/certificate-crd.yaml")
	stageFile(t, fs, "params-no-entry.libsonnet", "/params.libsonnet")

	y := NewYAML(app, "/certificate-crd.yaml", "/params.libsonnet")

	list, err := y.Objects("")
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
	app, fs := appMock("/")

	stageFile(t, fs, "certificate-crd.yaml", "/certificate-crd.yaml")
	stageFile(t, fs, "params-with-entry.libsonnet", "/params.libsonnet")

	y := NewYAML(app, "/certificate-crd.yaml", "/params.libsonnet")

	list, err := y.Objects("")
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

func TestYAML_SetParam(t *testing.T) {
	app, fs := appMock("/")

	stageFile(t, fs, "certificate-crd.yaml", "/certificate-crd.yaml")
	stageFile(t, fs, "params-no-entry.libsonnet", "/params.libsonnet")

	y := NewYAML(app, "/certificate-crd.yaml", "/params.libsonnet")

	err := y.SetParam([]string{"spec", "version"}, "v2", ParamOptions{})
	require.NoError(t, err)

	b, err := afero.ReadFile(fs, "/params.libsonnet")
	require.NoError(t, err)

	expected := testdata(t, "updated-yaml-params.libsonnet")

	require.Equal(t, string(expected), string(b))
}

func TestYAML_DeleteParam(t *testing.T) {
	app, fs := appMock("/")

	stageFile(t, fs, "certificate-crd.yaml", "/certificate-crd.yaml")
	stageFile(t, fs, "params-with-entry.libsonnet", "/params.libsonnet")

	y := NewYAML(app, "/certificate-crd.yaml", "/params.libsonnet")

	err := y.DeleteParam([]string{"spec", "version"}, ParamOptions{})
	require.NoError(t, err)

	b, err := afero.ReadFile(fs, "/params.libsonnet")
	require.NoError(t, err)

	expected := testdata(t, "params-delete-entry.libsonnet")

	require.Equal(t, string(expected), string(b))
}

func TestYAML_Summarize(t *testing.T) {
	app, fs := appMock("/")

	stageFile(t, fs, "rbac.yaml", "/components/rbac.yaml")
	stageFile(t, fs, "params-no-entry.libsonnet", "/components/params.libsonnet")

	y := NewYAML(app, "/components/rbac.yaml", "/components/params.libsonnet")

	list, err := y.Summarize()
	require.NoError(t, err)

	expected := []Summary{
		{
			ComponentName: "rbac",
			IndexStr:      "0",
			Type:          "yaml",
			APIVersion:    "rbac.authorization.k8s.io/v1beta1",
			Kind:          "ClusterRole",
			Name:          "cert-manager",
		},
		{
			ComponentName: "rbac",
			IndexStr:      "1",
			Type:          "yaml",
			APIVersion:    "rbac.authorization.k8s.io/v1beta1",
			Kind:          "ClusterRoleBinding",
			Name:          "cert-manager",
		},
	}

	require.Equal(t, expected, list)
}

func TestYAML_Summarize_json(t *testing.T) {
	app, fs := appMock("/")

	stageFile(t, fs, "certificate-crd.json", "/components/certificate-crd.json")
	stageFile(t, fs, "params-no-entry.libsonnet", "/components/params.libsonnet")

	y := NewYAML(app, "/components/certificate-crd.json", "/components/params.libsonnet")

	list, err := y.Summarize()
	require.NoError(t, err)

	expected := []Summary{
		{
			ComponentName: "certificate-crd",
			IndexStr:      "0",
			Type:          "json",
			APIVersion:    "apiextensions.k8s.io/v1beta1",
			Kind:          "CustomResourceDefinition",
			Name:          "certificates_certmanager_k8s_io",
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

	err := mergeMaps(m1, m2, nil)
	require.NoError(t, err)
	require.Equal(t, expected, m1)
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

func testdata(t *testing.T, name string) []byte {
	b, err := ioutil.ReadFile("testdata/" + name)
	require.NoError(t, err, "read testdata %s", name)
	return b
}
