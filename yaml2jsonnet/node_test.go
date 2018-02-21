package yaml2jsonnet

import (
	"testing"

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

	obj, err := ImportJsonnet("testdata/k8s.libsonnet")
	require.NoError(t, err)

	node := NewNode("root", obj)

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

func TestNode_Search(t *testing.T) {
	cases := []struct {
		name  string
		path  []string
		sr    SearchResult
		isErr bool
	}{
		{
			name: "empty path",
			sr: SearchResult{
				Fields: []string{
					"admissionregistration", "apiextensions", "apiregistration", "apps",
					"authentication", "authorization", "autoscaling", "batch", "certificates",
					"core", "extensions", "meta", "networking", "policy", "rbac",
					"scheduling", "settings", "storage", "hidden",
				},
			},
		},
		{
			name: "group level",
			path: []string{"apps"},
			sr: SearchResult{
				Fields: []string{
					"v1beta1", "v1beta2",
				},
				MatchedPath: []string{"apps"},
			},
		},
		{
			name: "version level",
			path: []string{"apps", "v1beta2"},
			sr: SearchResult{
				Fields: []string{
					"apiVersion", "controllerRevision", "daemonSet", "deployment", "replicaSet",
					"scale", "statefulSet",
				},
				MatchedPath: []string{"apps", "v1beta2"},
			},
		},
		{
			name: "kind level",
			path: []string{"apps", "v1beta2", "deployment"},
			sr: SearchResult{
				Fields:      []string{"kind", "mixin"},
				Functions:   []string{"new"},
				MatchedPath: []string{"apps", "v1beta2", "deployment"},
			},
		},
		{
			name:  "invalid path",
			path:  []string{"notfound"},
			isErr: true,
		},
		{
			name: "search mixin path",
			path: []string{"apps", "v1beta2", "deployment", "metadata"},
			sr: SearchResult{
				Fields: []string{"initializers"},
				Functions: []string{
					"withAnnotations", "withAnnotationsMixin", "withClusterName",
					"withFinalizers", "withFinalizersMixin", "withGenerateName",
					"withLabels", "withLabelsMixin", "withName", "withNamespace",
					"withOwnerReferences", "withOwnerReferencesMixin",
				},
				Types:       []string{"initializersType", "ownerReferencesType"},
				MatchedPath: []string{"apps", "v1beta2", "deployment", "mixin", "metadata"},
			},
		},
		{
			name: "search for item",
			path: []string{"apps", "v1beta2", "deployment", "metadata", "name"},
			sr: SearchResult{
				Setter:      "withName",
				MatchedPath: []string{"apps", "v1beta2", "deployment", "mixin", "metadata", "name"},
			},
		},
		{
			name: "search for object",
			path: []string{"apps", "v1beta2", "deployment", "metadata", "labels"},
			sr: SearchResult{
				Setter:      "withLabels",
				MatchedPath: []string{"apps", "v1beta2", "deployment", "mixin", "metadata", "labels"},
			},
		},
		{
			name: "breaks when searching in object fields",
			path: []string{"apps", "v1beta2", "deployment", "metadata", "labels", "app"},
			sr: SearchResult{
				Setter:      "withLabels",
				MatchedPath: []string{"apps", "v1beta2", "deployment", "mixin", "metadata", "labels"},
			},
		},
	}

	obj, err := ImportJsonnet("testdata/k8s.libsonnet")
	require.NoError(t, err)

	node := NewNode("root", obj)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sr, err := node.Search(tc.path...)
			if tc.isErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.sr, sr)
			}
		})
	}

}
