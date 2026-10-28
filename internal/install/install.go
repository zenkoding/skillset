// Package install links skills, instruction files and settings from a
// Source into each host's home directory (~/.claude, ~/.codex, ~/.agents,
// ~/.config/dcg).
//
// Every skill directory becomes one symlink per host skill directory. Drafts
// under in-progress/ are never installed. Instruction parts are rendered into
// <root>/build/ and linked from HOME. The installer records every link it
// creates in a manifest and, on later runs, removes manifest links that point
// into the source but are no longer wanted. Anything else in the way is a
// clash: planning and installing stop without changing anything. Settings
// files and directories (Source.Configs) are linked whole. JSON files the host
// writes to (Source.Copies) are copied instead. Later runs merge local edits
// with source changes key by key; a key both sides changed is a conflict that
// has to be resolved before install. The installer also adds .docs/ to git's
// global excludes file (see Ignore).
package install

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
)

// DraftDir is the skill-root subdirectory whose skills are never installed.
const DraftDir = "in-progress"

// SkillDirs are the per-host skill directories, relative to HOME.
var SkillDirs = []string{".claude/skills", ".codex/skills", ".agents/skills"}

// Instruction is one host's user-level instruction file. The installer
// renders the host's Source.Instructions parts into build/<Host>/<File> and
// links HOME/<Home> to it.
type Instruction struct {
	Home string // HOME-relative link path
	Host string // key in Source.Instructions, and build/<Host>/
	File string // rendered file name
}

// Instructions lists the instruction files for every host.
var Instructions = []Instruction{
	{Home: ".claude/CLAUDE.md", Host: "claude", File: "CLAUDE.md"},
	{Home: ".codex/AGENTS.md", Host: "codex", File: "AGENTS.md"},
}

// BuildDir holds rendered files, relative to the source root. It is gitignored.
const BuildDir = "build"

// Config is one settings file or directory linked from HOME into the repo.
// Entries whose repo path does not exist are skipped.
type Config struct {
	Home    string // HOME-relative link path
	Repo    string // repo-relative target
	Content []byte // copies only: what to install instead of reading Repo
}

// CopyStateDir keeps, per copy, what the repo held and what the installer
// wrote on the last install, relative to HOME. Merges use them as the base.
const CopyStateDir = StateDir + "/copies"

// ManifestPath is where the installer records its links, relative to HOME.
const ManifestPath = StateDir + "/installed-links.json"

// StateDir holds the installer's state, relative to HOME.
const StateDir = ".local/state/skillset"

// Link is one symlink the installer wants in place.
type Link struct {
	Path   string // absolute path under HOME
	Target string // absolute path inside the repo
	Exists bool   // already correct
}

// Plan is everything an install would do, computed without side effects.
type Plan struct {
	Renders  []Render
	Links    []Link
	Copies   []Copy
	Remove   []string // links to delete
	Kept     []string // stale entries reported and left alone
	Skipped  []string // skills not installed, with the reason
	Clashes  []string // paths in the way; non-empty blocks install
	Ignore   *Ignore
	manifest string
}

// Render is one generated file the installer writes inside the repo.
type Render struct {
	Path    string // absolute path under <repo>/build
	Content []byte
	Exists  bool // already up to date
}

// Copy is one settings file the installer writes under HOME.
type Copy struct {
	Path      string // absolute path under HOME
	Content   []byte // what install writes
	Exists    bool   // already up to date
	Conflicts []Conflict
	repo      []byte
	state     string          // state file prefix under CopyStateDir
	values    [4]any          // repo base, repo, last installed, local
	local     []byte          // nil when there is no local file to keep
	choices   map[string]bool // conflict key -> keep local
}

// Unresolved lists the conflicts no choice has settled yet.
func (c *Copy) Unresolved() []Conflict {
	var out []Conflict
	for _, k := range c.Conflicts {
		if _, ok := c.choices[k.Key]; !ok {
			out = append(out, k)
		}
	}
	return out
}

