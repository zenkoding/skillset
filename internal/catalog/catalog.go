// Package catalog reads sets from a catalog directory and turns a selection
// of them into an install.Source.
//
// A set is a directory sets/<name>/ with a set.toml manifest. Paths in the
// manifest are relative to the set directory and must stay inside the
// catalog. Unset paths default to the set layout:
//
//	sets/<name>/
//	  set.toml        name, description, hosts, requires, paths
//	  shared.md       instruction fragment for every host
//	  <host>.md       optional, one host only
//	  settings.json   settings fragment (Claude)
//	  skills/...
package catalog

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"

	"github.com/zenkoding/skillset/internal/fragment"
	"github.com/zenkoding/skillset/internal/install"
)

// SetsDir holds the sets, relative to the catalog root.
const SetsDir = "sets"

// Manifest is a set's set.toml file name.
const Manifest = "set.toml"

// Hosts are the hosts a set can support.
var Hosts = []string{"claude", "codex"}

// Set is one set's manifest with every path made absolute. Optional paths
// whose default file is missing are left empty.
type Set struct {
	Name         string            `toml:"name"`
	Description  string            `toml:"description"`
	Hosts        []string          `toml:"hosts"`
	Requires     []string          `toml:"requires"`
	Shared       string            `toml:"shared"`       // default shared.md
	Instructions map[string]string `toml:"instructions"` // host -> file, default <host>.md
	Skills       []string          `toml:"skills"`       // default skills
	Settings     string            `toml:"settings"`     // default settings.json
	Files        []File            `toml:"files"`
	Dir          string            `toml:"-"`
}

// File is a file or directory the set links into HOME.
type File struct {
	Home  string   `toml:"home"` // HOME-relative link path
	Path  string   `toml:"path"`
	Hosts []string `toml:"hosts"` // installed only for these hosts; empty means always
}

// Supports reports whether the set lists host.
func (s *Set) Supports(host string) bool {
	for _, h := range s.Hosts {
		if h == host {
			return true
		}
	}
	return false
}

// Catalog is a directory of sets.
type Catalog struct {
	Root string
	Sets map[string]*Set
}

