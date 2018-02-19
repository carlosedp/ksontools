package yaml2jsonnet

import (
	"fmt"
	"sort"
	"strings"
)

// LocalEntry is a local entry.
type LocalEntry struct {
	Path      string
	Setter    string
	ParamName string
	Arguments []string
}

// Locals are local definitions for jsonnet.
type Locals struct {
	Name    string
	Entries []LocalEntry
}

// NewLocals creates an instance of Locals.
func NewLocals(name string) *Locals {
	return &Locals{
		Name:    name,
		Entries: make([]LocalEntry, 0),
	}
}

// Add adds an entry to Locals.
func (l *Locals) Add(entry LocalEntry) {
	l.Entries = append(l.Entries, entry)
}

// Generate generates declarations for locals.
func (l *Locals) Generate() ([]Declaration, error) {
	candidates := make(map[string][]string)

	for _, entry := range l.Entries {
		parts := strings.Split(entry.Path, ".")

		for i := range parts {
			// find index which equals Name
			if parts[i] == l.Name {
				s := strings.Join(parts[i:], ".")
				candidates[s] = append(candidates[s],
					fmt.Sprintf("%s(params.%s)", entry.Setter, entry.ParamName))
				break
			}

		}
	}

	idMap := make(map[string]string)

	var candidateNames []string
	for cn := range candidates {
		candidateNames = append(candidateNames, cn)
	}

	sort.Strings(candidateNames)

	var decls []Declaration
	for _, applyBase := range candidateNames {
		setters := candidates[applyBase]
		apply := fmt.Sprintf("%s", applyBase)

		var mapIDs []string
		for id := range idMap {
			mapIDs = append(mapIDs, id)
		}

		sort.Strings(mapIDs)

		for _, id := range mapIDs {
			if strings.HasPrefix(applyBase, id) {
				apply = fmt.Sprintf("%s%s", idMap[id], strings.TrimPrefix(applyBase, id))
			}
		}

		for _, setter := range setters {
			apply = fmt.Sprintf("%s.%s", apply, setter)
		}

		parts := strings.Split(applyBase, ".")
		for i := range parts {
			if i > 0 {
				parts[i] = strings.Title(parts[i])
			}
		}

		name := strings.Join(parts, "")
		idMap[applyBase] = name

		decl := Declaration{
			Name:  name,
			Value: NewDeclarationApply(apply),
		}
		decls = append(decls, decl)
	}

	return decls, nil
}
