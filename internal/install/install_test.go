package install

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
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

func skill(name string) string { return "---\nname: " + name + "\n---\nbody\n" }

// fixture builds a repo with two skills, a draft, an empty skill, a name
// mismatch and both instruction files, plus an empty temp HOME.
func fixture(t *testing.T) (repo, home string) {
	t.Helper()
	repo, home = t.TempDir(), t.TempDir()
	write(t, filepath.Join(repo, "skills/eng/alpha/SKILL.md"), skill("alpha"))
	write(t, filepath.Join(repo, "skills/sdlc/beta/SKILL.md"), skill("beta"))
	write(t, filepath.Join(repo, "skills/in-progress/draft/SKILL.md"), skill("draft"))
	write(t, filepath.Join(repo, "skills/sdlc/empty/SKILL.md"), "")
	write(t, filepath.Join(repo, "skills/eng/wrong/SKILL.md"), skill("other"))
	write(t, filepath.Join(repo, "instructions/shared.md"), "shared\n")
	write(t, filepath.Join(repo, "instructions/claude.md"), "claude\n")
	write(t, filepath.Join(repo, "instructions/codex.md"), "codex\n")
	return repo, home
}

func install(t *testing.T, repo, home string) *Plan {
	t.Helper()
	p, err := makePlan(repo, home)
	if err != nil {
		t.Fatal(err)
	}
	if err := p.Apply(); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestInstall(t *testing.T) {
	repo, home := fixture(t)
	install(t, repo, home)
	tests := []struct {
		path   string
		target string // "" means must not exist
	}{
		{".claude/skills/alpha", "skills/eng/alpha"},
		{".codex/skills/alpha", "skills/eng/alpha"},
		{".agents/skills/beta", "skills/sdlc/beta"},
		{".claude/CLAUDE.md", "build/claude/CLAUDE.md"},
		{".codex/AGENTS.md", "build/codex/AGENTS.md"},
		{".claude/skills/draft", ""},
		{".agents/skills/draft", ""},
		{".claude/skills/empty", ""},
		{".claude/skills/wrong", ""},
		{".claude/skills/other", ""},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			path := filepath.Join(home, tt.path)
			if tt.target == "" {
				if _, err := os.Lstat(path); err == nil {
					t.Fatalf("%s should not be installed", tt.path)
				}
				return
			}
			if !linksTo(path, filepath.Join(repo, tt.target)) {
				t.Fatalf("%s does not link to %s", tt.path, tt.target)
			}
		})
	}
}

func TestNoDraftInstalled(t *testing.T) {
	repo, home := fixture(t)
	install(t, repo, home)
	for _, dir := range SkillDirs {
		entries, _ := os.ReadDir(filepath.Join(home, dir))
		for _, e := range entries {
			if t2, _ := linkTarget(filepath.Join(home, dir, e.Name())); filepath.Base(filepath.Dir(t2)) == DraftDir {
				t.Errorf("draft installed: %s/%s", dir, e.Name())
			}
		}
	}
}

func TestReinstallIsIdempotent(t *testing.T) {
	repo, home := fixture(t)
	install(t, repo, home)
	p := install(t, repo, home)
	for _, r := range p.Renders {
		if !r.Exists {
			t.Errorf("second run would re-render %s", r.Path)
		}
	}
	if len(p.Remove) != 0 || len(p.Clashes) != 0 {
		t.Fatalf("second run removes %v, clashes %v", p.Remove, p.Clashes)
	}
	for _, l := range p.Links {
		if !l.Exists {
			t.Errorf("second run would recreate %s", l.Path)
		}
	}
}

func TestStaleCleanup(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T, repo, home string) string // returns path to check
		removed bool
	}{
		{
			name: "manifest link to a deleted skill",
			setup: func(t *testing.T, repo, home string) string {
				write(t, filepath.Join(repo, "skills/eng/gone/SKILL.md"), skill("gone"))
				install(t, repo, home)
				os.Remove(filepath.Join(repo, "skills/eng/gone/SKILL.md"))
				return filepath.Join(home, ".claude/skills/gone")
			},
			removed: true,
		},
		{
			name: "real directory is left alone",
			setup: func(t *testing.T, repo, home string) string {
				path := filepath.Join(home, ".claude/skills/tour")
				write(t, filepath.Join(path, "notes.md"), "mine")
				return path
			},
		},
		{
			name: "unlisted link into repo is left alone",
			setup: func(t *testing.T, repo, home string) string {
				path := filepath.Join(home, ".claude/skills/someone-else")
				os.MkdirAll(filepath.Dir(path), 0o755)
				os.Symlink(filepath.Join(repo, "skills/eng/alpha"), path)
				return path
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, home := fixture(t)
			path := tt.setup(t, repo, home)
			install(t, repo, home)
			_, err := os.Lstat(path)
			if gone := err != nil; gone != tt.removed {
				t.Fatalf("removed=%v, want %v", gone, tt.removed)
			}
		})
	}
}

