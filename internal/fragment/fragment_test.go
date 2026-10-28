package fragment

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func decode(t *testing.T, s string) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("decode %s: %v", s, err)
	}
	return m
}

func TestMerge(t *testing.T) {
	type frag struct{ set, data string }
	tests := []struct {
		name    string
		frags   []frag
		want    string
		wantErr string
	}{
		{
			name: "no fragments",
			want: `{}`,
		},
		{
			name:  "single fragment",
			frags: []frag{{"a", `{"model":"opus","env":{"X":"1"}}`}},
			want:  `{"model":"opus","env":{"X":"1"}}`,
		},
		{
			name:  "nil data",
			frags: []frag{{"a", `null`}, {"b", `{"x":1}`}},
			want:  `{"x":1}`,
		},
		{
			name:  "disjoint keys",
			frags: []frag{{"a", `{"model":"opus"}`}, {"b", `{"theme":"dark"}`}},
			want:  `{"model":"opus","theme":"dark"}`,
		},
		{
			name: "nested objects merge recursively",
			frags: []frag{
				{"a", `{"permissions":{"defaultMode":"auto","env":{"A":"1"}}}`},
				{"b", `{"permissions":{"env":{"B":"2"}}}`},
			},
			want: `{"permissions":{"defaultMode":"auto","env":{"A":"1","B":"2"}}}`,
		},
		{
			name: "string lists union in first-seen order",
			frags: []frag{
				{"a", `{"permissions":{"allow":["Bash(ls)","Read"]}}`},
				{"b", `{"permissions":{"allow":["Read","Edit","Bash(ls)"]}}`},
				{"c", `{"permissions":{"allow":["Grep","Edit"]}}`},
			},
			want: `{"permissions":{"allow":["Bash(ls)","Read","Edit","Grep"]}}`,
		},
		{
			name:  "empty list unions with string list",
			frags: []frag{{"a", `{"l":[]}`}, {"b", `{"l":["x"]}`}},
			want:  `{"l":["x"]}`,
		},
		{
			name:  "duplicates inside one list collapse",
			frags: []frag{{"a", `{"l":["x","x"]}`}, {"b", `{"l":["y","x"]}`}},
			want:  `{"l":["x","y"]}`,
		},
		{
			name:    "scalar clash",
			frags:   []frag{{"a", `{"model":"opus"}`}, {"b", `{"model":"sonnet"}`}},
			wantErr: `key model set by both "a" and "b"`,
		},
		{
			name:    "equal scalars still clash",
			frags:   []frag{{"a", `{"model":"opus"}`}, {"b", `{"model":"opus"}`}},
			wantErr: `key model set by both "a" and "b"`,
		},
		{
			name: "non-string lists clash with nested path",
			frags: []frag{
				{"a", `{"hooks":{"Stop":[{"hooks":[{"command":"a.sh"}]}]}}`},
				{"b", `{"hooks":{"Stop":[{"hooks":[{"command":"b.sh"}]}]}}`},
			},
			wantErr: `key hooks.Stop set by both "a" and "b"`,
		},
		{
			name:    "object against scalar",
			frags:   []frag{{"a", `{"env":{"X":"1"}}`}, {"b", `{"env":"x"}`}},
			wantErr: `key env set by both "a" and "b"`,
		},
		{
			name:    "string list against scalar",
			frags:   []frag{{"a", `{"l":["x"]}`}, {"b", `{"l":"x"}`}},
			wantErr: `key l set by both "a" and "b"`,
		},
		{
			name: "clash names the set that first set the key",
			frags: []frag{
				{"a", `{"env":{"A":"1"}}`},
				{"b", `{"env":{"B":"2"}}`},
				{"c", `{"env":{"B":"3"}}`},
			},
			wantErr: `key env.B set by both "b" and "c"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var frags []Fragment
			for _, f := range tt.frags {
				frags = append(frags, Fragment{Set: f.set, Data: decode(t, f.data)})
			}
			got, err := Merge(frags)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("Merge() error = %v, want %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Merge() error = %v", err)
			}
			if want := decode(t, tt.want); !reflect.DeepEqual(got, want) {
				t.Errorf("Merge() = %v, want %v", got, want)
			}
		})
	}
}

func TestMergeDoesNotModifyInputs(t *testing.T) {
	a := decode(t, `{"env":{"A":"1"},"l":["x"]}`)
	b := decode(t, `{"env":{"B":"2"},"l":["y"]}`)
	got, err := Merge([]Fragment{{Set: "a", Data: a}, {Set: "b", Data: b}})
	if err != nil {
		t.Fatal(err)
	}
	got["env"].(map[string]any)["C"] = "3"
	if want := decode(t, `{"env":{"A":"1"},"l":["x"]}`); !reflect.DeepEqual(a, want) {
		t.Errorf("first fragment changed: %v", a)
	}
	if want := decode(t, `{"env":{"B":"2"},"l":["y"]}`); !reflect.DeepEqual(b, want) {
		t.Errorf("second fragment changed: %v", b)
	}
}
