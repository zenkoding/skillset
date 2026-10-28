package source

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Lock pins the catalog ref and the commit it resolved to.
type Lock struct {
	Ref    string `toml:"ref"`
	Commit string `toml:"commit"`
}

// ReadLock reads the lockfile at path. A missing file is a zero Lock.
func ReadLock(path string) (Lock, error) {
	var l Lock
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return Lock{}, nil
	}
	if err != nil {
		return Lock{}, fmt.Errorf("read lock: %w", err)
	}
	if _, err := toml.Decode(string(data), &l); err != nil {
		return Lock{}, fmt.Errorf("parse lock %s: %w", path, err)
	}
	return l, nil
}

// WriteLock writes l to path, creating parent dirs. The write goes through a
// temp file and a rename, so readers never see a partial lockfile.
func WriteLock(path string, l Lock) error {
	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(l); err != nil {
		return fmt.Errorf("encode lock: %w", err)
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create lock dir: %w", err)
	}
	tmp, err := os.CreateTemp(dir, ".lock-*.toml")
	if err != nil {
		return fmt.Errorf("create temp lock: %w", err)
	}
	defer func() { _ = os.Remove(tmp.Name()) }()
	if _, err := tmp.Write(buf.Bytes()); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("write lock: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("write lock: %w", err)
	}
	if err := os.Rename(tmp.Name(), path); err != nil {
		return fmt.Errorf("write lock: %w", err)
	}
	return nil
}
