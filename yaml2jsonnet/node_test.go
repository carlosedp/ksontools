package yaml2jsonnet

import (
	"testing"

	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/astext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
					"admissionregistration", "apps", "authentication", "authorization", "autoscaling",
					"batch", "certificates", "core", "extensions", "meta", "networking", "policy",
					"rbac", "scheduling", "settings", "storage", "hidden",
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
					"apiVersion", "controllerRevision", "controllerRevisionList", "daemonSet",
					"daemonSetList", "deployment", "deploymentList", "replicaSet",
					"replicaSetList", "scale", "statefulSet", "statefulSetList",
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

	astNode, err := ImportJsonnet("testdata/k8s.libsonnet")
	require.NoError(t, err)

	obj, ok := astNode.(*astext.Object)
	require.True(t, ok)

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
