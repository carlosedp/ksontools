package yaml2jsonnet

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/ast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNode_Search(t *testing.T) {
	cases := []struct {
		name        string
		path        []string
		matchedPath []string
		sr          SearchResult
		isErr       bool
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
			name:        "group level",
			path:        []string{"apps"},
			matchedPath: []string{"apps"},
			sr: SearchResult{
				Fields: []string{
					"v1beta1", "v1beta2",
				},
			},
		},
		{
			name:        "version level",
			path:        []string{"apps", "v1beta2"},
			matchedPath: []string{"apps", "v1beta2"},
			sr: SearchResult{
				Fields: []string{
					"apiVersion", "controllerRevision", "controllerRevisionList", "daemonSet",
					"daemonSetList", "deployment", "deploymentList", "replicaSet",
					"replicaSetList", "scale", "statefulSet", "statefulSetList",
				},
			},
		},
		{
			name:        "kind level",
			path:        []string{"apps", "v1beta2", "deployment"},
			matchedPath: []string{"apps", "v1beta2", "deployment"},
			sr: SearchResult{
				Fields:    []string{"kind", "mixin"},
				Functions: []string{"new"},
			},
		},
		{
			name:  "invalid path",
			path:  []string{"notfound"},
			isErr: true,
		},
		{
			name:        "search mixin path",
			path:        []string{"apps", "v1beta2", "deployment", "metadata"},
			matchedPath: []string{"apps", "v1beta2", "deployment", "mixin", "metadata"},
			sr: SearchResult{
				Fields: []string{"initializers"},
				Functions: []string{
					"withAnnotations", "withAnnotationsMixin", "withClusterName",
					"withFinalizers", "withFinalizersMixin", "withGenerateName",
					"withLabels", "withLabelsMixin", "withName", "withNamespace",
					"withOwnerReferences", "withOwnerReferencesMixin",
				},
				Types: []string{"initializersType", "ownerReferencesType"},
			},
		},
		{
			name:        "search for item",
			path:        []string{"apps", "v1beta2", "deployment", "metadata", "name"},
			matchedPath: []string{"apps", "v1beta2", "deployment", "mixin", "metadata", "name"},
			sr:          SearchResult{Setter: "withName"},
		},
		{
			name:        "search for object",
			path:        []string{"apps", "v1beta2", "deployment", "metadata", "labels"},
			matchedPath: []string{"apps", "v1beta2", "deployment", "mixin", "metadata", "labels"},
			sr:          SearchResult{Setter: "withLabels"},
		},
		{
			name:        "breaks when searching in object fields",
			path:        []string{"apps", "v1beta2", "deployment", "metadata", "labels", "app"},
			matchedPath: []string{"apps", "v1beta2", "deployment", "mixin", "metadata", "labels"},
			sr:          SearchResult{Setter: "withLabels"},
		},
	}

	astNode, err := ImportJsonnet("testdata/k8s.libsonnet")
	require.NoError(t, err)

	obj, ok := astNode.(*ast.Object)
	require.True(t, ok)

	node := NewNode("root", obj)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sr, matchedPath, err := node.Search(tc.path...)
			if tc.isErr {
				spew.Dump(sr)
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.matchedPath, matchedPath)
				assert.Equal(t, tc.sr, sr)
			}
		})
	}

}
