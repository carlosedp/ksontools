package component

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestJsonnet_Name(t *testing.T) {
	app, fs := appMock("/")

	files := []string{"guestbook-ui.jsonnet", "k.libsonnet", "k8s.libsonnet", "params.libsonnet"}
	for _, file := range files {
		stageFile(t, fs, "guestbook/"+file, "/components/"+file)
	}

	c := NewJsonnet(app, "/components/guestbook-ui.jsonnet", "/components/params.libsonnet")

	require.Equal(t, "guestbook-ui", c.Name())
}

func TestJsonnet_Params(t *testing.T) {
	app, fs := appMock("/")

	files := []string{"guestbook-ui.jsonnet", "k.libsonnet", "k8s.libsonnet", "params.libsonnet"}
	for _, file := range files {
		stageFile(t, fs, "guestbook/"+file, "/components/"+file)
	}

	c := NewJsonnet(app, "/components/guestbook-ui.jsonnet", "/components/params.libsonnet")

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

	c := NewJsonnet(app, "/components/guestbook-ui.jsonnet", "/components/params.libsonnet")

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

	c := NewJsonnet(app, "/components/guestbook-ui.jsonnet", "/components/params.libsonnet")

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

	c := NewJsonnet(app, "/components/guestbook-ui.jsonnet", "/components/params.libsonnet")

	err := c.DeleteParam([]string{"replicas"}, ParamOptions{})
	require.NoError(t, err)

	b, err := afero.ReadFile(fs, "/components/params.libsonnet")
	require.NoError(t, err)

	expected := testdata(t, "guestbook/delete-params.libsonnet")

	require.Equal(t, string(expected), string(b))
}
