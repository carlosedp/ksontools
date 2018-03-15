package component

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestJsonnet_Name(t *testing.T) {
	app, fs := appMock("/")

	files := []string{"guestbook-ui.jsonnet", "k.libsonnet", "k8s.libsonnet", "params.libsonnet"}
	for _, file := range files {
		stageFile(t, fs, "guestbook/"+file, "/components/"+file)
	}

	c := NewJsonnet(app, "", "/components/guestbook-ui.jsonnet", "/components/params.libsonnet")

	cases := []struct {
		name         string
		isNameSpaced bool
		expected     string
	}{
		{
			name:         "wants namespaced",
			isNameSpaced: true,
			expected:     "guestbook-ui",
		},
		{
			name:         "no namespace",
			isNameSpaced: false,
			expected:     "guestbook-ui",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, c.Name(tc.isNameSpaced))
		})
	}

}

func TestJsonnet_Objects(t *testing.T) {
	app, fs := appMock("/")

	files := []string{"guestbook-ui.jsonnet", "k.libsonnet", "k8s.libsonnet", "params.libsonnet"}
	for _, file := range files {
		stageFile(t, fs, "guestbook/"+file, "/components/"+file)
	}

	files = []string{"k.libsonnet", "k8s.libsonnet"}
	for _, file := range files {
		stageFile(t, fs, "guestbook/"+file, "/lib/v1.8.7/"+file)
	}

	c := NewJsonnet(app, "", "/components/guestbook-ui.jsonnet", "/components/params.libsonnet")

	paramsStr := testdata(t, "guestbook/params.libsonnet")

	list, err := c.Objects(string(paramsStr), "default")
	require.NoError(t, err)

	expected := []*unstructured.Unstructured{
		{
			Object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "Service",
				"metadata": map[string]interface{}{
					"name": "guiroot",
				},
				"spec": map[string]interface{}{
					"ports": []interface{}{
						map[string]interface{}{
							"port":       int64(80),
							"targetPort": int64(80),
						},
					},
					"selector": map[string]interface{}{
						"app": "guiroot",
					},
					"type": "ClusterIP",
				},
			},
		},
		{
			Object: map[string]interface{}{
				"apiVersion": "apps/v1beta1",
				"kind":       "Deployment",
				"metadata": map[string]interface{}{
					"name": "guiroot",
				},
				"spec": map[string]interface{}{
					"replicas": int64(1),
					"template": map[string]interface{}{
						"metadata": map[string]interface{}{
							"labels": map[string]interface{}{
								"app": "guiroot",
							},
						},
						"spec": map[string]interface{}{
							"containers": []interface{}{
								map[string]interface{}{
									"image": "gcr.io/heptio-images/ks-guestbook-demo:0.1",
									"name":  "guiroot",
									"ports": []interface{}{
										map[string]interface{}{
											"containerPort": int64(80),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	require.Equal(t, expected, list)
}

func TestJsonnet_Params(t *testing.T) {
	app, fs := appMock("/")

	files := []string{"guestbook-ui.jsonnet", "k.libsonnet", "k8s.libsonnet", "params.libsonnet"}
	for _, file := range files {
		stageFile(t, fs, "guestbook/"+file, "/components/"+file)
	}

	c := NewJsonnet(app, "", "/components/guestbook-ui.jsonnet", "/components/params.libsonnet")

	params, err := c.Params()
	require.NoError(t, err)

	expected := []NamespaceParameter{
		{
			Component: "guestbook-ui",
			Index:     "0",
			Key:       "containerPort",
			Value:     "80",
		},
		{
			Component: "guestbook-ui",
			Index:     "0",
			Key:       "image",
			Value:     `"gcr.io/heptio-images/ks-guestbook-demo:0.1"`,
		},
		{
			Component: "guestbook-ui",
			Index:     "0",
			Key:       "name",
			Value:     `"guiroot"`,
		},
		{
			Component: "guestbook-ui",
			Index:     "0",
			Key:       "obj",
			Value:     `{"a":"b"}`,
		},
		{
			Component: "guestbook-ui",
			Index:     "0",
			Key:       "replicas",
			Value:     "1",
		},
		{
			Component: "guestbook-ui",
			Index:     "0",
			Key:       "servicePort",
			Value:     "80",
		},
		{
			Component: "guestbook-ui",
			Index:     "0",
			Key:       "type",
			Value:     `"ClusterIP"`,
		},
	}

	require.Equal(t, expected, params)
}

func TestJsonnet_Summarize(t *testing.T) {
	app, fs := appMock("/")

	files := []string{"guestbook-ui.jsonnet", "k.libsonnet", "k8s.libsonnet", "params.libsonnet"}
	for _, file := range files {
		stageFile(t, fs, "guestbook/"+file, "/components/"+file)
	}

	c := NewJsonnet(app, "", "/components/guestbook-ui.jsonnet", "/components/params.libsonnet")

	got, err := c.Summarize()
	require.NoError(t, err)

	expected := []Summary{
		{ComponentName: "guestbook-ui", IndexStr: "0", Type: "jsonnet"},
	}

	require.Equal(t, expected, got)
}

func TestJsonnet_SetParam(t *testing.T) {
	app, fs := appMock("/")

	files := []string{"guestbook-ui.jsonnet", "k.libsonnet", "k8s.libsonnet", "params.libsonnet"}
	for _, file := range files {
		stageFile(t, fs, "guestbook/"+file, "/components/"+file)
	}

	c := NewJsonnet(app, "", "/components/guestbook-ui.jsonnet", "/components/params.libsonnet")

	err := c.SetParam([]string{"replicas"}, 4, ParamOptions{})
	require.NoError(t, err)

	b, err := afero.ReadFile(fs, "/components/params.libsonnet")
	require.NoError(t, err)

	expected := testdata(t, "guestbook/set-params.libsonnet")

	require.Equal(t, string(expected), string(b))
}

func TestJsonnet_DeleteParam(t *testing.T) {
	app, fs := appMock("/")

	files := []string{"guestbook-ui.jsonnet", "k.libsonnet", "k8s.libsonnet", "params.libsonnet"}
	for _, file := range files {
		stageFile(t, fs, "guestbook/"+file, "/components/"+file)
	}

	c := NewJsonnet(app, "", "/components/guestbook-ui.jsonnet", "/components/params.libsonnet")

	err := c.DeleteParam([]string{"replicas"}, ParamOptions{})
	require.NoError(t, err)

	b, err := afero.ReadFile(fs, "/components/params.libsonnet")
	require.NoError(t, err)

	expected := testdata(t, "guestbook/delete-params.libsonnet")

	require.Equal(t, string(expected), string(b))
}
