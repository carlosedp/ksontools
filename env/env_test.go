package env

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_upgradeParams(t *testing.T) {
	in := `local params = import "../../components/params.libsonnet";`
	expected := `local params = std.extVar("__ksonnet/params");`

	got := upgradeParams(in)
	require.Equal(t, expected, got)
}
