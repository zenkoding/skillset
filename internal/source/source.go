// Package source fetches release snapshots of the catalog repo into a local
// cache and records the pinned ref in a lockfile.
package source

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// semverTag matches release tags only: vX.Y.Z with no pre-release or build suffix.
var semverTag = regexp.MustCompile(`^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$`)

// LatestTag returns the newest semver tag (vX.Y.Z) of url, via git ls-remote.
// Non-semver and pre-release tags are ignored.
func LatestTag(ctx context.Context, url string) (string, error) {
	out, err := git(ctx, "", "ls-remote", "--tags", "--refs", "--", url)
	if err != nil {
		return "", fmt.Errorf("list tags of %s: %w", url, err)
	}

	var best string
	var bestVer [3]int
	for _, line := range strings.Split(out, "\n") {
		_, ref, ok := strings.Cut(strings.TrimSpace(line), "\t")
		if !ok {
			continue
		}
		tag := strings.TrimPrefix(ref, "refs/tags/")
		ver, ok := parseSemver(tag)
		if !ok {
			continue
		}
		if best == "" || less(bestVer, ver) {
			best, bestVer = tag, ver
		}
	}
	if best == "" {
		return "", fmt.Errorf("no release tag (vX.Y.Z) in %s", url)
	}
	return best, nil
}

// Fetch makes cacheDir/<ref> hold url at ref and returns that dir and the commit.
// An existing complete checkout is reused; old refs are never deleted.
//
// The checkout is built in a temp dir inside cacheDir and renamed into place,
// so a failed fetch leaves no half checkout behind.
func Fetch(ctx context.Context, url, ref, cacheDir string) (dir, commit string, err error) {
	if err := ValidRef(ref); err != nil {
		return "", "", err
	}
	dir = filepath.Join(cacheDir, ref)

	if info, statErr := os.Stat(dir); statErr == nil && info.IsDir() {
		commit, err = headCommit(ctx, dir)
		return dir, commit, err
	}

	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return "", "", fmt.Errorf("create cache dir: %w", err)
	}
	tmp, err := os.MkdirTemp(cacheDir, ".fetch-*")
	if err != nil {
		return "", "", fmt.Errorf("create temp dir: %w", err)
	}
	defer func() {
		// Only left over when something failed; after the rename it is gone.
		_ = os.RemoveAll(tmp)
	}()

	steps := [][]string{
		{"init", "--quiet"},
		{"fetch", "--quiet", "--depth", "1", "--", url, ref},
		{"-c", "advice.detachedHead=false", "checkout", "--quiet", "--detach", "FETCH_HEAD"},
	}
	for _, args := range steps {
		if _, err := git(ctx, tmp, args...); err != nil {
			return "", "", fmt.Errorf("fetch %s at %s: %w", url, ref, err)
		}
	}
	commit, err = headCommit(ctx, tmp)
	if err != nil {
		return "", "", err
	}

	if err := os.Rename(tmp, dir); err != nil {
		// Another fetch of the same ref may have won the race: use its checkout.
		if info, statErr := os.Stat(dir); statErr == nil && info.IsDir() {
			commit, err = headCommit(ctx, dir)
			return dir, commit, err
		}
		return "", "", fmt.Errorf("move checkout into place: %w", err)
	}
	return dir, commit, nil
}

// ValidRef rejects refs that would not map to a single directory under the cache.
func ValidRef(ref string) error {
	switch {
	case ref == "", ref == ".", ref == "..":
		return fmt.Errorf("invalid ref %q", ref)
	case strings.HasPrefix(ref, "-"), strings.HasPrefix(ref, "."):
		return fmt.Errorf("invalid ref %q: must not start with %q", ref, ref[:1])
	case strings.ContainsAny(ref, `/\`):
		return fmt.Errorf("invalid ref %q: must not contain a path separator", ref)
	}
	return nil
}

func headCommit(ctx context.Context, dir string) (string, error) {
	out, err := git(ctx, dir, "rev-parse", "HEAD")
	if err != nil {
		return "", fmt.Errorf("read commit of %s: %w", dir, err)
	}
	return strings.TrimSpace(out), nil
}

// git runs git in dir (or the current dir when empty) and returns stdout.
func git(ctx context.Context, dir string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	out, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return "", fmt.Errorf("git %s: %w\nstderr: %s", args[0], err, strings.TrimSpace(string(exitErr.Stderr)))
		}
		return "", fmt.Errorf("git %s: %w", args[0], err)
	}
	return string(out), nil
}

func parseSemver(tag string) ([3]int, bool) {
	m := semverTag.FindStringSubmatch(tag)
	if m == nil {
		return [3]int{}, false
	}
	var v [3]int
	for i := range v {
		n, err := strconv.Atoi(m[i+1])
		if err != nil {
			return [3]int{}, false
		}
		v[i] = n
	}
	return v, true
}

func less(a, b [3]int) bool {
	for i := range a {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
	}
	return false
}