func (c *Copy) remerge() error {
	if c.local == nil {
		c.Content, c.Exists = c.repo, false
		return nil
	}
	out, conflicts := mergeJSON(c.values[0], c.values[1], c.values[2], c.values[3], c.choices)
	c.Conflicts = conflicts
	switch {
	case reflect.DeepEqual(out, c.values[1]):
		c.Content = c.repo
	case reflect.DeepEqual(out, c.values[3]):
		c.Content = c.local
	default:
		data, err := encodeJSON(out)
		if err != nil {
			return err
		}
		c.Content = data
	}
	c.Exists = bytes.Equal(c.Content, c.local)
	return nil
}

// Resolve settles a conflict on the copy at path: keepLocal keeps the local
// value, otherwise the repo value wins.
func (p *Plan) Resolve(path, key string, keepLocal bool) error {
	for i := range p.Copies {
		c := &p.Copies[i]
		if c.Path != path {
			continue
		}
		if c.choices == nil {
			c.choices = map[string]bool{}
		}
		c.choices[key] = keepLocal
		return c.remerge()
	}
	return fmt.Errorf("no copy at %s", path)
}

// planCopy plans one copy; c.Repo is absolute.
func planCopy(home string, c Config, in func(string) bool) (cp Copy, clash bool, err error) {
	cp = Copy{Path: filepath.Join(home, c.Home), state: filepath.Join(home, CopyStateDir, c.Home)}
	if cp.repo = c.Content; cp.repo == nil {
		cp.repo, err = os.ReadFile(c.Repo)
	}
	if err != nil {
		return cp, false, err
	}
	if cp.values[1], err = decodeJSON(cp.repo); err != nil {
		return cp, false, fmt.Errorf("%s: %w", c.Repo, err)
	}
	fi, err := os.Lstat(cp.Path)
	switch {
	case err != nil:
		return cp, false, cp.remerge() // nothing there yet
	case fi.Mode()&fs.ModeSymlink != 0:
		if t, _ := linkTarget(cp.Path); in(t) {
			return cp, false, cp.remerge() // our link; the copy replaces it
		}
		return cp, true, nil
	case !fi.Mode().IsRegular():
		return cp, true, nil
	}
	if cp.local, err = os.ReadFile(cp.Path); err != nil {
		return cp, false, err
	}
	if cp.values[3], err = decodeJSON(cp.local); err != nil {
		return cp, true, nil // not JSON: nothing to merge
	}
	for i, suffix := range []string{".repo", ".installed"} {
		data, err := os.ReadFile(cp.state + suffix)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return cp, false, err
		}
		if cp.values[i*2], err = decodeJSON(data); err != nil {
			cp.values[i*2] = missing{} // unreadable base: merge as a first install
		}
	}
	return cp, false, cp.remerge()
}

// Source is what a plan installs. Relative paths are relative to Root, and
// links that point under Root belong to the installer.
type Source struct {
	Root string
	// Instructions maps a host (see Instructions) to the files rendered, in
	// order, into its instruction file. Hosts without parts are skipped.
	Instructions map[string][]string
	SkillRoots   []string // directories searched for skills
	Configs      []Config // linked whole
	Copies       []Config // copied and merged
	// SkillDirs are the HOME-relative directories skills are linked into;
	// nil means SkillDirs.
	SkillDirs []string
	// Roots are more directories whose links belong to the installer.
	Roots []string
}

func (s Source) path(p string) string {
	if filepath.IsAbs(p) {
		return p
	}
	return filepath.Join(s.Root, p)
}

func (s Source) render(host string) ([]byte, error) {
	var out []byte
	for _, part := range s.Instructions[host] {
		data, err := os.ReadFile(s.path(part))
		if err != nil {
			return nil, err
		}
		out = append(out, data...)
	}
	return out, nil
}

// Skill is one installable skill found in the repo.
type Skill struct {
	Name string
	Dir  string
}

func (s Source) discover() (skills []Skill, skipped []string, err error) {
	seen := map[string]string{}
	for _, root := range s.SkillRoots {
		if skipped, err = s.discoverIn(s.path(root), seen, &skills, skipped); err != nil {
			return nil, nil, err
		}
	}
	sort.Slice(skills, func(i, j int) bool { return skills[i].Name < skills[j].Name })
	return skills, skipped, nil
}

