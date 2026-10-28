package install

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GlobalIgnore is the pattern the installer adds to git's global excludes
// file, so the optional <repo>/.docs link to a repo's private project
// knowledge store is never committed.
const GlobalIgnore = ".docs/"

// Ignore is the global git ignore entry the installer ensures.
type Ignore struct {
	Path    string // git's global excludes file
	Pattern string
	Exists  bool // already listed
}

// planIgnore finds git's global excludes file for home and checks whether it
// already lists GlobalIgnore.
func planIgnore(home string) (*Ignore, error) {
	path, err := excludesFile(home)
	if err != nil {
		return nil, err
	}
	ig := &Ignore{Path: path, Pattern: GlobalIgnore}
	data, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	for _, line := range strings.Split(string(data), "\n") {
		if l := strings.TrimSpace(line); l == GlobalIgnore || l == strings.TrimSuffix(GlobalIgnore, "/") {
			ig.Exists = true
		}
	}
	return ig, nil
}

// excludesFile returns core.excludesFile from home's global git config, or
// git's default ($XDG_CONFIG_HOME/git/ignore, else ~/.config/git/ignore).
// For a home other than $HOME, the caller's XDG and global config overrides
// are dropped so only files under home count.
func excludesFile(home string) (string, error) {
	env := []string{"HOME=" + home}
	xdg := ""
	for _, kv := range os.Environ() {
		k, v, _ := strings.Cut(kv, "=")
		switch k {
		case "HOME":
			continue
		case "XDG_CONFIG_HOME", "GIT_CONFIG_GLOBAL":
			if home != os.Getenv("HOME") {
				continue
			}
			if k == "XDG_CONFIG_HOME" {
				xdg = v
			}
		}
		env = append(env, kv)
	}
	cmd := exec.Command("git", "config", "--global", "--type=path", "--get", "core.excludesFile")
	cmd.Env = env
	out, err := cmd.Output()
	var exit *exec.ExitError
	switch {
	case err == nil:
		return strings.TrimSpace(string(out)), nil // --type=path expands ~ against HOME
	case errors.As(err, &exit) && exit.ExitCode() == 1, errors.Is(err, exec.ErrNotFound):
		// Not set, or no git: git's default location.
	default:
		return "", err
	}
	if xdg == "" {
		xdg = filepath.Join(home, ".config")
	}
	return filepath.Join(xdg, "git", "ignore"), nil
}

// apply appends the pattern. It appends rather than rewrites so an excludes
// file that is a symlink (a dotfiles manager's) stays one.
func (ig *Ignore) apply() error {
	if ig.Exists {
		return nil
	}
	data, err := os.ReadFile(ig.Path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	line := ig.Pattern + "\n"
	if len(data) > 0 && !bytes.HasSuffix(data, []byte("\n")) {
		line = "\n" + line
	}
	if err := os.MkdirAll(filepath.Dir(ig.Path), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(ig.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	if _, err := f.WriteString(line); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}