func TestCollisions(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T, repo, home string)
		want  []string
	}{
		{
			name: "matt skill dir and foreign link",
			setup: func(t *testing.T, repo, home string) {
				write(t, filepath.Join(home, ".claude/skills/alpha/SKILL.md"), skill("alpha"))
				os.MkdirAll(filepath.Join(home, ".agents/skills"), 0o755)
				os.Symlink(t.TempDir(), filepath.Join(home, ".agents/skills/beta"))
			},
			want: []string{".agents/skills/beta", ".claude/skills/alpha"},
		},
		{
			name: "real instruction file",
			setup: func(t *testing.T, repo, home string) {
				write(t, filepath.Join(home, ".claude/CLAUDE.md"), "mine")
			},
			want: []string{".claude/CLAUDE.md"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, home := fixture(t)
			tt.setup(t, repo, home)
			p, err := makePlan(repo, home)
			if err != nil {
				t.Fatal(err)
			}
			var ce *ClashError
			if err := p.Apply(); !errors.As(err, &ce) {
				t.Fatalf("Apply error = %v, want ClashError", err)
			}
			if len(ce.Paths) != len(tt.want) {
				t.Fatalf("clashes %v, want %v", ce.Paths, tt.want)
			}
			for i, w := range tt.want {
				if ce.Paths[i] != filepath.Join(home, w) {
					t.Errorf("clash %d = %s, want %s", i, ce.Paths[i], w)
				}
			}
			// Nothing changed: no other link was created.
			if _, err := os.Lstat(filepath.Join(home, ".codex/skills/alpha")); err == nil {
				t.Error("install created links despite clashes")
			}
		})
	}
}

func TestUninstall(t *testing.T) {
	repo, home := fixture(t)
	install(t, repo, home)
	foreign := filepath.Join(home, ".claude/skills/beta")
	os.Remove(foreign)
	os.Symlink(t.TempDir(), foreign)

	dry, _, err := Uninstall(home, false, repo)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(filepath.Join(home, ".claude/skills/alpha")); err != nil {
		t.Fatal("dry run removed a link")
	}
	removed, kept, err := Uninstall(home, true, repo)
	if err != nil {
		t.Fatal(err)
	}
	if len(dry) != len(removed) || len(kept) != 1 || kept[0] != foreign {
		t.Fatalf("dry %d removed %d kept %v", len(dry), len(removed), kept)
	}
	tests := []struct {
		path string
		gone bool
	}{
		{".claude/skills/alpha", true},
		{".agents/skills/beta", true},
		{".codex/AGENTS.md", true},
		{".claude/skills/beta", false},
		{ManifestPath, true},
	}
	for _, tt := range tests {
		_, err := os.Lstat(filepath.Join(home, tt.path))
		if (err != nil) != tt.gone {
			t.Errorf("%s gone=%v, want %v", tt.path, err != nil, tt.gone)
		}
	}
}

func TestRenderInstructions(t *testing.T) {
	repo, home := fixture(t)
	plan, err := makePlan(repo, home)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(repo, BuildDir)); err == nil {
		t.Fatal("MakePlan wrote build/")
	}
	if err := plan.Apply(); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		link string
		want string
	}{
		{".claude/CLAUDE.md", "shared\nclaude\n"},
		{".codex/AGENTS.md", "shared\ncodex\n"},
	}
	for _, tt := range tests {
		t.Run(tt.link, func(t *testing.T) {
			got, err := os.ReadFile(filepath.Join(home, tt.link))
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
	// An edited source is re-rendered on the next run.
	write(t, filepath.Join(repo, "instructions/codex.md"), "codex v2\n")
	install(t, repo, home)
	got, _ := os.ReadFile(filepath.Join(home, ".codex/AGENTS.md"))
	if string(got) != "shared\ncodex v2\n" {
		t.Fatalf("not re-rendered: %q", got)
	}
}

func TestDuplicateSkillNames(t *testing.T) {
	repo, _ := fixture(t)
	write(t, filepath.Join(repo, "skills/sdlc/alpha/SKILL.md"), skill("alpha"))
	if _, _, err := repoSource(repo).discover(); err == nil {
		t.Fatal("want duplicate-name error")
	}
}
