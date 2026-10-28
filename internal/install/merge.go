package install

import (
	"bytes"
	"encoding/json"
	"reflect"
	"sort"
	"strings"
)

// missing stands for a key that is missing on one side of a merge.
type missing struct{}

// Conflict is one settings key both the repo and the local copy changed, to
// different values, since the last install.
type Conflict struct {
	Key   string // JSON pointer, e.g. /permissions/defaultMode
	Repo  string // the repo value as JSON, or "(absent)"
	Local string // the local value as JSON, or "(absent)"
}

// mergeJSON reconciles a copied JSON file. repoBase and installed are what
// the repo held and what the installer wrote last time (missing{} when
// unknown). Keys only one side changed take that side's change; string lists
// merge as sets; anything else both sides changed is a conflict, settled by
// choices (key -> keep local) or left at the repo value.
func mergeJSON(repoBase, repo, installed, local any, choices map[string]bool) (any, []Conflict) {
	var conflicts []Conflict
	out := merge("", repoBase, repo, installed, local, choices, &conflicts)
	return out, conflicts
}

func merge(key string, repoBase, repo, installed, local any, choices map[string]bool, conflicts *[]Conflict) any {
	switch {
	case reflect.DeepEqual(repoBase, repo), reflect.DeepEqual(repo, local):
		return local
	case reflect.DeepEqual(installed, local):
		return repo
	}
	rm, rok := repo.(map[string]any)
	lm, lok := local.(map[string]any)
	if rok && lok {
		bm, im := asMap(repoBase), asMap(installed)
		out := map[string]any{}
		for _, k := range unionKeys(rm, lm, bm, im) {
			v := merge(key+"/"+k, get(bm, k), get(rm, k), get(im, k), get(lm, k), choices, conflicts)
			if _, gone := v.(missing); !gone {
				out[k] = v
			}
		}
		return out
	}
	if rs, ok := stringList(repo); ok {
		if ls, ok := stringList(local); ok {
			bs, _ := stringList(repoBase)
			return mergeSets(bs, rs, ls)
		}
	}
	*conflicts = append(*conflicts, Conflict{Key: key, Repo: show(repo), Local: show(local)})
	if choices[key] {
		return local
	}
	return repo
}

// mergeSets keeps local's list, adds what the repo added and drops what the
// repo removed since repoBase.
func mergeSets(base, repo, local []string) []any {
	in := func(list []string, s string) bool {
		for _, x := range list {
			if x == s {
				return true
			}
		}
		return false
	}
	out := []any{}
	for _, s := range local {
		if in(base, s) && !in(repo, s) {
			continue
		}
		out = append(out, s)
	}
	for _, s := range repo {
		if !in(base, s) && !in(local, s) {
			out = append(out, s)
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

func asMap(v any) map[string]any {
	m, _ := v.(map[string]any)
	return m
}

func get(m map[string]any, k string) any {
	if v, ok := m[k]; ok {
		return v
	}
	return missing{}
}

func unionKeys(maps ...map[string]any) []string {
	seen := map[string]bool{}
	var keys []string
	for _, m := range maps {
		for k := range m {
			if !seen[k] {
				seen[k] = true
				keys = append(keys, k)
			}
		}
	}
	sort.Strings(keys)
	return keys
}

func show(v any) string {
	if _, ok := v.(missing); ok {
		return "(absent)"
	}
	data, err := encodeJSON(v)
	if err != nil {
		return "?"
	}
	return strings.Join(strings.Fields(string(data)), " ")
}

// decodeJSON parses data keeping numbers exact; empty or missing data is absent.
func decodeJSON(data []byte) (any, error) {
	if data == nil {
		return missing{}, nil
	}
	var v any
	d := json.NewDecoder(bytes.NewReader(data))
	d.UseNumber()
	if err := d.Decode(&v); err != nil {
		return nil, err
	}
	return v, nil
}

func encodeJSON(v any) ([]byte, error) {
	var b bytes.Buffer
	e := json.NewEncoder(&b)
	e.SetEscapeHTML(false)
	e.SetIndent("", "  ")
	if err := e.Encode(v); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
