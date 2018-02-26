package component

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGVK_Group(t *testing.T) {
	cases := []struct {
		name      string
		groupPath []string
		expected  []string
	}{
		{
			name:      "unqualified group",
			groupPath: []string{"apps"},
			expected:  []string{"apps"},
		},
		{
			name:      "qualified group",
			groupPath: []string{"apiextensions.k8s.io"},
			expected:  []string{"apiextensions"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gvk := GVK{
				GroupPath: tc.groupPath,
				Version:   "v1",
				Kind:      "deployment",
			}

			group := gvk.Group()
			require.Equal(t, tc.expected, group)
		})
	}
}

func TestGVK_Path(t *testing.T) {
	gvk := GVK{
		GroupPath: []string{"apps"},
		Version:   "v1",
		Kind:      "deployment",
	}

	expected := []string{"apps", "v1", "deployment"}
	require.Equal(t, expected, gvk.Path())
}
