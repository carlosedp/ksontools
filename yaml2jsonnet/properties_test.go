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
