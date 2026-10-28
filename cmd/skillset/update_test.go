package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zenkoding/skillset/internal/source"
)

// gitRun runs git in dir.
func gitRun(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

// commitCatalog commits a catalog whose one set's instructions say text, and tags
// it when tag is set.
func commitCatalog(t *testing.T, repo, text, tag string) {
	t.Helper()
	writeFile(t, filepath.Join(repo, "config/sets/a/set.toml"), "name = \"a\"\nhosts = [\"claude\"]")
	writeFile(t, filepath.Join(repo, "config/sets/a/shared.md"), text)
	gitRun(t, repo, "add", "-A")
	gitRun(t, repo, "commit", "-qm", text)
	if tag != "" {
		gitRun(t, repo, "tag", tag)
	}
}

func TestRemoteCatalog(t *testing.T) {
	repo, home := t.TempDir(), t.TempDir()
	gitRun(t, repo, "init", "-q", "-b", "main")
	commitCatalog(t, repo, "one\n", "v0.1.0")
	commitCatalog(t, repo, "two\n", "v0.2.0")
	url := "file://" + repo
	claudeMD := filepath.Join(home, ".claude/CLAUDE.md")
	lockPath := filepath.Join(home, LockPath)

	steps := []struct {
		name  string
		setup func()
		args  []string
		ref   string // pinned ref after the step
		text  string // CLAUDE.md after the step
	}{
		{"first apply pins the newest tag", nil, []string{"apply"}, "v0.2.0", "two\n"},
		{"apply keeps the pin", func() { commitCatalog(t, repo, "three\n", "v0.3.0") }, []string{"apply"}, "v0.2.0", "two\n"},
		{"update moves to the newest tag", nil, []string{"update"}, "v0.3.0", "three\n"},
		{"update --ref switches back", nil, []string{"update", "--ref", "v0.1.0"}, "v0.1.0", "one\n"},
		{"update --ref main", nil, []string{"update", "--ref", "main"}, "main", "three\n"},
		{"main is fetched fresh", func() { commitCatalog(t, repo, "four\n", "") }, []string{"update", "--ref", "main"}, "main", "four\n"},
	}
	for _, st := range steps {
		t.Run(st.name, func(t *testing.T) {
			if st.setup != nil {
				st.setup()
			}
			out, err := run(t, append(st.args, "--url", url, "--home", home)...)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(out, "catalog: "+url+" at "+st.ref+" (") || !strings.Contains(out, "installed "+st.ref) {
				t.Fatalf("output = %s", out)
			}
			lock, err := source.ReadLock(lockPath)
			if err != nil || lock.Ref != st.ref || lock.Commit == "" {
				t.Fatalf("lock = %+v, %v; want ref %s", lock, err, st.ref)
			}
			data, err := os.ReadFile(claudeMD)
			if err != nil || string(data) != st.text {
				t.Fatalf("CLAUDE.md = %q, %v; want %q", data, err, st.text)
			}
			if !strings.HasPrefix(mustReadlink(t, claudeMD), filepath.Join(home, CacheDir, st.ref)+"/") {
				t.Fatalf("CLAUDE.md links to %s, want inside the %s cache", mustReadlink(t, claudeMD), st.ref)
			}
			p, err := run(t, "plan", "--url", url, "--home", home)
			if err != nil || !strings.Contains(p, "0 to render, 0 to link") {
				t.Fatalf("plan after %s:\n%s\n%v", st.name, p, err)
			}
		})
	}
	for _, ref := range []string{"v0.1.0", "v0.2.0", "v0.3.0"} {
		if _, err := os.Stat(filepath.Join(home, CacheDir, ref)); err != nil {
			t.Errorf("cache for %s gone: %v", ref, err)
		}
	}
	// Refs that are not a single cache directory fail before the cache is touched.
	for _, ref := range []string{"..", "feature/x", ".hidden"} {
		if _, err := run(t, "update", "--ref", ref, "--url", url, "--home", home); err == nil || !strings.Contains(err.Error(), "invalid ref") {
			t.Fatalf("update --ref %s = %v, want invalid ref", ref, err)
		}
	}
	if _, err := os.Stat(filepath.Join(home, CacheDir, "main")); err != nil {
		t.Fatalf("cache touched by an invalid ref: %v", err)
	}

	// An unknown tag fails and leaves the pin alone.
	if _, err := run(t, "update", "--ref", "v9.9.9", "--url", url, "--home", home); err == nil {
		t.Fatal("update to a missing tag succeeded")
	}
	if lock, _ := source.ReadLock(lockPath); lock.Ref != "main" {
		t.Fatalf("lock after a failed update = %+v", lock)
	}
	// A failed fresh fetch keeps the cached branch, so links still work.
	if _, err := run(t, "update", "--ref", "main", "--url", "file://"+filepath.Join(repo, "gone"), "--home", home); err == nil {
		t.Fatal("update from a missing repo succeeded")
	}
	if data, err := os.ReadFile(claudeMD); err != nil || string(data) != "four\n" {
		t.Fatalf("CLAUDE.md after a failed fetch = %q, %v", data, err)
	}
	if stale, _ := filepath.Glob(filepath.Join(home, CacheDir, ".stale-*")); len(stale) > 0 {
		t.Fatalf("leftover %v", stale)
	}
	if _, err := run(t, "update", "--catalog", repo, "--home", home); err == nil {
		t.Error("update accepted --catalog")
	}
}

func mustReadlink(t *testing.T, path string) string {
	t.Helper()
	target, err := os.Readlink(path)
	if err != nil {
		t.Fatal(err)
	}
	return target
}

// TestSettingsFragments installs two sets with settings and checks they merge.
func TestSettingsFragments(t *testing.T) {
	root, home := t.TempDir(), t.TempDir()
	for p, body := range map[string]string{
		"sets/a/set.toml":      "name = \"a\"\nhosts = [\"claude\"]",
		"sets/a/settings.json": `{"permissions": {"allow": ["Read"]}, "model": "x"}`,
		"sets/b/set.toml":      "name = \"b\"\nhosts = [\"claude\"]",
		"sets/b/settings.json": `{"permissions": {"allow": ["Grep", "Read"]}}`,
	} {
		writeFile(t, filepath.Join(root, p), body)
	}
	if _, err := run(t, "apply", "--catalog", root, "--home", home); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(home, ".claude/settings.json"))
	if err != nil {
		t.Fatal(err)
	}
	want := "{\n  \"model\": \"x\",\n  \"permissions\": {\n    \"allow\": [\n      \"Read\",\n      \"Grep\"\n    ]\n  }\n}\n"
	if string(data) != want {
		t.Fatalf("settings.json = %s, want %s", data, want)
	}
}
