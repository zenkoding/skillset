package source

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLock(t *testing.T) {
	tests := []struct {
		name string
		lock Lock
	}{
		{name: "tag", lock: Lock{Ref: "v1.2.3", Commit: "0123456789abcdef0123456789abcdef01234567"}},
		{name: "branch", lock: Lock{Ref: "main", Commit: "fedcba9876543210fedcba9876543210fedcba98"}},
		{name: "zero", lock: Lock{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "nested", "lock.toml")
			if err := WriteLock(path, tt.lock); err != nil {
				t.Fatalf("WriteLock: %v", err)
			}
			got, err := ReadLock(path)
			if err != nil {
				t.Fatalf("ReadLock: %v", err)
			}
			if got != tt.lock {
				t.Errorf("round trip = %+v, want %+v", got, tt.lock)
			}
			entries, _ := os.ReadDir(filepath.Dir(path))
			if len(entries) != 1 {
				t.Errorf("lock dir has %d entries, want only lock.toml", len(entries))
			}
		})
	}
}

func TestReadLock(t *testing.T) {
	tests := []struct {
		name    string
		content *string
		want    Lock
		wantErr bool
	}{
		{name: "missing file", content: nil, want: Lock{}},
		{name: "valid", content: ptr("ref = \"v1.0.0\"\ncommit = \"abc\"\n"), want: Lock{Ref: "v1.0.0", Commit: "abc"}},
		{name: "invalid toml", content: ptr("ref = \n"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "lock.toml")
			if tt.content != nil {
				if err := os.WriteFile(path, []byte(*tt.content), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			got, err := ReadLock(path)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ReadLock() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ReadLock() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func ptr(s string) *string { return &s }
