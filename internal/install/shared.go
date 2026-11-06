package install

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// sharedLink matches a Markdown link target: an optional directory, then a
// file name ending in .md.
var sharedLink = regexp.MustCompile(`\]\(([^)#\s]*?)([^/)#\s]+\.md)\)`)

// renderShared plans one Render per shared reference a skill links: the
// shared file's content at the linked path inside the skill. A link that
// would land outside its skill is an error.
func (s Source) renderShared(skills []Skill) ([]Render, error) {
	if len(s.Shared) == 0 {
		return nil, nil
	}
	sources := map[string][]byte{}
	for _, f := range s.Shared {
		data, err := os.ReadFile(s.path(f))
		if err != nil {
			return nil, err
		}
		sources[filepath.Base(f)] = data
	}
	targets := map[string][]byte{}
	for _, sk := range skills {
		err := filepath.WalkDir(sk.Dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() || !strings.HasSuffix(path, ".md") {
				return err
			}
			if _, shared := sources[d.Name()]; shared {
				return nil // a rendered copy, not a source of links
			}
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			for _, m := range sharedLink.FindAllSubmatch(b, -1) {
				name := string(m[2])
				content, ok := sources[name]
				if !ok {
					continue
				}
				target := filepath.Join(filepath.Dir(path), string(m[1]), name)
				if !strings.HasPrefix(target, sk.Dir+string(filepath.Separator)) {
					return fmt.Errorf("skill %s: %s links %s outside the skill", sk.Name, path, name)
				}
				targets[target] = content
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	paths := make([]string, 0, len(targets))
	for p := range targets {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	renders := make([]Render, 0, len(paths))
	for _, p := range paths {
		cur, err := os.ReadFile(p)
		renders = append(renders, Render{Path: p, Content: targets[p], Exists: err == nil && bytes.Equal(cur, targets[p])})
	}
	return renders, nil
}
