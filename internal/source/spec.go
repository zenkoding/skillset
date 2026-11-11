package source

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// GitHub is the base URL that owner/repo shorthands expand to.
const GitHub = "https://github.com/"

var (
	shorthand = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.-]*/[A-Za-z0-9_.-]+$`)
	fullSHA   = regexp.MustCompile(`^[0-9a-f]{40}$`)
	shortSHA  = regexp.MustCompile(`^[0-9a-f]{4,39}$`)
)

// ErrRefNotFound is returned, wrapped, when a ref names no branch or tag.
var ErrRefNotFound = errors.New("not found")

// Parse splits a git source into its URL and ref. A source is a git URL
// (scheme://host/path or the scp-like host:path) or an owner/repo shorthand
// for GitHub, optionally followed by @ref. The ref starts at the first @
// after the URL's path begins, so user@host parts stay in the URL and a
// branch name may contain slashes. ref is "" when none is given.
func Parse(arg string) (url, ref string, err error) {
	if arg == "" {
		return "", "", errors.New("empty source")
	}
	if strings.HasPrefix(arg, "-") {
		return "", "", fmt.Errorf("%q is not a source: want a directory, a git URL or owner/repo", arg)
	}
	pathStart := 0
	if i := strings.Index(arg, "://"); i > 0 {
		j := strings.Index(arg[i+3:], "/")
		if j < 0 {
			return "", "", fmt.Errorf("%q is not a source: the URL has no path", arg)
		}
		pathStart = i + 3 + j
		// The config is written to disk and, at repo scope, committed, so
		// a URL must not carry a credential: a password, or any user part
		// in HTTP(S), where it is a token. The error leaves the URL out.
		if at := strings.LastIndex(arg[i+3:pathStart], "@"); at >= 0 {
			user := arg[i+3 : i+3+at]
			scheme := strings.ToLower(arg[:i])
			if strings.Contains(user, ":") || scheme == "http" || scheme == "https" {
				return "", "", errors.New("the source URL holds a credential: remove it and let git's credential helper log in")
			}
		}
	} else if c := strings.Index(arg, ":"); c > 0 && !strings.Contains(arg[:c], "/") {
		pathStart = c + 1
	}
	url = arg
	if k := strings.Index(arg[pathStart:], "@"); k >= 0 {
		url, ref = arg[:pathStart+k], arg[pathStart+k+1:]
		if ref == "" {
			return "", "", fmt.Errorf("%q: empty ref after @", arg)
		}
		if err := validRef(ref); err != nil {
			return "", "", err
		}
	}
	if pathStart == 0 {
		if !shorthand.MatchString(url) {
			return "", "", fmt.Errorf("%q is not a source: want a directory, a git URL or owner/repo", arg)
		}
		url = GitHub + url
	}
	return url, ref, nil
}

// Name is the repository name in url: its last path element without .git.
func Name(url string) string {
	url = strings.TrimRight(url, "/")
	if i := strings.LastIndexAny(url, "/:"); i >= 0 {
		url = url[i+1:]
	}
	return strings.TrimSuffix(url, ".git")
}

// validRef rejects refs git would refuse and anything git could read as an
// option. Unlike ValidRef it allows slashes: the ref never names a path.
func validRef(ref string) error {
	bad := ref == "" || strings.HasPrefix(ref, "-") || strings.HasPrefix(ref, "/") ||
		strings.HasSuffix(ref, "/") || strings.HasSuffix(ref, ".lock") ||
		strings.Contains(ref, "..") || strings.Contains(ref, "//") || strings.Contains(ref, "@{") ||
		strings.ContainsAny(ref, " ~^:?*[\\\x7f") || strings.IndexFunc(ref, func(r rune) bool { return r < ' ' }) >= 0
	if bad {
		return fmt.Errorf("invalid ref %q", ref)
	}
	return nil
}

// Resolve returns the commit ref points at in url: a branch head, a tag's
// commit, or a full commit hash as given. A name that is both a branch and
// a tag is an error.
func Resolve(ctx context.Context, url, ref string) (string, error) {
	if err := validRef(ref); err != nil {
		return "", err
	}
	if fullSHA.MatchString(ref) {
		return ref, nil
	}
	head, tag := "refs/heads/"+ref, "refs/tags/"+ref
	out, err := git(ctx, "", "ls-remote", "--", url, head, tag, tag+"^{}")
	if err != nil {
		return "", fmt.Errorf("resolve %s at %s: %w", url, ref, err)
	}
	found := map[string]string{}
	for _, line := range strings.Split(out, "\n") {
		if commit, name, ok := strings.Cut(strings.TrimSpace(line), "\t"); ok {
			found[name] = commit
		}
	}
	branch, isBranch := found[head]
	tagCommit, isTag := found[tag+"^{}"] // an annotated tag, peeled
	if !isTag {
		tagCommit, isTag = found[tag]
	}
	switch {
	case isBranch && isTag:
		return "", fmt.Errorf("ref %q in %s is both a branch and a tag: give the commit hash", ref, url)
	case isBranch:
		return branch, nil
	case isTag:
		return tagCommit, nil
	case shortSHA.MatchString(ref):
		return "", fmt.Errorf("ref %q %w in %s: give a commit as its full commit hash", ref, ErrRefNotFound, url)
	}
	return "", fmt.Errorf("ref %q %w in %s", ref, ErrRefNotFound, url)
}

// FetchCommit makes cacheDir/<commit> hold url at commit, a full commit
// hash, and returns that dir. A commit never changes, so a cached checkout
// is reused without the network.
func FetchCommit(ctx context.Context, url, commit, cacheDir string) (string, error) {
	if !fullSHA.MatchString(commit) {
		return "", fmt.Errorf("fetch %s: %q is not a full commit hash", url, commit)
	}
	dir, got, err := checkout(ctx, url, commit, cacheDir, commit)
	if err != nil {
		return "", err
	}
	if got != commit {
		return "", fmt.Errorf("fetch %s at %s: got commit %s", url, commit, got)
	}
	return dir, nil
}
