package yaml2jsonnet

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLocals_Generate(t *testing.T) {
	entries := []LocalEntry{
		{Path: "apps.v1beta2.deployment.mixin.metadata", Setter: "withLabels"},
		{Path: "apps.v1beta2.deployment.mixin.metadata", Setter: "withName"},
		{Path: "apps.v1beta2.deployment.mixin.spec", Setter: "withReplicas"},
		{Path: "apps.v1beta2.deployment.mixin.spec.selector", Setter: "withMatchLabels"},
		{Path: "apps.v1beta2.deployment.mixin.spec.template.metadata", Setter: "withLabels"},
		{Path: "apps.v1beta2.deployment.mixin.spec.template.spec", Setter: "withContainers"},
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
			Value: NewDeclarationApply("deployment.mixin.metadata.withLabels(foo).withName(foo)"),
		},
		{
			Name:  "deploymentMixinSpec",
			Value: NewDeclarationApply("deployment.mixin.spec.withReplicas(foo)"),
		},
		{
			Name:  "deploymentMixinSpecSelector",
			Value: NewDeclarationApply("deploymentMixinSpec.selector.withMatchLabels(foo)"),
		},
		{
			Name:  "deploymentMixinSpecTemplateMetadata",
			Value: NewDeclarationApply("deploymentMixinSpec.template.metadata.withLabels(foo)"),
		},
		{
			Name:  "deploymentMixinSpecTemplateSpec",
			Value: NewDeclarationApply("deploymentMixinSpec.template.spec.withContainers(foo)"),
		},
	}

	require.Equal(t, expected, decls)

}
