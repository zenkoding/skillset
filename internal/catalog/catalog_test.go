package catalog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zenkoding/skillset/internal/install"
)

func write(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// catalogOf writes files (path -> body) under a temp catalog root.
func catalogOf(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for p, body := range files {
		write(t, filepath.Join(root, p), body)
	}
	return root
}

func TestLoadErrors(t *testing.T) {
	tests := []struct {
		name  string
		files map[string]string
		want  string
	}{
		{"name differs from dir", map[string]string{"sets/a/set.toml": `name = "b"`}, `name "b" differs`},
		{"missing name", map[string]string{"sets/a/set.toml": `hosts = ["claude"]`}, "missing name"},
		{"unknown key", map[string]string{"sets/a/set.toml": "name = \"a\"\ncolor = \"red\""}, "unknown keys"},
		{"unknown host", map[string]string{"sets/a/set.toml": "name = \"a\"\nhosts = [\"vim\"]"}, `unknown host "vim"`},
		{"unknown requirement", map[string]string{"sets/a/set.toml": "name = \"a\"\nrequires = [\"x\"]"}, `requires unknown set "x"`},
		{"missing explicit path", map[string]string{"sets/a/set.toml": "name = \"a\"\nshared = \"nope.md\""}, "nope.md does not exist"},
		{"path leaves catalog", map[string]string{"sets/a/set.toml": "name = \"a\"\nskills = [\"../../..\"]"}, "leaves the catalog"},
		{"instructions for unlisted host", map[string]string{"sets/a/set.toml": "name = \"a\"\nhosts = [\"claude\"]\n[instructions]\ncodex = \"c.md\""}, `instructions for "codex"`},
		{"file without home", map[string]string{"sets/a/set.toml": "name = \"a\"\n[[files]]\npath = \"x\""}, "needs home and path"},
		{"absolute home", map[string]string{"sets/a/set.toml": "name = \"a\"\n[[files]]\nhome = \"/etc/passwd\"\npath = \"x\"", "sets/a/x": ""}, "must be a relative path under"},
		{"home escapes", map[string]string{"sets/a/set.toml": "name = \"a\"\n[[files]]\nhome = \"../x\"\npath = \"x\"", "sets/a/x": ""}, "must be a relative path under"},
		{"home escapes later", map[string]string{"sets/a/set.toml": "name = \"a\"\n[[files]]\nhome = \".claude/../../x\"\npath = \"x\"", "sets/a/x": ""}, "must be a relative path under"},
		{"home outside agent dirs", map[string]string{"sets/a/set.toml": "name = \"a\"\n[[files]]\nhome = \".ssh/authorized_keys\"\npath = \"x\"", "sets/a/x": ""}, "must be a relative path under"},
		{"home is an agent dir", map[string]string{"sets/a/set.toml": "name = \"a\"\n[[files]]\nhome = \".claude\"\npath = \"x\"", "sets/a/x": ""}, "must be a relative path under"},
		{"symlink leaves catalog", map[string]string{"sets/a/set.toml": "name = \"a\"\nshared = \"out.md\""}, "through a symlink"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := catalogOf(t, tt.files)
			if tt.name == "symlink leaves catalog" {
				outside := filepath.Join(t.TempDir(), "secret")
				write(t, outside, "")
				if err := os.Symlink(outside, filepath.Join(root, "sets/a/out.md")); err != nil {
					t.Fatal(err)
				}
			}
			_, err := Load(root)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("Load error = %v, want %q", err, tt.want)
			}
		})
	}
}

func TestLoadDefaults(t *testing.T) {
	root := catalogOf(t, map[string]string{
		"sets/a/set.toml":                 "name = \"a\"\nhosts = [\"claude\", \"codex\"]",
		"sets/a/shared.md":                "a shared\n",
		"sets/a/claude.md":                "a claude\n",
		"sets/a/skills/x/SKILL.md":        "---\nname: x\n---\n",
		"sets/a/settings.json":            "{}",
		"sets/bare/set.toml":              `name = "bare"`,
		"sets/not-a-set.txt":              "ignored",
		"sets/bare/unrelated/ignored.txt": "",
	})
	c, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(c.Names(), ","); got != "a,bare" {
		t.Fatalf("Names = %s", got)
	}
	a := c.Sets["a"]
	dir := filepath.Join(root, "sets/a")
	tests := []struct {
		name, got, want string
	}{
		{"shared", a.Shared, filepath.Join(dir, "shared.md")},
		{"claude", a.Instructions["claude"], filepath.Join(dir, "claude.md")},
		{"codex missing", a.Instructions["codex"], ""},
		{"skills", strings.Join(a.Skills, ","), filepath.Join(dir, "skills")},
		{"settings", a.Settings, filepath.Join(dir, "settings.json")},
		{"bare settings", c.Sets["bare"].Settings, ""},
		{"bare skills", strings.Join(c.Sets["bare"].Skills, ","), ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("got %q, want %q", tt.got, tt.want)
			}
		})
	}
}