// Names lists the catalog's set names, sorted.
func (c *Catalog) Names() []string {
	names := make([]string, 0, len(c.Sets))
	for n := range c.Sets {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// Load reads every root/sets/*/set.toml.
func Load(root string) (*Catalog, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(filepath.Join(root, SetsDir))
	if err != nil {
		return nil, err
	}
	c := &Catalog{Root: root, Sets: map[string]*Set{}}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		s, err := LoadSet(root, filepath.Join(root, SetsDir, e.Name()))
		if err != nil {
			return nil, err
		}
		if s.Name != e.Name() {
			return nil, fmt.Errorf("set %s: name %q differs from directory", e.Name(), s.Name)
		}
		c.Sets[s.Name] = s
	}
	for _, s := range c.Sets {
		for _, r := range s.Requires {
			if c.Sets[r] == nil {
				return nil, fmt.Errorf("set %s requires unknown set %q", s.Name, r)
			}
		}
	}
	return c, nil
}

// LoadSet reads dir/set.toml. Paths must stay inside root.
func LoadSet(root, dir string) (*Set, error) {
	s := &Set{Dir: dir}
	md, err := toml.DecodeFile(filepath.Join(dir, Manifest), s)
	if err != nil {
		return nil, err
	}
	if keys := md.Undecoded(); len(keys) > 0 {
		return nil, fmt.Errorf("%s: unknown keys %v", filepath.Join(dir, Manifest), keys)
	}
	return s, s.resolve(root)
}

// resolve checks the manifest and makes its paths absolute.
func (s *Set) resolve(root string) (err error) {
	dir := s.Dir
	fail := func(format string, args ...any) error {
		return fmt.Errorf("set %s: "+format, append([]any{s.Name}, args...)...)
	}
	if s.Name == "" {
		return fmt.Errorf("%s: missing name", filepath.Join(dir, Manifest))
	}
	for _, h := range s.Hosts {
		if !known(h) {
			return fail("unknown host %q", h)
		}
	}
	for h := range s.Instructions {
		if !s.Supports(h) {
			return fail("instructions for %q, which is not in hosts", h)
		}
	}
	realRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		return err
	}
	// path resolves p against the set dir. It must stay inside the catalog,
	// symlinks included. A missing explicit path is an error; a missing
	// default is dropped.
	path := func(p, def string) (string, error) {
		explicit := p != ""
		if !explicit {
			p = def
		}
		abs := filepath.Join(dir, p)
		if !within(abs, root) {
			return "", fail("path %s leaves the catalog", p)
		}
		if _, err := os.Stat(abs); errors.Is(err, fs.ErrNotExist) {
			if explicit {
				return "", fail("%s does not exist", p)
			}
			return "", nil
		} else if err != nil {
			return "", err
		}
		if real, err := filepath.EvalSymlinks(abs); err != nil {
			return "", err
		} else if !within(real, realRoot) {
			return "", fail("path %s leaves the catalog through a symlink", p)
		}
		return abs, nil
	}
	if s.Shared, err = path(s.Shared, "shared.md"); err != nil {
		return err
	}
	ins := map[string]string{}
	for _, h := range s.Hosts {
		p, err := path(s.Instructions[h], h+".md")
		if err != nil {
			return err
		}
		if p != "" {
			ins[h] = p
		}
	}
	s.Instructions = ins
	skills := s.Skills
	if len(skills) == 0 {
		skills = []string{""}
	}
	s.Skills = nil
	for _, sk := range skills {
		p, err := path(sk, "skills")
		if err != nil {
			return err
		}
		if p != "" {
			s.Skills = append(s.Skills, p)
		}
	}
	if s.Settings, err = path(s.Settings, "settings.json"); err != nil {
		return err
	}
	for i, f := range s.Files {
		if f.Home == "" || f.Path == "" {
			return fail("files entry %d needs home and path", i+1)
		}
		if !homeAllowed(f.Home) {
			return fail("files entry %d: home %s must be a relative path under %s", i+1, f.Home, strings.Join(HomeDirs, ", "))
		}
		if s.Files[i].Path, err = path(f.Path, ""); err != nil {
			return err
		}
		for _, h := range f.Hosts {
			if !known(h) {
				return fail("files entry %d: unknown host %q", i+1, h)
			}
		}
	}
	return nil
}

// HomeDirs are the HOME-relative directories a set's files may be linked
// into. A catalog fetched from a URL must not reach anything else, such as
// ~/.ssh.
var HomeDirs = []string{".claude/", ".codex/", ".agents/", ".config/"}

func homeAllowed(home string) bool {
	if !filepath.IsLocal(home) {
		return false // absolute, empty or escaping with ..
	}
	clean := filepath.ToSlash(filepath.Clean(home))
	for _, d := range HomeDirs {
		if strings.HasPrefix(clean, d) && clean != d {
			return true
		}
	}
	return false
}

// within reports whether path is root or inside it.
func within(path, root string) bool {
	return path == root || strings.HasPrefix(path, root+string(filepath.Separator))
}

func known(host string) bool {
	for _, h := range Hosts {
		if h == host {
			return true
		}
	}
	return false
}

