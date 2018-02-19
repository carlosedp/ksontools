package yaml2jsonnet

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProperties_Paths(t *testing.T) {

	var (
		deploymentBase = []string{"apps", "v1beta2", "deployment"}
		crdBase        = []string{"hidden", "apiextensions", "v1beta1", "customResourceDefinition"}
	)

	cases := []struct {
		name     string
		expected []PropertyPath
	}{
		{
			name: "testdata/deployment.yaml",
			expected: []PropertyPath{
				{Path: append(deploymentBase, "metadata", "labels", "app")},
				{Path: append(deploymentBase, "metadata", "name")},
				{Path: append(deploymentBase, "spec", "replicas")},
				{Path: append(deploymentBase, "spec", "selector", "matchLabels", "app")},
				{Path: append(deploymentBase, "spec", "template", "metadata", "labels", "app")},
				{Path: append(deploymentBase, "spec", "template", "spec", "containers")},
			},
		},
		{
			name: "testdata/certificate-crd.yaml",
			expected: []PropertyPath{
				{Path: append(crdBase, "metadata", "labels", "app")},
				{Path: append(crdBase, "metadata", "labels", "chart")},
				{Path: append(crdBase, "metadata", "labels", "heritage")},
				{Path: append(crdBase, "metadata", "labels", "release")},
				{Path: append(crdBase, "metadata", "name")},
				{Path: append(crdBase, "spec", "group")},
				{Path: append(crdBase, "spec", "names", "kind")},
				{Path: append(crdBase, "spec", "names", "plural")},
				{Path: append(crdBase, "spec", "scope")},
				{Path: append(crdBase, "spec", "version")},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(tc.name)
			require.NoError(t, err)
			defer f.Close()

			ts, props, err := importYaml(f)
			require.NoError(t, err)

			gvk, err := ts.GVK()
			require.NoError(t, err)

			got := props.Paths(gvk)
			require.Equal(t, tc.expected, got)
		})
	}

}

func TestProperties_Value(t *testing.T) {
	f, err := os.Open("testdata/deployment.yaml")
	require.NoError(t, err)
	defer f.Close()

	props := Properties{}

	_, props, err = importYaml(f)
	require.NoError(t, err)

	cases := []struct {
		name     string
		path     []string
		expected interface{}
	}{
		{
			name:     "string",
			path:     []string{"metadata", "name"},
			expected: "nginx-deployment",
		},
		{
			name:     "int",
			path:     []string{"spec", "replicas"},
			expected: 3,
		},
		{
			name:     "array",
			path:     []string{"spec", "template", "spec", "containers"},
			expected: []interface{}([]interface{}{map[interface{}]interface{}{"name": "nginx", "image": "nginx:1.7.9", "ports": []interface{}{map[interface{}]interface{}{"containerPort": 80}}}}),
		},
		{
			name:     "object",
			path:     []string{"metadata", "labels"},
			expected: map[interface{}]interface{}(map[interface{}]interface{}{"app": "nginx"}),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := props.Value(tc.path)
			require.NoError(t, err)

			require.Equal(t, tc.expected, got)
		})
	}

}
