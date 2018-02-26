package yaml2jsonnet

import (
	"github.com/pkg/errors"
)

var (
	// ErrNotFound is a not found error.
	ErrNotFound = errors.New("not found")
)

var (
	ignoredProps = []string{"mixin", "kind", "new", "mixinInstance"}
)

func stringInSlice(s string, sl []string) bool {
	for i := range sl {
		if sl[i] == s {
			return true
		}
	}

	return false
}
