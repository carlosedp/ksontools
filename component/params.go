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

type patchDoc struct {
	Components map[string]interface{} `json:"components"`
}

func patchJSON(jsonObject, patch, patchName string) (string, error) {
	vm := jsonnet.MakeVM()
	vm.TLACode("target", jsonObject)
	vm.TLACode("patch", patch)
	vm.TLAVar("patchName", patchName)

	return vm.EvaluateSnippet("snippet", snippetMergeComponentPatch)
}

var snippetMergeComponentPatch = `
function(target, patch, patchName)
	if std.objectHas(patch.components, patchName) then
		std.mergePatch(target, patch.components[patchName])
	else
		target
`
