package params

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	paramsSource, err := ioutil.ReadFile("testdata/params.libsonnet")
	require.NoError(t, err)

	componentName := "guestbook-ui"

	params := map[string]interface{}{
		"containerPort": 80,
		"image":         "gcr.io/heptio-images/ks-guestbook-demo:0.2",
		"name":          "guestbook-ui",
		"replicas":      5,
		"servicePort":   80,
		"type":          "NodePort",
	}

	got, err := Update(componentName, string(paramsSource), params)
	require.NoError(t, err)

	expected, err := ioutil.ReadFile("testdata/updated.libsonnet")
	require.NoError(t, err)

	fmt.Printf("got:\n%s\n", got)
	fmt.Printf("expected:\n%s\n", expected)

	require.Equal(t, string(expected), got)
}
