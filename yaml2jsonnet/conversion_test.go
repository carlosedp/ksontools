package yaml2jsonnet

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_generateComponentFileName(t *testing.T) {
	cases := []struct {
		name     string
		expected string
	}{
		{
			name:     "file.yaml",
			expected: "file",
		},
		{
			name:     "complex-file.yaml",
			expected: "complexFile",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := generateComponentName(tc.name)
			require.Equal(t, tc.expected, got)
		})
	}
}
