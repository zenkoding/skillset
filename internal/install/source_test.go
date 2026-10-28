package install

import "path/filepath"

// repoSource is a catalog-like layout for tests: instructions/, skills/,
// claude/ and dcg/ under repo.
func repoSource(repo string) Source {
	src := Source{
		Root:         repo,
		Instructions: map[string][]string{},
		SkillRoots:   []string{"skills"},
		Configs: []Config{
			{Home: ".claude/statusline.sh", Repo: "claude/statusline.sh"},
			{Home: ".claude/output-styles", Repo: "claude/output-styles"},
			{Home: ".config/dcg/config.toml", Repo: "dcg/config.toml"},
			{Home: ".config/dcg/packs", Repo: "dcg/packs"},
		},
		Copies: []Config{{Home: ".claude/settings.json", Repo: "claude/settings.json"}},
	}
	for _, ins := range Instructions {
		src.Instructions[ins.Host] = []string{"instructions/shared.md", filepath.Join("instructions", ins.Host+".md")}
	}
	return src
}

func makePlan(repo, home string) (*Plan, error) {
	return MakeSourcePlan(repoSource(repo), home)
}