// Resolve returns the named sets plus everything they require, each after
// the sets it requires; ties keep name order.
func (c *Catalog) Resolve(names []string) ([]*Set, error) {
	var out []*Set
	const (
		visiting = 1
		done     = 2
	)
	state := map[string]int{}
	var visit func(name string, path []string) error
	visit = func(name string, path []string) error {
		s := c.Sets[name]
		if s == nil {
			return fmt.Errorf("unknown set %q", name)
		}
		switch state[name] {
		case done:
			return nil
		case visiting:
			return fmt.Errorf("requires cycle: %s", strings.Join(append(path, name), " -> "))
		}
		state[name] = visiting
		reqs := append([]string(nil), s.Requires...)
		sort.Strings(reqs)
		for _, r := range reqs {
			if err := visit(r, append(path, name)); err != nil {
				return err
			}
		}
		state[name] = done
		out = append(out, s)
		return nil
	}
	sorted := append([]string(nil), names...)
	sort.Strings(sorted)
	for _, n := range sorted {
		if err := visit(n, nil); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// HostSkillDirs are the HOME-relative skill directories each host reads.
var HostSkillDirs = map[string][]string{
	"claude": {".claude/skills"},
	"codex":  {".codex/skills", ".agents/skills"},
}

// Source turns resolved sets into what the installer installs for hosts:
// per host, each set's shared part then its host part, for sets that list
// the host. Sets outside the catalog (the local set) are owned roots too.
// Settings fragments merge; two sets setting the same key is an error.
func (c *Catalog) Source(sets []*Set, hosts []string) (install.Source, error) {
	src := install.Source{Root: c.Root, Instructions: map[string][]string{}, SkillDirs: []string{}}
	want := map[string]bool{}
	for _, h := range hosts {
		if !known(h) {
			return src, fmt.Errorf("unknown host %q", h)
		}
		want[h] = true
	}
	for _, h := range Hosts {
		if want[h] {
			src.SkillDirs = append(src.SkillDirs, HostSkillDirs[h]...)
		}
	}
	var frags []fragment.Fragment
	var first string
	for _, s := range sets {
		if !strings.HasPrefix(s.Dir, c.Root+string(filepath.Separator)) {
			src.Roots = append(src.Roots, s.Dir)
		}
		for _, h := range Hosts {
			if !want[h] || !s.Supports(h) {
				continue
			}
			if s.Shared != "" {
				src.Instructions[h] = append(src.Instructions[h], s.Shared)
			}
			if p := s.Instructions[h]; p != "" {
				src.Instructions[h] = append(src.Instructions[h], p)
			}
		}
		src.SkillRoots = append(src.SkillRoots, s.Skills...)
		for _, f := range s.Files {
			if !anyOf(f.Hosts, want) {
				continue
			}
			src.Configs = append(src.Configs, install.Config{Home: f.Home, Repo: f.Path})
		}
		if s.Settings != "" && s.Supports("claude") && want["claude"] {
			data, err := decodeObject(s.Settings)
			if err != nil {
				return src, fmt.Errorf("set %s: %w", s.Name, err)
			}
			if first == "" {
				first = s.Settings
			}
			frags = append(frags, fragment.Fragment{Set: s.Name, Data: data})
		}
	}
	switch len(frags) {
	case 0:
	case 1: // installed as is, keeping its formatting
		src.Copies = append(src.Copies, install.Config{Home: settingsHome, Repo: first})
	default:
		merged, err := fragment.Merge(frags)
		if err != nil {
			return src, err
		}
		content, err := encodeJSON(merged)
		if err != nil {
			return src, err
		}
		src.Copies = append(src.Copies, install.Config{Home: settingsHome, Repo: first, Content: content})
	}
	return src, nil
}

// settingsHome is where settings fragments are installed, relative to HOME.
const settingsHome = ".claude/settings.json"

func decodeObject(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var v map[string]any
	d := json.NewDecoder(bytes.NewReader(data))
	d.UseNumber()
	if err := d.Decode(&v); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return v, nil
}

func encodeJSON(v any) ([]byte, error) {
	var b bytes.Buffer
	e := json.NewEncoder(&b)
	e.SetEscapeHTML(false)
	e.SetIndent("", "  ")
	if err := e.Encode(v); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// anyOf reports whether hosts is empty or holds a wanted host.
func anyOf(hosts []string, want map[string]bool) bool {
	if len(hosts) == 0 {
		return true
	}
	for _, h := range hosts {
		if want[h] {
			return true
		}
	}
	return false
}
