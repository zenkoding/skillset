package install

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAsk(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"keep local after a bad answer", "x\nl\n", `"b"`},
		{"keep repo", "r\n", `"c"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, home := t.TempDir(), t.TempDir()
			src := filepath.Join(repo, "claude/settings.json")
			dst := filepath.Join(home, ".claude/settings.json")
			if err := os.MkdirAll(filepath.Join(repo, "skills"), 0o755); err != nil {
				t.Fatal(err)
			}
			write(t, src, `{"m": "a"}`)
			install(t, repo, home)
			write(t, dst, `{"m": "b"}`)
			write(t, src, `{"m": "c"}`)
			p, err := makePlan(repo, home)
			if err != nil {
				t.Fatal(err)
			}
			var out bytes.Buffer
			if err := p.Ask(strings.NewReader(tt.input), &out); err != nil {
				t.Fatal(err)
			}
			if err := p.Apply(); err != nil {
				t.Fatal(err)
			}
			data, _ := os.ReadFile(dst)
			if !strings.Contains(string(data), tt.want) {
				t.Fatalf("settings.json = %s, want m = %s", data, tt.want)
			}
		})
	}
}
