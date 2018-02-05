package yaml2jsonnet

import (
	"fmt"
	"sort"
)

type Properties map[interface{}]interface{}

func (p Properties) Paths(gvk GVK) []PropertyPath {
	ch := make(chan PropertyPath)

	go func() {
		base := []string{gvk.Group, gvk.Version, gvk.Kind}
		iterateMap(ch, base, p)
		close(ch)
	}()

	var out []PropertyPath
	for pr := range ch {
		out = append(out, pr)
	}

	return out
}

type PropertyPath struct {
	Path []string
}

func iterateMap(ch chan PropertyPath, base []string, m map[interface{}]interface{}) {
	localBase := make([]string, len(base))
	copy(localBase, base)

	var keys []interface{}
	for k := range m {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		a := keys[i].(string)
		b := keys[j].(string)

		return a < b
	})

	for i := range keys {
		name := keys[i].(string)
		switch t := m[name].(type) {
		default:
			panic(fmt.Sprintf("not sure what to do with %T", t))
		case map[interface{}]interface{}:
			newBase := append(localBase, name)
			iterateMap(ch, newBase, t)
		case string, int, []interface{}:
			ch <- PropertyPath{
				Path: append(base, name),
			}

		}
	}

}

func intersection(a []string, b []string) [][]string {
	if isStringSliceEqual(a, b) {
		return [][]string{a}
	}

	var inter []string

	low, high := a, b
	if len(a) > len(b) {
		low = b
		high = a
	}

	done := false
	for i, l := range low {
		for j, h := range high {
			f1 := i + 1
			f2 := j + 1
			if l == h {
				inter = append(inter, h)
				if f1 < len(low) && f2 < len(high) {
					if low[f1] != high[f2] {
						done = true
					}
				}
				high = high[:j+copy(high[j:], high[j+1:])]
				break
			}
		}
		if done {
			break
		}
	}

	if isStringSliceEqual(a, inter) {
		return [][]string{a}
	} else if isStringSliceEqual(b, inter) {
		return [][]string{b}
	}

	return [][]string{a, b}
}

func isStringSliceEqual(a, b []string) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
