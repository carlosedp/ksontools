package component

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNamespace_Components(t *testing.T) {
	app, fs := appMock("/app")

	stageFile(t, fs, "certificate-crd.yaml", "/app/components/ns1/certificate-crd.yaml")
	stageFile(t, fs, "params-with-entry.libsonnet", "/app/components/ns1/params.libsonnet")
	stageFile(t, fs, "params-no-entry.libsonnet", "/app/components/params.libsonnet")

	cases := []struct {
		name   string
		nsName string
		count  int
	}{
		{
			name: "no components",
		},
		{
			name:   "with components",
			nsName: "ns1",
			count:  1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			ns, err := GetNamespace(app, tc.nsName)
			require.NoError(t, err)

			assert.Equal(t, tc.nsName, ns.Name())
			components, err := ns.Components()
			require.NoError(t, err)

			assert.Len(t, components, tc.count)
		})
	}

}
