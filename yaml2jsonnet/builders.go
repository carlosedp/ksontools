package yaml2jsonnet

import (
	"sort"
	"strings"

	"github.com/pkg/errors"
)

type ctorArgument struct {
	setter     string
	paramName  string
	paramValue interface{}
}

func buildConstructors(m map[string]documentValues) (map[string][]ctorArgument, error) {

	groups := make(map[string][]ctorArgument)

	for paramPath, dv := range m {
		ns, setter, err := parseSetterNamespace(dv.setter)
		if err != nil {
			return nil, errors.Wrap(err, "parse setter namespace")
		}

		if _, ok := groups[ns]; !ok {
			groups[ns] = make([]ctorArgument, 0)
		}

		ca := ctorArgument{
			setter:     setter,
			paramName:  paramName(paramPath),
			paramValue: dv.value,
		}

		groups[ns] = append(groups[ns], ca)
	}

	for k := range groups {
		sort.Slice(groups[k], func(i, j int) bool {
			return groups[k][i].setter < groups[k][j].setter
		})
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
