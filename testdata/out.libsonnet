local k = import "k.libsonnet";
local deployment = "k.apps.v1beta2.deployment";

local deploymentInstance = deployment.new();

local chainedApply = a.b.c.withFoo(foo).withBar(bar);

{}