func (s Source) discoverIn(root string, seen map[string]string, skills *[]Skill, skipped []string) ([]string, error) {
	repo := s.Root
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && path != root && d.Name() == DraftDir && filepath.Dir(path) == root {
			return filepath.SkipDir
		}
		if d.IsDir() || d.Name() != "SKILL.md" {
			return nil
		}
		dir := filepath.Dir(path)
		name := filepath.Base(dir)
		rel, _ := filepath.Rel(repo, dir)
		fm, err := frontmatterName(path)
		switch {
		case err != nil:
			return err
		case fm == "":
			skipped = append(skipped, rel+": empty or no name (draft)")
		case fm != name:
			skipped = append(skipped, fmt.Sprintf("%s: frontmatter name %q differs from directory", rel, fm))
		case seen[name] != "":
			return fmt.Errorf("duplicate skill name %q: %s and %s", name, seen[name], rel)
		default:
			seen[name] = rel
			*skills = append(*skills, Skill{Name: name, Dir: dir})
		}
		return filepath.SkipDir
	})
	return skipped, err
}

func frontmatterName(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }() // read-only; close error is irrelevant
	s := bufio.NewScanner(f)
	if !s.Scan() || strings.TrimSpace(s.Text()) != "---" {
		return "", s.Err()
	}
	for s.Scan() {
		line := s.Text()
		if strings.TrimSpace(line) == "---" {
			break
		}
		if v, ok := strings.CutPrefix(line, "name:"); ok {
			return strings.Trim(strings.TrimSpace(v), `"'`), nil
		}
	}
	return "", s.Err()
}

// MakeSourcePlan computes what Install would do for src into home.
func MakeSourcePlan(src Source, home string) (*Plan, error) {
	repo, err := filepath.Abs(src.Root)
	if err != nil {
		return nil, err
	}
	src.Root = repo
	skills, skipped, err := src.discover()
	if err != nil {
		return nil, err
	}
	p := &Plan{Skipped: skipped, manifest: filepath.Join(home, ManifestPath)}
	if p.Ignore, err = planIgnore(home); err != nil {
		return nil, err
	}
	in := intoRepo(repo, src.Roots...)
	skillDirs := src.SkillDirs
	if skillDirs == nil {
		skillDirs = SkillDirs
	}

	want := map[string]string{}
	for _, dir := range skillDirs {
		for _, s := range skills {
			want[filepath.Join(home, dir, s.Name)] = s.Dir
		}
	}
	candidates := map[string]bool{}
	for _, ins := range Instructions {
		if len(src.Instructions[ins.Host]) == 0 {
			continue
		}
		content, err := src.render(ins.Host)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, err
		}
		out := filepath.Join(repo, BuildDir, ins.Host, ins.File)
		cur, err := os.ReadFile(out)
		p.Renders = append(p.Renders, Render{Path: out, Content: content, Exists: err == nil && string(cur) == string(content)})
		want[filepath.Join(home, ins.Home)] = out
	}
	for _, c := range src.Configs {
		target := src.path(c.Repo)
		if _, err := os.Stat(target); errors.Is(err, fs.ErrNotExist) {
			continue
		} else if err != nil {
			return nil, err
		}
		path := filepath.Join(home, c.Home)
		want[path] = target
		// A link into the source already there is ours to replace.
		candidates[path] = true
	}

	// Candidates for cleanup: the manifest and the links found above.
	old, err := readManifest(p.manifest)
	if err != nil {
		return nil, err
	}
	for _, path := range old.Links {
		candidates[path] = true
	}
	removable := map[string]bool{}
	for _, c := range src.Copies {
		c.Repo = src.path(c.Repo)
		if _, err := os.Stat(c.Repo); c.Content == nil && errors.Is(err, fs.ErrNotExist) {
			continue
		}
		cp, clash, err := planCopy(home, c, in)
		if err != nil {
			return nil, err
		}
		if clash {
			p.Clashes = append(p.Clashes, cp.Path)
			continue
		}
		p.Copies = append(p.Copies, cp)
	}

	for path := range candidates {
		if target, ok := want[path]; ok && linksTo(path, target) {
			continue
		}
		switch ownedState(path, in) {
		case owned:
			removable[path] = true
			p.Remove = append(p.Remove, path)
		case foreign:
			if _, ok := want[path]; !ok {
				p.Kept = append(p.Kept, path)
			}
		}
	}

	for path, target := range want {
		l := Link{Path: path, Target: target}
		if _, err := os.Lstat(path); err == nil {
			switch {
			case linksTo(path, target):
				l.Exists = true
			case removable[path]:
			default:
				p.Clashes = append(p.Clashes, path)
			}
		}
		p.Links = append(p.Links, l)
	}
	sort.Slice(p.Links, func(i, j int) bool { return p.Links[i].Path < p.Links[j].Path })
	sort.Strings(p.Remove)
	sort.Strings(p.Kept)
	sort.Strings(p.Clashes)
	return p, nil
}

