package catalog

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// TestSkillLinksFollowSets checks that a skill only names skills from its own
// set or from sets that set requires, directly or not. A name followed by
// "when installed" in the same clause ("`tdd` when installed", "the `grill`
// skill when installed") may point anywhere: it is a soft link the skill can
// do without.
func TestSkillLinksFollowSets(t *testing.T) {
	root, err := filepath.Abs("../../config")
	if err != nil {
		t.Fatal(err)
	}
	c, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	shared, err := filepath.Glob(filepath.Join(root, SharedDir, "*.md"))
	if err != nil {
		t.Fatal(err)
	}
	copies := map[string]bool{}
	for _, f := range shared {
		copies[filepath.Base(f)] = true
	}

	setOf := map[string]string{} // skill name -> set
	dirOf := map[string]string{} // skill name -> directory
	for _, name := range c.Names() {
		for _, dir := range c.Sets[name].Skills {
			for skill, d := range skillsIn(t, dir) {
				if prev, ok := setOf[skill]; ok {
					t.Fatalf("skill %s is in both %s and %s", skill, prev, name)
				}
				setOf[skill], dirOf[skill] = name, d
			}
		}
	}
	if len(setOf) == 0 {
		t.Fatal("no skills found in any set")
	}

	for skill, dir := range dirOf {
		reach := c.reachable(setOf[skill])
		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() || filepath.Ext(path) != ".md" || copies[d.Name()] {
				return err
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			rel, _ := filepath.Rel(root, path)
			for _, s := range sentences(string(data)) {
				for _, other := range names(setOf) {
					if reach[setOf[other]] || !hardMention(s, other) {
						continue
					}
					t.Errorf("%s names %s (set %s), which set %s does not require; add \"when installed\" or move a skill:\n    %s", rel, other, setOf[other], setOf[skill], strings.TrimSpace(s))
				}
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}

// reachable lists set name and every set it requires, directly or not.
func (c *Catalog) reachable(name string) map[string]bool {
	seen := map[string]bool{}
	var visit func(string)
	visit = func(n string) {
		if seen[n] {
			return
		}
		seen[n] = true
		for _, r := range c.Sets[n].Requires {
			visit(r)
		}
	}
	visit(name)
	return seen
}

var frontmatterNameRE = regexp.MustCompile(`(?m)^name:\s*(\S+)\s*$`)

// skillsIn maps each skill name under dir to its directory.
func skillsIn(t *testing.T, dir string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || d.Name() != "SKILL.md" {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if m := frontmatterNameRE.FindSubmatch(data); m != nil {
			out[string(m[1])] = filepath.Dir(path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func names(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// sentences splits Markdown into lines, then lines into sentences.
func sentences(text string) []string {
	var out []string
	for _, line := range strings.Split(text, "\n") {
		start := 0
		for i := 0; i+1 < len(line); i++ {
			if strings.ContainsRune(".!?", rune(line[i])) && line[i+1] == ' ' {
				out = append(out, line[start:i+1])
				start = i + 2
			}
		}
		out = append(out, line[start:])
	}
	return out
}

var softRE = regexp.MustCompile(`^[^.;,]{0,40}?when installed`)

// hardMention reports whether s names skill as `skill`, /skill or a link into
// its directory, without "when installed" right after it.
func hardMention(s, skill string) bool {
	q := regexp.QuoteMeta(skill)
	re := regexp.MustCompile("`" + q + "`" + `|(?:^|[\s("'` + "`" + `])/` + q + `(?:$|[^\w/-])|\]\([^)]*/` + q + `/[^)]*\)`)
	for _, m := range re.FindAllStringIndex(s, -1) {
		if !softRE.MatchString(s[m[1]:]) {
			return true
		}
	}
	return false
}
