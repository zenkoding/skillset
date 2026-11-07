package catalog

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestProjectKnowledgeHook runs the session-start hook in a scratch repo and
// checks it prints the private store and catalog, and nothing from the repo.
func TestProjectKnowledgeHook(t *testing.T) {
	if _, err := exec.LookPath("bash"); err != nil {
		t.Skip("bash not installed")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	hook, err := filepath.Abs("../../config/claude/hooks/project-knowledge.sh")
	if err != nil {
		t.Fatal(err)
	}
	store := filepath.Join(t.TempDir(), "projects", "github-com-foo-bar")

	tests := []struct {
		name    string
		repo    bool
		catalog string // private README.md; empty means none
		want    []string
		notWant []string
	}{
		{name: "outside a repo", notWant: []string{"Project knowledge"}},
		{
			name:    "repo without a catalog",
			repo:    true,
			want:    []string{"Private store: " + store},
			notWant: []string{"PUBLIC-CATALOG", "Private catalog"},
		},
		{
			name:    "repo with a catalog",
			repo:    true,
			catalog: "- [a.md](a.md): PRIVATE-LINE\n",
			want:    []string{"Private store: " + store, "PRIVATE-LINE"},
			notWant: []string{"PUBLIC-CATALOG"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			if tt.repo {
				for _, args := range [][]string{
					{"init", "-q"},
					{"remote", "add", "origin", "git@github.com:Foo/Bar.git"},
				} {
					if out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).CombinedOutput(); err != nil {
						t.Fatalf("git %v: %v\n%s", args, err, out)
					}
				}
				// A public catalog from the old layout must not be printed.
				if err := os.MkdirAll(filepath.Join(dir, "docs", "agents"), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(dir, "docs", "agents", "README.md"), []byte("PUBLIC-CATALOG\n"), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			os.RemoveAll(store)
			if tt.catalog != "" {
				if err := os.MkdirAll(store, 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(store, "README.md"), []byte(tt.catalog), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			cmd := exec.Command("bash", hook)
			cmd.Dir = dir
			cmd.Env = append(os.Environ(),
				"CLAUDE_PROJECT_DIR="+dir,
				"AGENTS_DATA_HOME="+filepath.Dir(filepath.Dir(store)),
				"GIT_CEILING_DIRECTORIES="+filepath.Dir(dir),
			)
			out, err := cmd.Output()
			if err != nil {
				t.Fatalf("hook failed: %v", err)
			}
			for _, w := range tt.want {
				if !strings.Contains(string(out), w) {
					t.Errorf("output lacks %q:\n%s", w, out)
				}
			}
			for _, w := range tt.notWant {
				if strings.Contains(string(out), w) {
					t.Errorf("output has %q:\n%s", w, out)
				}
			}
		})
	}
}
