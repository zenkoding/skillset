package source

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// remote is a local bare repo with a work tree used to add commits and tags.
type remote struct {
	t    *testing.T
	url  string
	work string
}

func newRemote(t *testing.T) *remote {
	t.Helper()
	root := t.TempDir()
	r := &remote{t: t, url: filepath.Join(root, "catalog.git"), work: filepath.Join(root, "work")}
	r.run(root, "init", "--quiet", "--bare", r.url)
	r.run(root, "init", "--quiet", "--initial-branch", "main", r.work)
	r.run(r.work, "remote", "add", "origin", r.url)
	return r
}

func (r *remote) run(dir string, args ...string) string {
	r.t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=Test", "GIT_AUTHOR_EMAIL=test@example.com",
		"GIT_COMMITTER_NAME=Test", "GIT_COMMITTER_EMAIL=test@example.com",
		"GIT_CONFIG_NOSYSTEM=1", "HOME="+r.t.TempDir(),
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		r.t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

// commit writes content to set.txt, commits, pushes main and returns the commit.
func (r *remote) commit(content string) string {
	r.t.Helper()
	if err := os.WriteFile(filepath.Join(r.work, "set.txt"), []byte(content), 0o644); err != nil {
		r.t.Fatal(err)
	}
	r.run(r.work, "add", "set.txt")
	r.run(r.work, "commit", "--quiet", "-m", content)
	r.run(r.work, "push", "--quiet", "origin", "main")
	return r.run(r.work, "rev-parse", "HEAD")
}

// tag tags HEAD (annotated when annotated is true) and pushes the tag.
func (r *remote) tag(name string, annotated bool) {
	r.t.Helper()
	if annotated {
		r.run(r.work, "tag", "-a", "-m", name, name)
	} else {
		r.run(r.work, "tag", name)
	}
	r.run(r.work, "push", "--quiet", "origin", name)
}

func TestLatestTag(t *testing.T) {
	tests := []struct {
		name    string
		tags    []string
		want    string
		wantErr bool
	}{
		{name: "single", tags: []string{"v0.1.0"}, want: "v0.1.0"},
		{name: "numeric not lexical", tags: []string{"v1.9.0", "v1.10.0", "v1.2.3"}, want: "v1.10.0"},
		{name: "major wins", tags: []string{"v1.99.99", "v2.0.0"}, want: "v2.0.0"},
		{name: "skips pre-release and build", tags: []string{"v1.0.0", "v2.0.0-rc.1", "v1.1.0+build"}, want: "v1.0.0"},
		{name: "skips non-semver", tags: []string{"v1.0.0", "latest", "1.5.0", "v1.5", "release-2"}, want: "v1.0.0"},
		{name: "no tags", tags: nil, wantErr: true},
		{name: "only non-release tags", tags: []string{"v3.0.0-beta", "nightly"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newRemote(t)
			r.commit("one")
			for i, tag := range tt.tags {
				r.tag(tag, i%2 == 0)
			}

			got, err := LatestTag(context.Background(), r.url)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("LatestTag() = %q, want error", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("LatestTag() error: %v", err)
			}
			if got != tt.want {
				t.Errorf("LatestTag() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLatestTagUnreachable(t *testing.T) {
	if _, err := LatestTag(context.Background(), filepath.Join(t.TempDir(), "missing.git")); err == nil {
		t.Fatal("LatestTag() on missing repo: want error")
	}
}

func readSet(t *testing.T, dir string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(dir, "set.txt"))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestFetch(t *testing.T) {
	ctx := context.Background()
	r := newRemote(t)
	c1 := r.commit("one")
	r.tag("v1.0.0", true)
	c2 := r.commit("two")
	r.tag("v1.1.0", false)
	cache := filepath.Join(t.TempDir(), "cache")

	dir1, got1, err := Fetch(ctx, r.url, "v1.0.0", cache)
	if err != nil {
		t.Fatalf("Fetch v1.0.0: %v", err)
	}
	if dir1 != filepath.Join(cache, "v1.0.0") || got1 != c1 || readSet(t, dir1) != "one" {
		t.Fatalf("Fetch v1.0.0 = (%q, %q), content %q; want commit %q, content one", dir1, got1, readSet(t, dir1), c1)
	}

	dir2, got2, err := Fetch(ctx, r.url, "v1.1.0", cache)
	if err != nil {
		t.Fatalf("Fetch v1.1.0: %v", err)
	}
	if got2 != c2 || readSet(t, dir2) != "two" {
		t.Fatalf("Fetch v1.1.0 commit %q content %q; want %q two", got2, readSet(t, dir2), c2)
	}
	if readSet(t, dir1) != "one" {
		t.Error("second fetch changed the first ref's checkout")
	}

	t.Run("branch ref", func(t *testing.T) {
		dir, commit, err := Fetch(ctx, r.url, "main", cache)
		if err != nil {
			t.Fatalf("Fetch main: %v", err)
		}
		if commit != c2 || readSet(t, dir) != "two" {
			t.Errorf("Fetch main commit %q content %q; want %q two", commit, readSet(t, dir), c2)
		}
	})

	t.Run("reuses existing checkout", func(t *testing.T) {
		// Prove no network: a dead URL still resolves from the cache.
		dir, commit, err := Fetch(ctx, filepath.Join(t.TempDir(), "gone.git"), "v1.0.0", cache)
		if err != nil {
			t.Fatalf("Fetch cached: %v", err)
		}
		if dir != dir1 || commit != c1 {
			t.Errorf("Fetch cached = (%q, %q), want (%q, %q)", dir, commit, dir1, c1)
		}
	})

	t.Run("failed fetch leaves nothing", func(t *testing.T) {
		if _, _, err := Fetch(ctx, r.url, "v9.9.9", cache); err == nil {
			t.Fatal("Fetch of unknown ref: want error")
		}
		entries, err := os.ReadDir(cache)
		if err != nil {
			t.Fatal(err)
		}
		var names []string
		for _, e := range entries {
			names = append(names, e.Name())
		}
		if got, want := strings.Join(names, ","), "main,v1.0.0,v1.1.0"; got != want {
			t.Errorf("cache entries = %s, want %s", got, want)
		}
	})
}

func TestFetchDashURL(t *testing.T) {
	cache := t.TempDir()
	marker := filepath.Join(t.TempDir(), "ran")
	if _, _, err := Fetch(context.Background(), "--upload-pack=touch "+marker, "v1.0.0", cache); err == nil {
		t.Fatal("Fetch() with dash URL: want error")
	}
	if _, err := os.Stat(marker); err == nil {
		t.Error("dash URL was parsed as --upload-pack")
	}
}

func TestFetchRejectsBadRef(t *testing.T) {
	for _, ref := range []string{"", ".", "..", "-c", ".hidden", "feat/x", `a\b`} {
		t.Run(ref, func(t *testing.T) {
			cache := t.TempDir()
			if _, _, err := Fetch(context.Background(), "unused", ref, cache); err == nil {
				t.Errorf("Fetch(%q): want error", ref)
			}
		})
	}
}
