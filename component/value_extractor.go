package component

import (
	"strings"

	"github.com/bryanl/woowoo/node"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/astext"
	"github.com/pkg/errors"
)

// Values are values extracted from a manifest.
type Values struct {
	Lookup []string
	Setter string
	Value  interface{}
}

// ValueExtractor extracts Values from a manifest.
type ValueExtractor struct {
	object *node.Node
}

// NewValueExtractor creates an instance of ValueExtractor.
func NewValueExtractor(root *astext.Object) *ValueExtractor {
	return &ValueExtractor{
		object: node.New("root", root),
	}
}

// Extract extracts values from an object.
func (ve *ValueExtractor) Extract(gvk GVK, props Properties) (map[string]Values, error) {
	m := make(map[string]Values)
	cache := make(map[string]bool)

	paths := props.Paths(gvk)
	for _, path := range paths {
		item, err := ve.object.Search2(path.Path...)
		if err != nil {
			continue
		}

		var manifestPath []string
		var found bool
		for _, p := range item.Path {
			if p == gvk.Kind {
				found = true
				continue
			}

			if !found {
				continue
			}

			manifestPath = append(manifestPath, p)
		}

		cachedPath := strings.Join(manifestPath, ".")
		if _, ok := cache[cachedPath]; ok {
			continue
		}

		cache[cachedPath] = true

		v, err := props.Value(manifestPath)
		if err != nil {
			return nil, errors.Wrapf(err, "retrieve values for %s", strings.Join(manifestPath, "."))
		}

		lookupPath := manifestPath
		if manifestPath[0] == "mixin" {
			lookupPath = manifestPath[1:]
		}

		p := strings.Join(item.Path, ".")
		dv := Values{
			Lookup: lookupPath,
			Setter: item.Name,
			Value:  v,
		}

		m[p] = dv
	}

	return m, nil
}
