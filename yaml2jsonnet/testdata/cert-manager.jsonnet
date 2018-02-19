local k = import "k.libsonnet";
local stdlib = import "stdlib.libsonnet";
local createCustomResourceDefinition(params) = {
};
local crd = createCustomResourceDefinition(params);

k.core.v1.list.new([abcCustomResourceDefinition])