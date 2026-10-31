package catalog

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// sharedLink matches a Markdown link target to a shared reference file.
var sharedLink = regexp.MustCompile(`\]\(([^)#\s]*?)([^/)#\s]+\.md)\)`)

// TestSharedReferencesInSync checks that every copy of a config/shared file
// inside a skill matches its source, and that every skill linking one carries
// the copy at the linked path. Skills install one by one, so each needs its
// own copy; run `make sync/shared` after editing a source.
func TestSharedReferencesInSync(t *testing.T) {
	root, err := filepath.Abs("../../config")
	if err != nil {
		t.Fatal(err)
	}
	sharedDir := filepath.Join(root, "shared")
	entries, err := os.ReadDir(sharedDir)
	if err != nil {
		t.Fatal(err)
	}
	sources := map[string][]byte{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		b, err := os.ReadFile(filepath.Join(sharedDir, e.Name()))
		if err != nil {
			t.Fatal(err)
		}
		sources[e.Name()] = b
	}
	if len(sources) == 0 {
		t.Fatal("no shared sources in config/shared")
	}

	copies := 0
	skills := filepath.Join(root, "skills")
	err = filepath.WalkDir(skills, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".md") {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		parts := strings.SplitN(filepath.ToSlash(rel), "/", 4) // skills/<group>/<skill>/...
		if len(parts) < 4 {
			return nil
		}
		skillDir := filepath.Join(root, parts[0], parts[1], parts[2])
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if src, ok := sources[d.Name()]; ok {
			copies++
			if !bytes.Equal(b, src) {
				t.Errorf("%s differs from shared/%s; run make sync/shared", rel, d.Name())
			}
		}
		for _, m := range sharedLink.FindAllSubmatch(b, -1) {
			name := string(m[2])
			if _, ok := sources[name]; !ok {
				continue
			}
			target := filepath.Join(filepath.Dir(path), string(m[1]), name)
			if !strings.HasPrefix(target, skillDir+string(filepath.Separator)) {
				t.Errorf("%s links %s outside its skill", rel, name)
				continue
			}
			if _, err := os.Stat(target); err != nil {
				t.Errorf("%s links %s%s, which is missing; copy shared/%s there", rel, m[1], name, name)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if copies == 0 {
		t.Error("no skill carries a shared reference")
	}
}
