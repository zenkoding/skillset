package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCommands(t *testing.T) {
	root := t.TempDir()
	for p, body := range map[string]string{
		"sets/a/set.toml":  "name = \"a\"\nhosts = [\"claude\"]",
		"sets/a/shared.md": "a\n",
	} {
		path := filepath.Join(root, p)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	home := t.TempDir()
	tests := []struct {
		name  string
		args  []string
		out   string // substring of stdout
		error string // substring of the error; "" means success
	}{
		{"unknown set", []string{"plan", "--catalog", root, "--home", home, "--set", "nope"}, "", `unknown set "nope"`},
		{"plan", []string{"plan", "--catalog", root, "--home", home}, "1 to render, 1 to link", ""},
		{"apply", []string{"apply", "--catalog", root, "--home", home}, "installed", ""},
		{"plan after apply", []string{"plan", "--catalog", root, "--home", home, "--set", "a"}, "0 to render, 0 to link", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			cmd := newRootCmd(strings.NewReader(""), &out)
			cmd.SetArgs(tt.args)
			err := cmd.Execute()
			if tt.error != "" {
				if err == nil || !strings.Contains(err.Error(), tt.error) {
					t.Fatalf("error = %v, want %q", err, tt.error)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(out.String(), tt.out) {
				t.Fatalf("output = %s, want %q", out.String(), tt.out)
			}
		})
	}
	data, err := os.ReadFile(filepath.Join(home, ".claude/CLAUDE.md"))
	if err != nil || string(data) != "a\n" {
		t.Fatalf("CLAUDE.md = %q, %v", data, err)
	}
}

// run runs skillset with args and returns its output.
func run(t *testing.T, args ...string) (string, error) {
	t.Helper()
	var out bytes.Buffer
	cmd := newRootCmd(strings.NewReader(""), &out)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return out.String(), err
}

func TestInitThenApply(t *testing.T) {
	root, home := t.TempDir(), t.TempDir()
	for p, body := range map[string]string{
		"sets/base/set.toml":       "name = \"base\"\nhosts = [\"claude\", \"codex\"]",
		"sets/base/shared.md":      "base\n",
		"sets/dcg/set.toml":        "name = \"dcg\"\nhosts = [\"claude\", \"codex\"]\nrequires = [\"base\"]",
		"sets/dcg/shared.md":       "dcg\n",
		"sets/other/set.toml":      "name = \"other\"\nhosts = [\"claude\"]",
		"sets/other/shared.md":     "other\n",
		"sets/other/settings.json": "{}",
	} {
		writeFile(t, filepath.Join(root, p), body)
	}
	writeFile(t, filepath.Join(home, ".config/skillset/local/shared.md"), "mine\n")

	if _, err := run(t, "init", "--catalog", root, "--home", home); err == nil || !strings.Contains(err.Error(), "pass --host or --set") {
		t.Fatalf("init without a terminal = %v", err)
	}
	if _, err := run(t, "init", "--catalog", root, "--home", home, "--host", "vim", "--set", "dcg"); err == nil {
		t.Fatal("init accepted an unknown host")
	}
	out, err := run(t, "init", "--catalog", root, "--home", home, "--host", "claude,codex", "--set", "dcg")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "added required sets: base") {
		t.Fatalf("init output = %s", out)
	}
	cfg, _ := os.ReadFile(filepath.Join(home, ".config/skillset/config.toml"))
	if !strings.Contains(string(cfg), `sets = ["base", "dcg"]`) {
		t.Fatalf("config = %s", cfg)
	}
	if _, err := run(t, "apply", "--catalog", root, "--home", home); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		path string
		want string // "" means must not exist
	}{
		{".claude/CLAUDE.md", "base\ndcg\nmine\n"},
		{".codex/AGENTS.md", "base\ndcg\nmine\n"},
		{".claude/settings.json", ""}, // only "other" has settings
	}
	for _, tt := range tests {
		data, err := os.ReadFile(filepath.Join(home, tt.path))
		if tt.want == "" {
			if err == nil {
				t.Errorf("%s should not be installed", tt.path)
			}
			continue
		}
		if string(data) != tt.want {
			t.Errorf("%s = %q, %v; want %q", tt.path, data, err, tt.want)
		}
	}

	// A set the catalog no longer has points at the config.
	writeFile(t, filepath.Join(home, ".config/skillset/config.toml"), "sets = [\"gone\"]\n")
	if _, err := run(t, "plan", "--catalog", root, "--home", home); err == nil || !strings.Contains(err.Error(), "config.toml; run skillset init") {
		t.Fatalf("plan with a stale set = %v", err)
	}
	if _, err := run(t, "init", "--catalog", root, "--home", home, "--host", "claude,codex", "--set", "dcg"); err != nil {
		t.Fatal(err)
	}

	// Dropping a host removes its links; --host alone keeps the saved sets.
	if _, err := run(t, "init", "--catalog", root, "--home", home, "--host", "claude"); err != nil {
		t.Fatal(err)
	}
	cfg, _ = os.ReadFile(filepath.Join(home, ".config/skillset/config.toml"))
	if !strings.Contains(string(cfg), `hosts = ["claude"]`) || !strings.Contains(string(cfg), `sets = ["base", "dcg"]`) {
		t.Fatalf("config after --host only = %s", cfg)
	}
	if _, err := run(t, "apply", "--catalog", root, "--home", home); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(filepath.Join(home, ".codex/AGENTS.md")); err == nil {
		t.Error(".codex/AGENTS.md still installed after dropping codex")
	}
}

func writeFile(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestUninstall(t *testing.T) {
	root, home := t.TempDir(), t.TempDir()
	writeFile(t, filepath.Join(root, "sets/a/set.toml"), "name = \"a\"\nhosts = [\"claude\"]")
	writeFile(t, filepath.Join(root, "sets/a/shared.md"), "a\n")
	writeFile(t, filepath.Join(root, "sets/a/skills/x/SKILL.md"), "---\nname: x\n---\nbody\n")
	if _, err := run(t, "apply", "--catalog", root, "--home", home); err != nil {
		t.Fatal(err)
	}
	paths := []string{".claude/CLAUDE.md", ".claude/skills/x"}
	out, err := run(t, "uninstall", "--dry-run", "--catalog", root, "--home", home)
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range paths {
		if !strings.Contains(out, "would remove "+filepath.Join(home, p)) {
			t.Errorf("dry run output misses %s:\n%s", p, out)
		}
		if _, err := os.Lstat(filepath.Join(home, p)); err != nil {
			t.Errorf("dry run removed %s", p)
		}
	}
	if _, err := run(t, "uninstall", "--catalog", root, "--home", home); err != nil {
		t.Fatal(err)
	}
	for _, p := range paths {
		if _, err := os.Lstat(filepath.Join(home, p)); err == nil {
			t.Errorf("%s still installed", p)
		}
	}
}