// ClashError lists every path that blocks an install.
type ClashError struct{ Paths []string }

func (e *ClashError) Error() string {
	return "refusing to overwrite existing paths:\n  " + strings.Join(e.Paths, "\n  ")
}

// ConflictError lists settings keys that still need a choice.
type ConflictError struct{ Keys []string }

func (e *ConflictError) Error() string {
	return "unresolved settings conflicts (run apply in a terminal to choose):\n  " + strings.Join(e.Keys, "\n  ")
}

// Err reports the plan's clashes or unresolved conflicts, if any.
func (p *Plan) Err() error {
	if len(p.Clashes) > 0 {
		return &ClashError{Paths: p.Clashes}
	}
	var keys []string
	for i := range p.Copies {
		for _, k := range p.Copies[i].Unresolved() {
			keys = append(keys, p.Copies[i].Path+" "+k.Key)
		}
	}
	if len(keys) > 0 {
		return &ConflictError{Keys: keys}
	}
	return nil
}

// Apply carries out the plan. It refuses to run if there are clashes.
func (p *Plan) Apply() error {
	if err := p.Err(); err != nil {
		return err
	}
	for _, r := range p.Renders {
		if r.Exists {
			continue
		}
		if err := writeAtomic(r.Path, r.Content); err != nil {
			return err
		}
	}
	for _, path := range p.Remove {
		if err := os.Remove(path); err != nil {
			return err
		}
	}
	var made []string
	for _, l := range p.Links {
		made = append(made, l.Path)
		if l.Exists {
			continue
		}
		if err := os.MkdirAll(filepath.Dir(l.Path), 0o755); err != nil {
			return err
		}
		if err := os.Symlink(l.Target, l.Path); err != nil {
			return err
		}
	}
	copies := map[string]string{}
	for _, c := range p.Copies {
		if !c.Exists {
			if err := writeAtomic(c.Path, c.Content); err != nil {
				return err
			}
		}
		if err := writeAtomic(c.state+".repo", c.repo); err != nil {
			return err
		}
		if err := writeAtomic(c.state+".installed", c.Content); err != nil {
			return err
		}
		copies[c.Path] = hash(c.Content)
	}
	if err := p.Ignore.apply(); err != nil {
		return err
	}
	return writeManifest(p.manifest, manifest{Links: made, Copies: copies})
}

// String renders the plan for humans.
func (p *Plan) String() string {
	var b strings.Builder
	section := func(title string, items []string) {
		if len(items) == 0 {
			return
		}
		fmt.Fprintf(&b, "%s:\n", title)
		for _, s := range items {
			fmt.Fprintf(&b, "  %s\n", s)
		}
	}
	var render, fresh []string
	for _, r := range p.Renders {
		if r.Exists {
			fresh = append(fresh, r.Path)
		} else {
			render = append(render, r.Path)
		}
	}
	section("render (inside the repo)", render)
	section("rendered, up to date", fresh)
	var create, ok []string
	for _, l := range p.Links {
		s := l.Path + " -> " + l.Target
		if l.Exists {
			ok = append(ok, s)
		} else {
			create = append(create, s)
		}
	}
	section("link", create)
	section("unchanged", ok)
	var copied, same []string
	for _, c := range p.Copies {
		if c.Exists {
			same = append(same, c.Path)
		} else {
			copied = append(copied, c.Path)
		}
	}
	section("copy", copied)
	section("copied, up to date", same)
	var conflicts []string
	for _, c := range p.Copies {
		for _, k := range c.Conflicts {
			conflicts = append(conflicts, fmt.Sprintf("%s %s: repo %s, local %s", c.Path, k.Key, k.Repo, k.Local))
		}
	}
	section("CONFLICT (choose at install)", conflicts)
	if ig := p.Ignore; ig != nil {
		if ig.Exists {
			section("git ignore, up to date", []string{ig.Pattern + " in " + ig.Path})
		} else {
			section("git ignore", []string{"add " + ig.Pattern + " to " + ig.Path})
		}
	}
	section("remove (old links into this repo)", p.Remove)
	section("left alone (not ours)", p.Kept)
	section("skipped skills", p.Skipped)
	section("CLASH (nothing will be changed)", p.Clashes)
	fmt.Fprintf(&b, "summary: %d to render, %d to link, %d to copy, %d unchanged, %d to remove, %d clashes\n",
		len(render), len(create), len(copied), len(ok)+len(same), len(p.Remove), len(p.Clashes))
	return b.String()
}

