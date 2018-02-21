package yaml2jsonnet

import (
	"sort"
	"strings"

	"github.com/pkg/errors"
)

func buildConstructors(m map[string]documentValues) (map[string][]string, error) {

	groups := make(map[string][]string)

	for _, dv := range m {
		ns, setter, err := parseSetterNamespace(dv.setter)
		if err != nil {
			return nil, errors.Wrap(err, "parse setter namespace")
		}

		if _, ok := groups[ns]; !ok {
			groups[ns] = make([]string, 0)
		}

		groups[ns] = append(groups[ns], setter)
	}

	for k := range groups {
		sort.Strings(groups[k])
	}

	return groups, nil
}

func parseSetterNamespace(s string) (string, string, error) {
	parts := strings.Split(s, ".")
	if len(parts) < 2 {
		return "", "", errors.Errorf("unexpected namespaced setter %q", s)
	}

	high := len(parts) - 1

	ns := strings.Join(parts[:high], ".")
	setter := parts[high]

	return ns, setter, nil
}
