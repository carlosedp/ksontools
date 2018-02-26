package node

import (
	"testing"

	"github.com/bryanl/woowoo/jsonnetutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNode_Search2(t *testing.T) {
	cases := []struct {
		name  string
		path  []string
		item  *Item
		isErr bool
	}{
		{
			name:  "empty path",
			isErr: true,
		},
		{
			name: "search for group",
			path: []string{"apps"},
			item: &Item{
				Type: ItemTypeObject,
				Path: []string{"apps"},
			},
		},
		{
			name: "search for version",
			path: []string{"apps", "v1beta2"},
			item: &Item{
				Type: ItemTypeObject,
				Path: []string{"apps", "v1beta2"},
			},
		},
		{
			name: "search for kind",
			path: []string{"apps", "v1beta2", "deployment"},
			item: &Item{
				Type: ItemTypeObject,
				Path: []string{"apps", "v1beta2", "deployment"},
			},
		},
		{
			name: "search for metadata path",
			path: []string{"apps", "v1beta2", "deployment", "metadata"},
			item: &Item{
				Type: ItemTypeObject,
				Path: []string{"apps", "v1beta2", "deployment", "mixin", "metadata"},
			},
		},
		{
			name: "search for metadata name",
			path: []string{"apps", "v1beta2", "deployment", "metadata", "name"},
			item: &Item{
				Type: ItemTypeSetter,
				Name: "apps.v1beta2.deployment.mixin.metadata.withName",
				Path: []string{"apps", "v1beta2", "deployment", "mixin", "metadata", "name"},
			},
		},
		{
			name: "search for object in metadata labels",
			path: []string{"apps", "v1beta2", "deployment", "metadata", "labels", "app"},
			item: &Item{
				Type: ItemTypeSetter,
				Name: "apps.v1beta2.deployment.mixin.metadata.withLabels",
				Path: []string{"apps", "v1beta2", "deployment", "mixin", "metadata", "labels"},
			},
		},
	}

	obj, err := jsonnetutil.Import("testdata/k8s.libsonnet")
	require.NoError(t, err)

	node := New("root", obj)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			item, err := node.Search2(tc.path...)
			if tc.isErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.item, item)
			}
		})
	}
}
