package install

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSharedReferencesRendered(t *testing.T) {
	tests := []struct {
		name    string
		link    string // link text in alpha's SKILL.md
		want    string // rendered path relative to alpha's dir; "" for none
		wantErr string
	}{
		{name: "references dir", link: "[sizing](references/sizing.md)", want: "references/sizing.md"},
		{name: "same dir", link: "[sizing](sizing.md)", want: "sizing.md"},
		{name: "not shared", link: "[other](references/other.md)"},
		{name: "outside the skill", link: "[sizing](../sizing.md)", wantErr: "outside the skill"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, home := t.TempDir(), t.TempDir()
			write(t, filepath.Join(repo, "skills/eng/alpha/SKILL.md"), skill("alpha")+tt.link+"\n")
			write(t, filepath.Join(repo, "shared/sizing.md"), "# Sizing\n")
			src := Source{Root: repo, SkillRoots: []string{"skills"}, Shared: []string{"shared/sizing.md"}}

			p, err := MakeSourcePlan(src, home)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("err = %v, want %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if tt.want == "" {
				if len(p.Renders) != 0 {
					t.Fatalf("renders = %v, want none", p.Renders)
				}
				return
			}
			if err := p.Apply(); err != nil {
				t.Fatal(err)
			}
			got, err := os.ReadFile(filepath.Join(repo, "skills/eng/alpha", tt.want))
			if err != nil || string(got) != "# Sizing\n" {
				t.Fatalf("rendered copy = %q, %v", got, err)
			}
			again, err := MakeSourcePlan(src, home)
			if err != nil {
				t.Fatal(err)
			}
			for _, r := range again.Renders {
				if !r.Exists {
					t.Errorf("%s not up to date after install", r.Path)
				}
			}
		})
	}
}

// A shared file that links another shared file brings it along, next to
// itself, so the link works in every skill that carries the first one.
func TestSharedReferencesFollowLinks(t *testing.T) {
	repo, home := t.TempDir(), t.TempDir()
	write(t, filepath.Join(repo, "skills/eng/alpha/SKILL.md"), skill("alpha")+"[loop](references/loop.md)\n")
	write(t, filepath.Join(repo, "shared/loop.md"), "# Loop\n[sizing](sizing.md)\n")
	write(t, filepath.Join(repo, "shared/sizing.md"), "# Sizing\n[loop](loop.md)\n")
	src := Source{Root: repo, SkillRoots: []string{"skills"}, Shared: []string{"shared/loop.md", "shared/sizing.md"}}

	p, err := MakeSourcePlan(src, home)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, r := range p.Renders {
		rel, _ := filepath.Rel(filepath.Join(repo, "skills/eng/alpha"), r.Path)
		got = append(got, rel)
	}
	want := []string{"references/loop.md", "references/sizing.md"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("renders = %v, want %v", got, want)
	}
}
