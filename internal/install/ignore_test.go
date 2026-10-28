package install

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIgnore(t *testing.T) {
	tests := []struct {
		name      string
		gitconfig string // ~/.gitconfig, "" for none
		file      string // excludes file, relative to HOME
		before    string // its content before install, "" for no file
		want      string
		exists    bool // already listed before install
	}{
		{name: "default location, no file", file: ".config/git/ignore", want: ".docs/\n"},
		{name: "default location, no trailing newline", file: ".config/git/ignore", before: "*.swp", want: "*.swp\n.docs/\n"},
		{
			name:      "configured file",
			gitconfig: "[core]\n\texcludesFile = ~/.gitignore_global\n",
			file:      ".gitignore_global",
			before:    ".DS_Store\n",
			want:      ".DS_Store\n.docs/\n",
		},
		{name: "already listed", file: ".config/git/ignore", before: "  .docs\n", want: "  .docs\n", exists: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, home := fixture(t)
			if tt.gitconfig != "" {
				write(t, filepath.Join(home, ".gitconfig"), tt.gitconfig)
			}
			path := filepath.Join(home, tt.file)
			if tt.before != "" {
				write(t, path, tt.before)
			}
			p, err := makePlan(repo, home)
			if err != nil {
				t.Fatal(err)
			}
			if p.Ignore.Path != path || p.Ignore.Exists != tt.exists {
				t.Fatalf("Ignore = %+v, want path %s, exists %v", p.Ignore, path, tt.exists)
			}
			if !strings.Contains(p.String(), path) {
				t.Fatalf("plan does not mention %s:\n%s", path, p)
			}
			// A second install must not add the pattern again.
			install(t, repo, home)
			install(t, repo, home)
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if string(data) != tt.want {
				t.Fatalf("%s = %q, want %q", tt.file, data, tt.want)
			}
		})
	}
}

func TestIgnoreKeepsSymlink(t *testing.T) {
	repo, home := fixture(t)
	real := filepath.Join(home, "dotfiles/ignore")
	write(t, real, "*.log\n")
	link := filepath.Join(home, ".config/git/ignore")
	if err := os.MkdirAll(filepath.Dir(link), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(real, link); err != nil {
		t.Fatal(err)
	}
	install(t, repo, home)
	if !linksTo(link, real) {
		t.Fatalf("%s is no longer a link to %s", link, real)
	}
	if data, _ := os.ReadFile(real); string(data) != "*.log\n.docs/\n" {
		t.Fatalf("%s = %q", real, data)
	}
}
