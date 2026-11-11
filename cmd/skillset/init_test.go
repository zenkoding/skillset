package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/zenkoding/skillset/internal/catalog"
)

// sourceRepo makes a git repository at dir holding a catalog with sets
// base, dcg (requires base) and other. It returns the commits of v1.0.0
// (an annotated tag), the branch feature/x and main.
func sourceRepo(t *testing.T, dir string) (tag, branch, main string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	gitRun(t, dir, "init", "-q", "-b", "main")
	for p, body := range map[string]string{
		"config/sets/base/set.toml":  "name = \"base\"\nhosts = [\"claude\", \"codex\"]",
		"config/sets/base/shared.md": "base\n",
		"config/sets/dcg/set.toml":   "name = \"dcg\"\nhosts = [\"claude\", \"codex\"]\nrequires = [\"base\"]",
		"config/sets/dcg/shared.md":  "dcg\n",
	} {
		writeFile(t, filepath.Join(dir, p), body)
	}
	gitRun(t, dir, "add", "-A")
	gitRun(t, dir, "commit", "-qm", "one")
	gitRun(t, dir, "tag", "-a", "-m", "v1.0.0", "v1.0.0")
	tag = revParse(t, dir, "HEAD")
	writeFile(t, filepath.Join(dir, "config/sets/other/set.toml"), "name = \"other\"\nhosts = [\"claude\"]")
	writeFile(t, filepath.Join(dir, "config/sets/other/shared.md"), "other\n")
	gitRun(t, dir, "add", "-A")
	gitRun(t, dir, "commit", "-qm", "two")
	gitRun(t, dir, "branch", "feature/x")
	branch = revParse(t, dir, "HEAD")
	writeFile(t, filepath.Join(dir, "config/sets/base/shared.md"), "base 3\n")
	gitRun(t, dir, "commit", "-qam", "three")
	return tag, branch, revParse(t, dir, "HEAD")
}

func revParse(t *testing.T, dir, rev string) string {
	t.Helper()
	out, err := exec.Command("git", "-C", dir, "rev-parse", rev).Output()
	if err != nil {
		t.Fatal(err)
	}
	return strings.TrimSpace(string(out))
}

// isolateGit hides the global and system git config, whose insteadOf
// and signing rules would change what the tests fetch and commit, and
// lets git use only file:// remotes, so a test never reaches the network.
func isolateGit(t *testing.T) {
	t.Helper()
	t.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_ALLOW_PROTOCOL", "file")
}

// insteadOf makes git fetch every GitHub URL, HTTPS or SSH, from local
// repositories under root, so tests of the shorthand and the default
// source need no network.
func insteadOf(t *testing.T, root string) {
	t.Helper()
	t.Setenv("GIT_CONFIG_COUNT", "2")
	t.Setenv("GIT_CONFIG_KEY_0", "url.file://"+root+"/.insteadOf")
	t.Setenv("GIT_CONFIG_VALUE_0", "https://github.com/")
	t.Setenv("GIT_CONFIG_KEY_1", "url.file://"+root+"/.insteadOf")
	t.Setenv("GIT_CONFIG_VALUE_1", "git@github.com:")
}

