// Package fragment combines the settings.json fragments shipped by sets into
// one settings object.
package fragment

import (
	"fmt"
	"sort"
	"strings"
)

// Fragment is one set's settings fragment.
type Fragment struct {
	Set  string         // set name, for errors
	Data map[string]any // decoded JSON object
}

// Merge combines fragments in order. Objects merge recursively and string
// lists merge as sets (union, first-seen order). Any other key set by two
// fragments is an error naming both sets and the key path, even when the
// values are equal. The inputs are not modified.
func Merge(frags []Fragment) (map[string]any, error) {
	out := map[string]any{}
	owners := map[string]string{} // key path -> set that first set it
	for _, f := range frags {
		if err := mergeObject(out, f.Data, "", f.Set, owners); err != nil {
			return nil, err
		}
	}
	return out, nil
}

func mergeObject(dst, src map[string]any, prefix, set string, owners map[string]string) error {
	keys := make([]string, 0, len(src))
	for k := range src {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		path := k
		if prefix != "" {
			path = prefix + "." + k
		}
		v := src[k]
		cur, exists := dst[k]
		if !exists {
			dst[k] = clone(v)
			owners[path] = set
			continue
		}
		if cm, ok := cur.(map[string]any); ok {
			if vm, ok := v.(map[string]any); ok {
				if err := mergeObject(cm, vm, path, set, owners); err != nil {
					return err
				}
				continue
			}
		}
		if cs, ok := stringList(cur); ok {
			if vs, ok := stringList(v); ok {
				dst[k] = union(cs, vs)
				continue
			}
		}
		return fmt.Errorf("fragment: key %s set by both %q and %q", path, owner(owners, path), set)
	}
	return nil
}

// owner finds the set that first set path, or the nearest object above it.
func owner(owners map[string]string, path string) string {
	for {
		if set, ok := owners[path]; ok {
			return set
		}
		i := strings.LastIndex(path, ".")
		if i < 0 {
			return ""
		}
		path = path[:i]
	}
}

func union(a, b []string) []any {
	seen := map[string]bool{}
	out := []any{}
	for _, list := range [][]string{a, b} {
		for _, s := range list {
			if !seen[s] {
				seen[s] = true
				out = append(out, s)
			}
		}
	}
	return out
}

func stringList(v any) ([]string, bool) {
	list, ok := v.([]any)
	if !ok {
		return nil, false
	}
	out := make([]string, 0, len(list))
	for _, x := range list {
		s, ok := x.(string)
		if !ok {
			return nil, false
		}
		out = append(out, s)
	}
	return out, true
}

// clone deep-copies decoded JSON so the result never aliases a fragment.
func clone(v any) any {
	switch v := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(v))
		for k, x := range v {
			out[k] = clone(x)
		}
		return out
	case []any:
		out := make([]any, len(v))
		for i, x := range v {
			out[i] = clone(x)
		}
		return out
	default:
		return v
	}
}
