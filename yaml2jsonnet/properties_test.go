package yaml2jsonnet

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProperties_Paths(t *testing.T) {
	f, err := os.Open("testdata/deployment.yaml")
	require.NoError(t, err)
	defer f.Close()

	props := Properties{}

	ts, err := importYaml(f, props)
	require.NoError(t, err)

	gvk, err := ts.GVK()
	require.NoError(t, err)

	expected := []PropertyPath{
		{Path: []string{"apps", "v1beta2", "deployment", "metadata", "labels", "app"}},
		{Path: []string{"apps", "v1beta2", "deployment", "metadata", "name"}},
		{Path: []string{"apps", "v1beta2", "deployment", "spec", "replicas"}},
		{Path: []string{"apps", "v1beta2", "deployment", "spec", "selector", "matchLabels", "app"}},
		{Path: []string{"apps", "v1beta2", "deployment", "spec", "template", "metadata", "labels", "app"}},
		{Path: []string{"apps", "v1beta2", "deployment", "spec", "template", "spec", "containers"}},
	}

	got := props.Paths(gvk)
	require.Equal(t, expected, got)
}

func TestProperties_Value(t *testing.T) {
	f, err := os.Open("testdata/deployment.yaml")
	require.NoError(t, err)
	defer f.Close()

	props := Properties{}

	_, err = importYaml(f, props)
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