func TestResolve(t *testing.T) {
	c := &Catalog{Sets: map[string]*Set{
		"base":  {Name: "base"},
		"dcg":   {Name: "dcg", Requires: []string{"base"}},
		"tools": {Name: "tools", Requires: []string{"dcg", "base"}},
		"loopa": {Name: "loopa", Requires: []string{"loopb"}},
		"loopb": {Name: "loopb", Requires: []string{"loopa"}},
	}}
	tests := []struct {
		name  string
		pick  []string
		want  string
		error string
	}{
		{"requires come first", []string{"tools"}, "base,dcg,tools", ""},
		{"no duplicates", []string{"dcg", "tools", "base"}, "base,dcg,tools", ""},
		{"unknown", []string{"nope"}, "", `unknown set "nope"`},
		{"cycle", []string{"loopa"}, "", "requires cycle: loopa -> loopb -> loopa"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sets, err := c.Resolve(tt.pick)
			if tt.error != "" {
				if err == nil || !strings.Contains(err.Error(), tt.error) {
					t.Fatalf("error = %v, want %q", err, tt.error)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			var names []string
			for _, s := range sets {
				names = append(names, s.Name)
			}
			if got := strings.Join(names, ","); got != tt.want {
				t.Fatalf("order = %s, want %s", got, tt.want)
			}
		})
	}
}

// TestInstall installs two sets and checks the per-host instructions, skills
// and files.
func TestInstall(t *testing.T) {
	root := catalogOf(t, map[string]string{
		"sets/base/set.toml":          "name = \"base\"\nhosts = [\"claude\", \"codex\"]",
		"sets/base/shared.md":         "base\n",
		"sets/base/codex.md":          "base codex\n",
		"sets/base/skills/x/SKILL.md": "---\nname: x\n---\n",
		"sets/cl/set.toml":            "name = \"cl\"\nhosts = [\"claude\"]\nrequires = [\"base\"]\n[[files]]\nhome = \".claude/line.sh\"\npath = \"line.sh\"",
		"sets/cl/shared.md":           "cl\n",
		"sets/cl/claude.md":           "cl claude\n",
		"sets/cl/line.sh":             "#!/bin/sh\n",
		"sets/cl/settings.json":       `{"a": 1}`,
	})
	home := t.TempDir()
	c, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	sets, err := c.Resolve([]string{"cl"})
	if err != nil {
		t.Fatal(err)
	}
	src, err := c.Source(sets, Hosts)
	if err != nil {
		t.Fatal(err)
	}
	p, err := install.MakeSourcePlan(src, home)
	if err != nil {
		t.Fatal(err)
	}
	if err := p.Apply(); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		path, want string
	}{
		{".claude/CLAUDE.md", "base\ncl\ncl claude\n"},
		{".codex/AGENTS.md", "base\nbase codex\n"},
		{".claude/line.sh", "#!/bin/sh\n"},
		{".claude/settings.json", `{"a": 1}`},
		{".agents/skills/x/SKILL.md", "---\nname: x\n---\n"},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join(home, tt.path))
			if err != nil {
				t.Fatal(err)
			}
			if string(data) != tt.want {
				t.Fatalf("%s = %q, want %q", tt.path, data, tt.want)
			}
		})
	}
}

func TestSourceSettings(t *testing.T) {
	tests := []struct {
		name    string
		a, b    string // settings.json of sets a and b; "" means none
		content string // merged content; "" means read from the file
		error   string
	}{
		{"one fragment is copied as is", `{"x": 1}`, "", "", ""},
		{"fragments merge", `{"x": 1, "allow": ["a"]}`, `{"y": 2.5, "allow": ["b", "a"]}`, "{\n  \"allow\": [\n    \"a\",\n    \"b\"\n  ],\n  \"x\": 1,\n  \"y\": 2.5\n}\n", ""},
		{"same key fails", `{"x": 1}`, `{"x": 1}`, "", `"a" and "b"`},
		{"bad JSON", `{`, "", "", "set a:"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			var sets []*Set
			for _, f := range []struct{ name, body string }{{"a", tt.a}, {"b", tt.b}} {
				s := &Set{Name: f.name, Hosts: []string{"claude"}, Dir: filepath.Join(dir, f.name)}
				if f.body != "" {
					s.Settings = filepath.Join(dir, f.name+".json")
					write(t, s.Settings, f.body)
				}
				sets = append(sets, s)
			}
			src, err := (&Catalog{Root: dir}).Source(sets, Hosts)
			if tt.error != "" {
				if err == nil || !strings.Contains(err.Error(), tt.error) {
					t.Fatalf("error = %v, want %q", err, tt.error)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if len(src.Copies) != 1 {
				t.Fatalf("copies = %+v, want one", src.Copies)
			}
			if got := string(src.Copies[0].Content); got != tt.content {
				t.Fatalf("content = %q, want %q", got, tt.content)
			}
		})
	}
}

// TestRepoCatalog checks that the legacy set in config/ installs cleanly
// into an empty home: both instruction files, settings and every file it
// lists.
func TestRepoCatalog(t *testing.T) {
	root, err := filepath.Abs("../../config")
	if err != nil {
		t.Fatal(err)
	}
	c, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	sets, err := c.Resolve([]string{"legacy"})
	if err != nil {
		t.Fatal(err)
	}
	src, err := c.Source(sets, Hosts)
	if err != nil {
		t.Fatal(err)
	}
	home := t.TempDir()
	p, err := install.MakeSourcePlan(src, home)
	if err != nil {
		t.Fatal(err)
	}
	if err := p.Err(); err != nil {
		t.Fatal(err)
	}
	links := map[string]bool{}
	for _, l := range p.Links {
		links[l.Path] = true
	}
	want := []string{".claude/CLAUDE.md", ".codex/AGENTS.md", ".claude/skills/i-have-adhd", ".agents/skills/i-have-adhd"}
	for _, f := range sets[0].Files {
		want = append(want, f.Home)
	}
	for _, w := range want {
		if !links[filepath.Join(home, w)] {
			t.Errorf("no link for %s", w)
		}
	}
	if len(p.Copies) != 1 || p.Copies[0].Path != filepath.Join(home, ".claude/settings.json") {
		t.Errorf("copies = %+v, want .claude/settings.json", p.Copies)
	}
}
