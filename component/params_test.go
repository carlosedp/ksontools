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
