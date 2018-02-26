package component

var (
	// TODO: might need something in ksonnet lib to look this up
	groupMappings = map[string][]string{
		"apiextensions.k8s.io":      []string{"apiextensions"},
		"rbac.authorization.k8s.io": []string{"rbac"},
	}
)

// GVK is a group, version, kind descriptor.
type GVK struct {
	GroupPath []string
	Version   string
	Kind      string
}

// Group returns the group this GVK represents.
func (gvk *GVK) Group() []string {
	g, ok := groupMappings[gvk.GroupPath[0]]
	if !ok {
		return gvk.GroupPath
	}

	return g
}

// Path returns the path of the current descriptor as a slice of strings.
func (gvk *GVK) Path() []string {
	return append(gvk.Group(), gvk.Version, gvk.Kind)
}
