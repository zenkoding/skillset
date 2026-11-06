package catalog

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zenkoding/skillset/internal/install"
)

// TestSharedReferences plans an install of the real catalog and checks that
// every config/shared file reaches at least one skill, and that .gitignore
// keeps the copies the installer writes out of the repo.
func TestSharedReferences(t *testing.T) {
	root, err := filepath.Abs("../../config")
	if err != nil {
		t.Fatal(err)
	}
	c, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	sets, err := c.Resolve(c.Names())
	if err != nil {
		t.Fatal(err)
	}
	src, err := c.Source(sets, Hosts)
	if err != nil {
		t.Fatal(err)
	}
	if len(src.Shared) == 0 {
		t.Fatal("no shared references in config/shared")
	}
	plan, err := install.MakeSourcePlan(src, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	reached := map[string]bool{}
	for _, r := range plan.Renders {
		if strings.HasPrefix(r.Path, filepath.Join(root, "skills")+string(filepath.Separator)) {
			reached[filepath.Base(r.Path)] = true
		}
	}
	attrs, err := os.ReadFile(filepath.Join(root, "..", ".gitignore"))
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range src.Shared {
		name := filepath.Base(f)
		if !reached[name] {
			t.Errorf("no skill links shared/%s", name)
		}
		if want := "config/skills/**/" + name; !bytes.Contains(attrs, []byte(want+"\n")) {
			t.Errorf(".gitignore lacks %q, so installed copies of %s would show as untracked", want, name)
		}
	}
}