type state int

const (
	absent state = iota
	owned        // a symlink into the repo, or a dir holding only such links
	foreign
)

func intoRepo(repo string, more ...string) func(string) bool {
	var roots []string
	for _, r := range append([]string{repo}, more...) {
		roots = append(roots, r)
		if real, err := filepath.EvalSymlinks(r); err == nil && real != r {
			roots = append(roots, real)
		}
	}
	return func(target string) bool {
		for _, r := range roots {
			if target == r || strings.HasPrefix(target, r+string(filepath.Separator)) {
				return true
			}
		}
		return false
	}
}

func linkTarget(path string) (string, bool) {
	t, err := os.Readlink(path)
	if err != nil {
		return "", false
	}
	if !filepath.IsAbs(t) {
		t = filepath.Join(filepath.Dir(path), t)
	}
	return filepath.Clean(t), true
}

func linksTo(path, target string) bool {
	t, ok := linkTarget(path)
	return ok && t == filepath.Clean(target)
}

func ownedState(path string, in func(string) bool) state {
	fi, err := os.Lstat(path)
	if err != nil {
		return absent
	}
	if fi.Mode()&fs.ModeSymlink != 0 {
		if t, _ := linkTarget(path); in(t) {
			return owned
		}
		return foreign
	}
	return foreign
}

func hash(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// Uninstall lists the manifest links that still point into root or roots and
// the copies that are unchanged since install. With apply set it removes them
// and the manifest; anything else is left alone.
func Uninstall(home string, apply bool, root string, roots ...string) (removed, kept []string, err error) {
	root, err = filepath.Abs(root)
	if err != nil {
		return nil, nil, err
	}
	manifest := filepath.Join(home, ManifestPath)
	m, err := readManifest(manifest)
	if err != nil {
		return nil, nil, err
	}
	in := intoRepo(root, roots...)
	for _, path := range sortedKeys(m.Copies) {
		fi, err := os.Lstat(path)
		if err != nil {
			continue
		}
		cur, err := os.ReadFile(path)
		if !fi.Mode().IsRegular() || err != nil || hash(cur) != m.Copies[path] {
			kept = append(kept, path)
			continue
		}
		removed = append(removed, path)
	}
	for _, path := range m.Links {
		fi, err := os.Lstat(path)
		if err != nil {
			continue
		}
		if t, _ := linkTarget(path); fi.Mode()&fs.ModeSymlink == 0 || !in(t) {
			kept = append(kept, path)
			continue
		}
		removed = append(removed, path)
	}
	if !apply {
		return removed, kept, nil
	}
	for _, path := range removed {
		if err := os.Remove(path); err != nil {
			return nil, nil, err
		}
	}
	if err := os.Remove(manifest); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, nil, err
	}
	if err := os.RemoveAll(filepath.Join(home, CopyStateDir)); err != nil {
		return nil, nil, err
	}
	return removed, kept, nil
}

// manifest records what the installer made: links, and a hash of each copy.
type manifest struct {
	Links  []string
	Copies map[string]string `json:",omitempty"`
}

func readManifest(path string) (manifest, error) {
	var m manifest
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return m, nil
	}
	if err != nil {
		return m, err
	}
	if err := json.Unmarshal(data, &m); err != nil {
		return m, fmt.Errorf("manifest %s: %w", path, err)
	}
	return m, nil
}

func writeManifest(path string, m manifest) error {
	sort.Strings(m.Links)
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return writeAtomic(path, append(data, '\n'))
}

// writeAtomic writes data to a temp file beside path and renames it over path.
func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func writeAtomic(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
