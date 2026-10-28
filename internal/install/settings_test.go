package install

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// withSettings adds the settings files to a fixture repo.
func withSettings(t *testing.T, repo string) {
	t.Helper()
	write(t, filepath.Join(repo, "claude/settings.json"), `{"model": "opus"}`+"\n")
	write(t, filepath.Join(repo, "claude/statusline.sh"), "#!/bin/sh\n")
	write(t, filepath.Join(repo, "claude/output-styles/dry.md"), "dry\n")
	write(t, filepath.Join(repo, "dcg/config.toml"), "x = 1\n")
	write(t, filepath.Join(repo, "dcg/packs/p.yaml"), "p\n")
}

func TestConfigLinks(t *testing.T) {
	repo, home := fixture(t)
	withSettings(t, repo)
	install(t, repo, home)
	tests := []struct{ path, target string }{
		{".claude/statusline.sh", "claude/statusline.sh"},
		{".claude/output-styles", "claude/output-styles"},
		{".config/dcg/config.toml", "dcg/config.toml"},
		{".config/dcg/packs", "dcg/packs"},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if !linksTo(filepath.Join(home, tt.path), filepath.Join(repo, tt.target)) {
				t.Fatalf("%s does not link to %s", tt.path, tt.target)
			}
		})
	}
	fi, err := os.Lstat(filepath.Join(home, ".claude/settings.json"))
	if err != nil || !fi.Mode().IsRegular() {
		t.Fatalf("settings.json is not a copied file: %v", err)
	}
}

func TestSettingsMerge(t *testing.T) {
	tests := []struct {
		name      string
		repo1     string // repo at first install
		local     string // local file after first install
		repo2     string // repo at second install
		want      string // settings.json after second install
		conflicts []string
	}{
		{
			name:  "untouched copy follows the repo",
			repo1: `{"a": 1}`, local: `{"a": 1}`, repo2: `{"a": 2}`,
			want: `{"a": 2}`,
		},
		{
			name:  "local edit kept when the repo is unchanged",
			repo1: `{"a": 1}`, local: `{"a": 1, "b": true}`, repo2: `{"a": 1}`,
			want: `{"a": 1, "b": true}`,
		},
		{
			name:  "edits to different keys merge",
			repo1: `{"a": 1, "m": "x"}`, local: `{"a": 1, "m": "y"}`, repo2: `{"a": 2, "m": "x"}`,
			want: `{"a": 2, "m": "y"}`,
		},
		{
			name:  "repo removal applies to an untouched key",
			repo1: `{"a": 1, "b": 2}`, local: `{"a": 1, "b": 2, "c": 3}`, repo2: `{"a": 1}`,
			want: `{"a": 1, "c": 3}`,
		},
		{
			name:  "string lists merge as sets",
			repo1: `{"allow": ["x", "y"]}`, local: `{"allow": ["x", "y", "local"]}`, repo2: `{"allow": ["x", "repo"]}`,
			want: `{"allow": ["x", "local", "repo"]}`,
		},
		{
			name:  "same key changed both ways conflicts",
			repo1: `{"p": {"mode": "a"}}`, local: `{"p": {"mode": "b"}}`, repo2: `{"p": {"mode": "c"}}`,
			want:      `{"p": {"mode": "c"}}`,
			conflicts: []string{"/p/mode"},
		},
		{
			name:  "same change on both sides is no conflict",
			repo1: `{"a": 1}`, local: `{"a": 5}`, repo2: `{"a": 5}`,
			want: `{"a": 5}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, home := fixture(t)
			src, dst := filepath.Join(repo, "claude/settings.json"), filepath.Join(home, ".claude/settings.json")
			write(t, src, tt.repo1)
			install(t, repo, home)
			write(t, dst, tt.local)
			write(t, src, tt.repo2)
			p, err := makePlan(repo, home)
			if err != nil {
				t.Fatal(err)
			}
			var keys []string
			for _, c := range p.Copies[0].Conflicts {
				keys = append(keys, c.Key)
			}
			if len(keys) != len(tt.conflicts) || (len(keys) > 0 && keys[0] != tt.conflicts[0]) {
				t.Fatalf("conflicts %v, want %v", keys, tt.conflicts)
			}
			if len(keys) > 0 {
				var ce *ConflictError
				if err := p.Apply(); !errors.As(err, &ce) {
					t.Fatalf("Apply error = %v, want ConflictError", err)
				}
				if err := p.Resolve(dst, keys[0], false); err != nil {
					t.Fatal(err)
				}
			}
			if err := p.Apply(); err != nil {
				t.Fatal(err)
			}
			assertJSON(t, dst, tt.want)
		})
	}
}

func TestSettingsConflictKeepLocal(t *testing.T) {
	repo, home := fixture(t)
	src, dst := filepath.Join(repo, "claude/settings.json"), filepath.Join(home, ".claude/settings.json")
	write(t, src, `{"a": 1, "b": 1}`)
	install(t, repo, home)
	write(t, dst, `{"a": 2, "b": 1}`)
	write(t, src, `{"a": 3, "b": 9}`)
	p, err := makePlan(repo, home)
	if err != nil {
		t.Fatal(err)
	}
	if err := p.Resolve(dst, "/a", true); err != nil {
		t.Fatal(err)
	}
	if err := p.Apply(); err != nil {
		t.Fatal(err)
	}
	assertJSON(t, dst, `{"a": 2, "b": 9}`)
	// The choice sticks: the next run sees nothing new on either side.
	install(t, repo, home)
	assertJSON(t, dst, `{"a": 2, "b": 9}`)
}

func TestSettingsClash(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T, dst string)
	}{
		{"foreign link", func(t *testing.T, dst string) {
			os.MkdirAll(filepath.Dir(dst), 0o755)
			os.Symlink(filepath.Join(t.TempDir(), "s.json"), dst)
		}},
		{"not JSON", func(t *testing.T, dst string) { write(t, dst, "not json") }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, home := fixture(t)
			withSettings(t, repo)
			tt.setup(t, filepath.Join(home, ".claude/settings.json"))
			p, err := makePlan(repo, home)
			if err != nil {
				t.Fatal(err)
			}
			var ce *ClashError
			if err := p.Apply(); !errors.As(err, &ce) {
				t.Fatalf("Apply error = %v, want ClashError", err)
			}
		})
	}
}

func TestUninstallSettings(t *testing.T) {
	tests := []struct {
		name    string
		edit    bool
		removed bool
	}{
		{"unchanged copy is removed", false, true},
		{"edited copy is kept", true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, home := fixture(t)
			withSettings(t, repo)
			install(t, repo, home)
			dst := filepath.Join(home, ".claude/settings.json")
			if tt.edit {
				write(t, dst, `{"model": "mine"}`)
			}
			if _, _, err := Uninstall(home, true, repo); err != nil {
				t.Fatal(err)
			}
			_, err := os.Lstat(dst)
			if gone := err != nil; gone != tt.removed {
				t.Fatalf("removed=%v, want %v", gone, tt.removed)
			}
			if _, err := os.Stat(filepath.Join(home, CopyStateDir)); err == nil {
				t.Error("copy state left behind")
			}
		})
	}
}

func assertJSON(t *testing.T, path, want string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	got, err := decodeJSON(data)
	if err != nil {
		t.Fatal(err)
	}
	w, _ := decodeJSON([]byte(want))
	if show(got) != show(w) {
		t.Fatalf("%s = %s, want %s", path, show(got), show(w))
	}
}
