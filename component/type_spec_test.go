package component

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTypeSpec(t *testing.T) {
	cases := []struct {
		name    string
		kind    string
		version string
		isErr   bool
	}{
		{
			name:    "with kind and version",
			kind:    "Deployment",
			version: "v1",
		},
		{
			name:    "missing kind or version",
			version: "v1",
			isErr:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewTypeSpec(tc.version, tc.kind)

			if tc.isErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestTypeSpec_Group(t *testing.T) {
	cases := []struct {
		name    string
		version string
		group   []string
	}{
		{
			name:    "with an explicit group in the version",
			version: "group/v1",
			group:   []string{"group"},
		},
		{
			name:    "without an explicit group in the version",
			version: "v1",
			group:   []string{"core"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ts, err := NewTypeSpec(tc.version, "kind")
			require.NoError(t, err)

			require.Equal(t, tc.group, ts.Group())
		})
	}
}

func TestTypeSpec_Version(t *testing.T) {
	cases := []struct {
		name     string
		version  string
		expected string
	}{
		{
			name:     "with an explicit group in the version",
			version:  "group/v1",
			expected: "v1",
		},
		{
			name:     "without an explicit group in the version",
			version:  "v1",
			expected: "v1",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ts, err := NewTypeSpec(tc.version, "kind")
			require.NoError(t, err)

			require.Equal(t, tc.expected, ts.Version())
		})
	}
}

func TestTypeSpec_Kind(t *testing.T) {
	cases := []struct {
		name     string
		kind     string
		expected string
	}{
		{
			name:     "with an explicit group in the version",
			kind:     "Group",
			expected: "group",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ts, err := NewTypeSpec("v1", tc.kind)
			require.NoError(t, err)

			require.Equal(t, tc.expected, ts.Kind())
		})
	}
}