func TestInitSources(t *testing.T) {
	isolateGit(t)
	github := t.TempDir()
	insteadOf(t, github)
	repo := filepath.Join(github, "zenkoding", "skillset")
	tag, branch, main := sourceRepo(t, repo)
	url := "file://" + repo
	// A repository whose only branch is trunk, not main.
	trunk := filepath.Join(github, "o", "trunk")
	if err := os.MkdirAll(trunk, 0o755); err != nil {
		t.Fatal(err)
	}
	gitRun(t, trunk, "init", "-q", "-b", "trunk")
	gitRun(t, trunk, "commit", "-q", "--allow-empty", "-m", "one")

	tests := []struct {
		name      string
		saved     string // config body before init, "" for none
		args      []string
		savedSets bool // pass no --set, so init starts from the saved sets
		want      catalog.Source
		commit    string // printed commit, "" for a local source
		error     string
	}{
		{name: "default", args: nil, want: catalog.Source{Name: "skillset", URL: "https://github.com/zenkoding/skillset", Ref: "main"}, commit: main},
		{name: "tag", args: []string{url + "@v1.0.0"}, want: catalog.Source{Name: "skillset", URL: url, Ref: "v1.0.0"}, commit: tag},
		{name: "commit", args: []string{url + "@" + branch}, want: catalog.Source{Name: "skillset", URL: url, Ref: branch}, commit: branch},
		{name: "branch with a slash", args: []string{url + "@feature/x"}, want: catalog.Source{Name: "skillset", URL: url, Ref: "feature/x"}, commit: branch},
		{name: "url without ref", args: []string{url}, want: catalog.Source{Name: "skillset", URL: url, Ref: "main"}, commit: main},
		{name: "shorthand", args: []string{"zenkoding/skillset@feature/x"}, want: catalog.Source{Name: "skillset", URL: "https://github.com/zenkoding/skillset", Ref: "feature/x"}, commit: branch},
		{name: "scp-like", args: []string{"git@github.com:zenkoding/skillset@v1.0.0"}, want: catalog.Source{Name: "skillset", URL: "git@github.com:zenkoding/skillset", Ref: "v1.0.0"}, commit: tag},
		{name: "local directory", args: []string{repo}, want: catalog.Source{Name: "skillset", Path: repo}},
		{
			name:   "keeps the saved source",
			saved:  "schema = 1\nhosts = [\"claude\"]\nsets = [\"base\"]\n[[sources]]\nname = \"skillset\"\nurl = \"" + url + "\"\nref = \"feature/x\"\n",
			want:   catalog.Source{Name: "skillset", URL: url, Ref: "feature/x"},
			commit: branch,
		},
		{
			// other is not at v1.0.0: a new source drops the saved sets it lacks.
			name:      "new source drops saved sets it lacks",
			saved:     "schema = 1\nhosts = [\"claude\"]\nsets = [\"dcg\", \"other\"]\n[[sources]]\nname = \"skillset\"\nurl = \"" + url + "\"\nref = \"main\"\n",
			args:      []string{url + "@v1.0.0"},
			savedSets: true,
			want:      catalog.Source{Name: "skillset", URL: url, Ref: "v1.0.0"},
			commit:    tag,
		},
		{
			name:      "saved source keeps a missing set an error",
			saved:     "schema = 1\nhosts = [\"claude\"]\nsets = [\"other\"]\n[[sources]]\nname = \"skillset\"\nurl = \"" + url + "\"\nref = \"v1.0.0\"\n",
			savedSets: true,
			error:     `unknown set "other"`,
		},
		{name: "main assumed", args: []string{"file://" + trunk}, error: "no @ref given, so main was assumed"},
		{name: "unknown ref", args: []string{url + "@nope"}, error: `ref "nope" not found`},
		{name: "short commit", args: []string{url + "@" + branch[:7]}, error: "full commit hash"},
		{name: "set missing at the ref", args: []string{url + "@v1.0.0", "--set", "other"}, error: `unknown set "other"`},
		{name: "not a source", args: []string{"nonsense"}, error: "not a source"},
		{name: "catalog flag", args: []string{"--catalog", filepath.Join(repo, "config")}, error: "skillset init <source>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			home := t.TempDir()
			if tt.saved != "" {
				writeFile(t, catalog.ConfigPath(home), tt.saved)
			}
			args := append([]string{"init", "--home", home, "--host", "claude"}, tt.args...)
			if !tt.savedSets && !contains(tt.args, "--set") {
				args = append(args, "--set", "dcg")
			}
			out, err := run(t, args...)
			if tt.error != "" {
				if err == nil || !strings.Contains(err.Error(), tt.error) {
					t.Fatalf("error = %v, want %q", err, tt.error)
				}
				if tt.saved == "" {
					if _, statErr := os.Stat(catalog.ConfigPath(home)); statErr == nil {
						t.Fatal("a failed init wrote a config")
					}
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			cfg, ok, err := catalog.ReadConfig(home)
			if err != nil || !ok {
				t.Fatalf("ReadConfig = %v, %v", ok, err)
			}
			// Only the chosen set: install adds what it requires.
			want := catalog.Config{Schema: 1, Hosts: []string{"claude"}, Sets: []string{"dcg"}, Sources: []catalog.Source{tt.want}}
			if !reflect.DeepEqual(cfg, want) {
				t.Fatalf("config = %+v\nwant %+v", cfg, want)
			}
			if tt.commit != "" && !strings.Contains(out, tt.commit[:12]) {
				t.Fatalf("output does not show commit %s:\n%s", tt.commit[:12], out)
			}
		})
	}
}

func TestInitDcg(t *testing.T) {
	isolateGit(t)
	repo := t.TempDir()
	sourceRepo(t, repo)
	git, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string
		dcg  bool
		warn bool
	}{
		{name: "missing", dcg: false, warn: true},
		{name: "installed", dcg: true, warn: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// PATH holds git and fake installers that leave a marker if run.
			bin, marker := t.TempDir(), filepath.Join(t.TempDir(), "ran")
			if err := os.Symlink(git, filepath.Join(bin, "git")); err != nil {
				t.Fatal(err)
			}
			fakes := []string{"curl", "bash", "sh"}
			if tt.dcg {
				fakes = append(fakes, "dcg")
			}
			for _, name := range fakes {
				writeFile(t, filepath.Join(bin, name), "#!/bin/sh\necho "+name+" >> "+marker+"\n")
				if err := os.Chmod(filepath.Join(bin, name), 0o755); err != nil {
					t.Fatal(err)
				}
			}
			t.Setenv("PATH", bin)

			out, err := run(t, "init", repo, "--home", t.TempDir(), "--host", "claude", "--set", "base")
			if err != nil {
				t.Fatal(err)
			}
			warned := strings.Contains(out, "warning: dcg is not installed") &&
				strings.Contains(out, "destructive_command_guard/main/install.sh")
			if warned != tt.warn {
				t.Fatalf("warned = %v, want %v:\n%s", warned, tt.warn, out)
			}
			if ran, err := os.ReadFile(marker); err == nil {
				t.Fatalf("init ran %s", ran)
			}
		})
	}
}
