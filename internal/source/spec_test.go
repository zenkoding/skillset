package source

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		arg   string
		url   string
		ref   string // "" means none given
		name  string
		error string
	}{
		{arg: "zenkoding/skillset", url: "https://github.com/zenkoding/skillset", name: "skillset"},
		{arg: "zenkoding/skillset@v1.2.0", url: "https://github.com/zenkoding/skillset", ref: "v1.2.0", name: "skillset"},
		{arg: "zenkoding/skillset@feature/x", url: "https://github.com/zenkoding/skillset", ref: "feature/x", name: "skillset"},
		{arg: "https://github.com/o/r", url: "https://github.com/o/r", name: "r"},
		{arg: "https://github.com/o/r.git@main", url: "https://github.com/o/r.git", ref: "main", name: "r"},
		{arg: "ssh://user@host.example/o/r@feature/a/b", url: "ssh://user@host.example/o/r", ref: "feature/a/b", name: "r"},
		{arg: "ssh://git@host.example:22/o/r.git@v1", url: "ssh://git@host.example:22/o/r.git", ref: "v1", name: "r"},
		{arg: "git@github.com:o/r.git", url: "git@github.com:o/r.git", name: "r"},
		{arg: "git@github.com:o/r.git@feature/x", url: "git@github.com:o/r.git", ref: "feature/x", name: "r"},
		{arg: "file:///srv/git/sets.git@0123456789abcdef0123456789abcdef01234567", url: "file:///srv/git/sets.git", ref: "0123456789abcdef0123456789abcdef01234567", name: "sets"},
		{arg: "", error: "empty source"},
		{arg: "skillset", error: "not a source"},
		{arg: "o/r@", error: "empty ref"},
		{arg: "o/r@-c", error: "invalid ref"},
		{arg: "o/r@a..b", error: "invalid ref"},
		{arg: "-o/r", error: "not a source"},
		{arg: "https://user:s3cret@github.com/o/r@main", error: "credential"},
		{arg: "https://s3cret@github.com/o/r", error: "credential"},
		{arg: "ssh://git:s3cret@host.example/o/r", error: "credential"},
	}
	for _, tt := range tests {
		t.Run(tt.arg, func(t *testing.T) {
			url, ref, err := Parse(tt.arg)
			if tt.error != "" {
				if err == nil || !strings.Contains(err.Error(), tt.error) {
					t.Fatalf("Parse = %q, %q, %v; want error %q", url, ref, err, tt.error)
				}
				if strings.Contains(err.Error(), "s3cret") {
					t.Fatalf("error shows the secret: %v", err)
				}
				return
			}
			if err != nil || url != tt.url || ref != tt.ref {
				t.Fatalf("Parse = %q, %q, %v; want %q, %q", url, ref, err, tt.url, tt.ref)
			}
			if got := Name(url); got != tt.name {
				t.Fatalf("Name(%q) = %q, want %q", url, got, tt.name)
			}
		})
	}
}

func TestResolveAndFetchCommit(t *testing.T) {
	ctx := context.Background()
	r := newRemote(t)
	c1 := r.commit("one")
	r.tag("v1.0.0", true)
	r.tag("plain", false)
	r.run(r.work, "push", "--quiet", "origin", "main:feature/x")
	c2 := r.commit("two")

	tests := []struct {
		ref   string
		want  string
		error string
	}{
		{ref: "main", want: c2},
		{ref: "feature/x", want: c1},
		{ref: "v1.0.0", want: c1}, // annotated: the commit, not the tag object
		{ref: "plain", want: c1},
		{ref: c1, want: c1},
		{ref: "x", error: `ref "x" not found`},
		{ref: c1[:7], error: "full commit hash"},
		{ref: "-c", error: "invalid ref"},
	}
	for _, tt := range tests {
		t.Run(tt.ref, func(t *testing.T) {
			got, err := Resolve(ctx, r.url, tt.ref)
			if tt.error != "" {
				if err == nil || !strings.Contains(err.Error(), tt.error) {
					t.Fatalf("Resolve = %q, %v; want error %q", got, err, tt.error)
				}
				return
			}
			if err != nil || got != tt.want {
				t.Fatalf("Resolve = %q, %v; want %q", got, err, tt.want)
			}
		})
	}

	t.Run("fetch", func(t *testing.T) {
		cache := filepath.Join(t.TempDir(), "cache")
		dir, err := FetchCommit(ctx, r.url, c1, cache)
		if err != nil {
			t.Fatal(err)
		}
		if dir != filepath.Join(cache, c1) || readSet(t, dir) != "one" {
			t.Fatalf("FetchCommit = %q, content %q", dir, readSet(t, dir))
		}
		// A cached commit needs no network.
		again, err := FetchCommit(ctx, filepath.Join(t.TempDir(), "gone.git"), c1, cache)
		if err != nil || again != dir {
			t.Fatalf("FetchCommit cached = %q, %v", again, err)
		}
		if _, err := FetchCommit(ctx, r.url, strings.Repeat("0", 40), cache); err == nil {
			t.Fatal("FetchCommit of a missing commit: want error")
		}
		entries, _ := os.ReadDir(cache)
		if len(entries) != 1 {
			t.Fatalf("cache holds %d entries after a failed fetch, want 1", len(entries))
		}
		if _, err := FetchCommit(ctx, r.url, "main", cache); err == nil {
			t.Fatal("FetchCommit of a branch name: want error")
		}
	})
}
