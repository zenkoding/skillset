package catalog

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"github.com/BurntSushi/toml"
)

// ConfigDir holds the user's config and local set, relative to HOME.
const ConfigDir = ".config/skillset"

// LocalDir is the user's own set, relative to HOME. It is applied last.
const LocalDir = ConfigDir + "/local"

// LocalName is the local set's name.
const LocalName = "local"

// Schema is the config format WriteConfig writes. A config without a
// schema predates it and holds only hosts and sets.
const Schema = 1

// Config is what the user asked for: the source, hosts and sets. Sets
// lists only the sets the user chose, not the sets those require.
type Config struct {
	Schema  int      `toml:"schema,omitempty"`
	Hosts   []string `toml:"hosts"`
	Sets    []string `toml:"sets"`
	Sources []Source `toml:"sources,omitempty"`
}

// Source is where sets come from: a git URL at a ref (a branch, tag or
// commit), or a local directory, which has no ref.
type Source struct {
	Name string `toml:"name"`
	URL  string `toml:"url,omitempty"`
	Ref  string `toml:"ref,omitempty"`
	Path string `toml:"path,omitempty"`
}

// header heads every config WriteConfig writes.
const header = "# Managed by skillset. Change it with `skillset install`.\n"

func (c Config) validate() error {
	if c.Schema > Schema {
		return fmt.Errorf("schema %d is newer than this skillset reads (%d): upgrade skillset", c.Schema, Schema)
	}
	if len(c.Sources) > 1 {
		return fmt.Errorf("%d sources: a config holds one source for now", len(c.Sources))
	}
	for _, s := range c.Sources {
		if s.Name == "" {
			return errors.New("a source has no name")
		}
		if (s.URL == "") == (s.Path == "") {
			return fmt.Errorf("source %q needs a url or a path, not both", s.Name)
		}
		if s.URL != "" && s.Ref == "" {
			return fmt.Errorf("source %q has a url but no ref", s.Name)
		}
		if s.Path != "" && s.Ref != "" {
			return fmt.Errorf("source %q is a local path and takes no ref", s.Name)
		}
	}
	return nil
}

// ConfigPath is the config file under home.
func ConfigPath(home string) string {
	return filepath.Join(home, ConfigDir, "config.toml")
}

// ReadConfig reads home's config. ok is false when there is none.
func ReadConfig(home string) (cfg Config, ok bool, err error) {
	path := ConfigPath(home)
	md, err := toml.DecodeFile(path, &cfg)
	if errors.Is(err, fs.ErrNotExist) {
		return Config{}, false, nil
	}
	if err != nil {
		return Config{}, false, err
	}
	if keys := md.Undecoded(); len(keys) > 0 {
		return Config{}, false, fmt.Errorf("%s: unknown keys %v", path, keys)
	}
	if err := cfg.validate(); err != nil {
		return Config{}, false, fmt.Errorf("%s: %w", path, err)
	}
	return cfg, true, nil
}

// WriteConfig writes cfg to home's config file in the current schema:
// keys in a fixed order, lists sorted, no timestamps, so the same choices
// always give the same bytes and two people's changes merge cleanly.
func WriteConfig(home string, cfg Config) error {
	cfg.Schema = Schema
	cfg.Hosts = sorted(cfg.Hosts)
	cfg.Sets = sorted(cfg.Sets)
	cfg.Sources = append([]Source(nil), cfg.Sources...)
	sort.Slice(cfg.Sources, func(i, j int) bool { return cfg.Sources[i].Name < cfg.Sources[j].Name })
	if err := cfg.validate(); err != nil {
		return err
	}
	var b bytes.Buffer
	b.WriteString(header)
	enc := toml.NewEncoder(&b)
	enc.Indent = ""
	if err := enc.Encode(cfg); err != nil {
		return err
	}
	path := ConfigPath(home)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, b.Bytes(), 0o644)
}

func sorted(s []string) []string {
	out := append([]string(nil), s...)
	sort.Strings(out)
	return out
}

// LoadLocal reads home's local set, or returns nil when its directory is
// missing. Without a set.toml it is a set named local for every host.
// Its requires must name catalog sets.
func (c *Catalog) LoadLocal(home string) (*Set, error) {
	dir := filepath.Join(home, LocalDir)
	if fi, err := os.Stat(dir); errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else if !fi.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", dir)
	}
	var s *Set
	if _, err := os.Stat(filepath.Join(dir, Manifest)); errors.Is(err, fs.ErrNotExist) {
		s = &Set{Name: LocalName, Hosts: Hosts, Dir: dir}
		if err := s.resolve(dir); err != nil {
			return nil, err
		}
	} else if s, err = LoadSet(dir, dir); err != nil {
		return nil, err
	}
	if s.Name != LocalName {
		return nil, fmt.Errorf("%s: name must be %q", filepath.Join(dir, Manifest), LocalName)
	}
	for _, r := range s.Requires {
		if c.Sets[r] == nil {
			return nil, fmt.Errorf("local set requires unknown set %q", r)
		}
	}
	return s, nil
}

// Selection is what to install: hosts and sets in order, local set last.
type Selection struct {
	Hosts []string
	Sets  []*Set
}

// Select resolves names (or every set when names is empty) plus what they
// and the local set require, then appends the local set. Empty hosts means
// every host.
func (c *Catalog) Select(names, hosts []string, local *Set) (Selection, error) {
	if len(names) == 0 {
		names = c.Names()
	}
	if len(hosts) == 0 {
		hosts = Hosts
	}
	if local != nil {
		names = append(append([]string(nil), names...), local.Requires...)
	}
	sets, err := c.Resolve(names)
	if err != nil {
		return Selection{}, err
	}
	if local != nil {
		sets = append(sets, local)
	}
	return Selection{Hosts: hosts, Sets: sets}, nil
}
