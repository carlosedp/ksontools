package yaml2jsonnet

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLocals_Generate(t *testing.T) {
	entries := []LocalEntry{
		{Path: "apps.v1beta2.deployment.mixin.metadata", Setter: "withLabels", ParamName: "labels"},
		{Path: "apps.v1beta2.deployment.mixin.metadata", Setter: "withName", ParamName: "name"},
		{Path: "apps.v1beta2.deployment.mixin.spec", Setter: "withReplicas", ParamName: "replicas"},
		{Path: "apps.v1beta2.deployment.mixin.spec.selector", Setter: "withMatchLabels", ParamName: "labels"},
		{Path: "apps.v1beta2.deployment.mixin.spec.template.metadata", Setter: "withLabels", ParamName: "labels"},
		{Path: "apps.v1beta2.deployment.mixin.spec.template.spec", Setter: "withContainers", ParamName: "containers"},
	}

	l := NewLocals("deployment")

	for _, entry := range entries {
		l.Add(entry)
	}

	decls, err := l.Generate()
	require.NoError(t, err)

	expected := []Declaration{
		{
			Name:  "deploymentMixinMetadata",
			Value: NewDeclarationApply("deployment.mixin.metadata.withLabels(params.labels).withName(params.name)"),
		},
		{
			Name:  "deploymentMixinSpec",
			Value: NewDeclarationApply("deployment.mixin.spec.withReplicas(params.replicas)"),
		},
		{
			Name:  "deploymentMixinSpecSelector",
			Value: NewDeclarationApply("deploymentMixinSpec.selector.withMatchLabels(params.labels)"),
		},
		{
			Name:  "deploymentMixinSpecTemplateMetadata",
			Value: NewDeclarationApply("deploymentMixinSpec.template.metadata.withLabels(params.labels)"),
		},
		{
			Name:  "deploymentMixinSpecTemplateSpec",
			Value: NewDeclarationApply("deploymentMixinSpec.template.spec.withContainers(params.containers)"),
		},
	}

	require.Equal(t, expected, decls)

}
