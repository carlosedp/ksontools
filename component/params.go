package component

import (
	jsonnet "github.com/google/go-jsonnet"
)

func applyGlobals(params string) (string, error) {
	vm := jsonnet.MakeVM()

	vm.ExtCode("params", params)
	return vm.EvaluateSnippet("snippet", snippetMapGlobal)
}

var snippetMapGlobal = `
local params = std.extVar("params");
local applyGlobal = function(key, value) std.mergePatch(value, params.global);

{
	components: std.mapWithKey(applyGlobal, params.components)
}
`
