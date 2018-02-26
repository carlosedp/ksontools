local params = std.extVar("__ksonnet/params").components.certificateCrd;
local k = import "k.libsonnet";
local createCustomResourceDefinitionMetadata(labels, name) = k.apiextensions.v1beta1.customResourceDefinition.mixin.metadata.withLabels(labels).withName(name);
local createCustomResourceDefinitionSpec(group, scope, version) = k.apiextensions.v1beta1.customResourceDefinition.mixin.spec.withGroup(group).withScope(scope).withVersion(version);
local createCustomResourceDefinitionSpecNames(kind, plural) = k.apiextensions.v1beta1.customResourceDefinition.mixin.spec.names.withKind(kind).withPlural(plural);
local createCertificateCrd(params) =
  local customResourceDefinitionMetadata = createCustomResourceDefinitionMetadata(params.crdMetadataLabels, params.crdMetadataName);
  local customResourceDefinitionSpec = createCustomResourceDefinitionSpec(params.crdSpecGroup, params.crdSpecScope, params.crdSpecVersion);
  local customResourceDefinitionSpecNames = createCustomResourceDefinitionSpecNames(params.crdSpecNamesKind, params.crdSpecNamesPlural);

  k.apiextensions.v1beta1.customResourceDefinition.new() + customResourceDefinitionMetadata + customResourceDefinitionSpec + customResourceDefinitionSpecNames;
local certificateCrd = createCertificateCrd(params);

certificateCrd