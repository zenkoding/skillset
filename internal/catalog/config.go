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

// Config is the user's choice of hosts and sets.
type Config struct {
	Hosts []string `toml:"hosts"`
	Sets  []string `toml:"sets"`
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
	return cfg, true, nil
}

// WriteConfig writes cfg to home's config file, sorted.
func WriteConfig(home string, cfg Config) error {
	cfg.Hosts = sorted(cfg.Hosts)
	cfg.Sets = sorted(cfg.Sets)
	var b bytes.Buffer
	b.WriteString("# Written by skillset init.\n")
	if err := toml.NewEncoder(&b).Encode(cfg); err != nil {
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
