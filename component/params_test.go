package component

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_applyGlobals(t *testing.T) {
	myParams, err := ioutil.ReadFile("testdata/params-global.libsonnet")
	require.NoError(t, err)

	got, err := applyGlobals(string(myParams))
	require.NoError(t, err)

	expected, err := ioutil.ReadFile("testdata/params-global-expected.json")
	require.NoError(t, err)

	require.Equal(t, string(expected), got)
}

func Test_patchJSON(t *testing.T) {
	jsonObject, err := ioutil.ReadFile("testdata/rbac-1.json")
	require.NoError(t, err)

	patch, err := ioutil.ReadFile("testdata/patch.json")
	require.NoError(t, err)

	got, err := patchJSON(string(jsonObject), string(patch), "rbac-1")
	require.NoError(t, err)

	expected, err := ioutil.ReadFile("testdata/rbac-1-patched.json")
	require.NoError(t, err)

	require.Equal(t, string(expected), got)
}
