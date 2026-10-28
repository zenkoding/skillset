package catalog

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/zenkoding/skillset/internal/install"
)

func TestConfigRoundTrip(t *testing.T) {
	home := t.TempDir()
	if _, ok, err := ReadConfig(home); ok || err != nil {
		t.Fatalf("ReadConfig on empty home = %v, %v", ok, err)
	}
	want := Config{Hosts: []string{"claude", "codex"}, Sets: []string{"base", "dcg"}}
	if err := WriteConfig(home, Config{Hosts: []string{"codex", "claude"}, Sets: []string{"dcg", "base"}}); err != nil {
		t.Fatal(err)
	}
	got, ok, err := ReadConfig(home)
	if err != nil || !ok || !reflect.DeepEqual(got, want) {
		t.Fatalf("ReadConfig = %+v, %v, %v; want %+v", got, ok, err, want)
	}
	write(t, ConfigPath(home), "sets = [\"a\"]\ncolour = 1\n")
	if _, _, err := ReadConfig(home); err == nil || !strings.Contains(err.Error(), "unknown keys") {
		t.Fatalf("ReadConfig with unknown key = %v", err)
	}
}

func TestLoadLocal(t *testing.T) {
	c := &Catalog{Sets: map[string]*Set{"base": {Name: "base"}}}
	tests := []struct {
		name  string
		files map[string]string // under HOME
		hosts string            // "" when no set is expected
		error string
	}{
		{"no dir", nil, "", ""},
		{"no manifest", map[string]string{LocalDir + "/shared.md": "mine\n"}, "claude,codex", ""},
		{"manifest", map[string]string{LocalDir + "/set.toml": "name = \"local\"\nhosts = [\"codex\"]\nrequires = [\"base\"]"}, "codex", ""},
		{"wrong name", map[string]string{LocalDir + "/set.toml": `name = "mine"`}, "", `name must be "local"`},
		{"unknown requirement", map[string]string{LocalDir + "/set.toml": "name = \"local\"\nrequires = [\"x\"]"}, "", `unknown set "x"`},
		{"path leaves local", map[string]string{LocalDir + "/set.toml": "name = \"local\"\nshared = \"../config.toml\""}, "", "leaves the catalog"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			home := catalogOf(t, tt.files)
			s, err := c.LoadLocal(home)
			if tt.error != "" {
				if err == nil || !strings.Contains(err.Error(), tt.error) {
					t.Fatalf("error = %v, want %q", err, tt.error)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if tt.hosts == "" {
				if s != nil {
					t.Fatalf("got set %+v, want none", s)
				}
				return
			}
			if got := strings.Join(s.Hosts, ","); got != tt.hosts {
				t.Fatalf("hosts = %s, want %s", got, tt.hosts)
			}
		})
	}
}

func TestSelect(t *testing.T) {
	c := &Catalog{Sets: map[string]*Set{
		"base": {Name: "base"},
		"dcg":  {Name: "dcg", Requires: []string{"base"}},
		"misc": {Name: "misc"},
	}}
	local := &Set{Name: LocalName, Requires: []string{"misc"}}
	tests := []struct {
		name  string
		names []string
		hosts []string
		local *Set
		sets  string
		host  string
	}{
		{"defaults to everything", nil, nil, nil, "base,dcg,misc", "claude,codex"},
		{"requires added", []string{"dcg"}, []string{"codex"}, nil, "base,dcg", "codex"},
		{"local last with its requires", []string{"dcg"}, nil, local, "base,dcg,misc,local", "claude,codex"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sel, err := c.Select(tt.names, tt.hosts, tt.local)
			if err != nil {
				t.Fatal(err)
			}
			var names []string
			for _, s := range sel.Sets {
				names = append(names, s.Name)
			}
			if got := strings.Join(names, ","); got != tt.sets {
				t.Fatalf("sets = %s, want %s", got, tt.sets)
			}
			if got := strings.Join(sel.Hosts, ","); got != tt.host {
				t.Fatalf("hosts = %s, want %s", got, tt.host)
			}
		})
	}
}

// TestHostsAndLocal installs for Codex only, with a local set, and checks
// that nothing Claude-only lands and the local part comes last.
func TestHostsAndLocal(t *testing.T) {
	root := catalogOf(t, map[string]string{
		"sets/base/set.toml":          "name = \"base\"\nhosts = [\"claude\", \"codex\"]\n[[files]]\nhome = \".claude/line.sh\"\npath = \"line.sh\"\nhosts = [\"claude\"]\n[[files]]\nhome = \".config/tool.toml\"\npath = \"tool.toml\"",
		"sets/base/shared.md":         "base\n",
		"sets/base/line.sh":           "",
		"sets/base/tool.toml":         "",
		"sets/base/settings.json":     "{}",
		"sets/base/skills/x/SKILL.md": "---\nname: x\n---\n",
	})
	home := catalogOf(t, map[string]string{
		LocalDir + "/shared.md":         "mine\n",
		LocalDir + "/skills/y/SKILL.md": "---\nname: y\n---\n",
	})
	c, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	local, err := c.LoadLocal(home)
	if err != nil {
		t.Fatal(err)
	}
	sel, err := c.Select(nil, []string{"codex"}, local)
	if err != nil {
		t.Fatal(err)
	}
	src, err := c.Source(sel.Sets, sel.Hosts)
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
		path string
		want string // "" means must not exist
	}{
		{".codex/AGENTS.md", "base\nmine\n"},
		{".codex/skills/x/SKILL.md", "---\nname: x\n---\n"},
		{".agents/skills/y/SKILL.md", "---\nname: y\n---\n"},
		{".config/tool.toml", "\x00"},
		{".claude/CLAUDE.md", ""},
		{".claude/skills/x", ""},
		{".claude/line.sh", ""},
		{".claude/settings.json", ""},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join(home, tt.path))
			switch {
			case tt.want == "":
				if err == nil {
					t.Fatalf("%s should not be installed", tt.path)
				}
			case err != nil:
				t.Fatal(err)
			case tt.want != "\x00" && string(data) != tt.want:
				t.Fatalf("%s = %q, want %q", tt.path, data, tt.want)
			}
		})
	}
	// A second plan sees the local links as ours: nothing to do.
	p, err = install.MakeSourcePlan(src, home)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(p.String(), "0 to link") || len(p.Remove)+len(p.Kept)+len(p.Clashes) > 0 {
		t.Fatalf("second plan not clean:\n%s", p)
	}
}
