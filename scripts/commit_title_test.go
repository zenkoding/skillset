package scripts

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCheckCommitTitle(t *testing.T) {
	tests := []struct {
		name  string
		title string
		want  bool
	}{
		{name: "type only", title: "fix: handle an empty config", want: true},
		{name: "type and scope", title: "feat(skills): add guide", want: true},
		{name: "breaking change", title: "refactor(make)!: drop the plan target", want: true},
		{name: "scope with slash", title: "docs(contributing/issues): trim", want: true},
		{name: "empty", title: "", want: false},
		{name: "no type", title: "Add the engineering workflow skills", want: false},
		{name: "unknown type", title: "feature: add guide", want: false},
		{name: "uppercase type", title: "Fix: handle an empty config", want: false},
		{name: "no space after colon", title: "fix:handle an empty config", want: false},
		{name: "empty scope", title: "fix(): handle an empty config", want: false},
		{name: "no summary", title: "fix: ", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := exec.Command("./check-commit-title.sh", tt.title).Run()
			if got := err == nil; got != tt.want {
				t.Errorf("check-commit-title.sh %q passed = %v, want %v (err: %v)", tt.title, got, tt.want, err)
			}
		})
	}
}

func TestCommitMsgHook(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    bool
	}{
		{name: "conventional", message: "feat(skills): add guide\n\nWhat.\n\nWhy.\n", want: true},
		{name: "title after template comments", message: "\n# a comment\n\nfix: handle an empty config\n# more\n", want: true},
		{name: "not conventional", message: "Add guide\n", want: false},
		{name: "bad title after comments", message: "# feat: only a comment\nAdd guide\n", want: false},
		{name: "only comments", message: "\n# feat: only a comment\n", want: true},
		{name: "merge commit", message: "Merge branch 'main' into topic\n", want: true},
		{name: "git revert", message: "Revert \"feat: add guide\"\n", want: true},
		{name: "fixup", message: "fixup! feat: add guide\n", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := filepath.Join(t.TempDir(), "COMMIT_EDITMSG")
			if err := os.WriteFile(file, []byte(tt.message), 0o600); err != nil {
				t.Fatal(err)
			}

			err := exec.Command("../.githooks/commit-msg", file).Run()
			if got := err == nil; got != tt.want {
				t.Errorf("commit-msg on %q passed = %v, want %v (err: %v)", tt.message, got, tt.want, err)
			}
		})
	}
}
