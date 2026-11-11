package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/zenkoding/skillset/internal/catalog"
	"github.com/zenkoding/skillset/internal/source"
)

// DefaultSource is the source init records when none is given or saved:
// this repository's main branch.
var DefaultSource = catalog.Source{Name: "skillset", URL: "https://github.com/zenkoding/skillset", Ref: "main"}

// DcgInstall is the command that installs dcg. skillset prints it and
// never runs it.
const DcgInstall = `curl -fsSL "https://raw.githubusercontent.com/Dicklesworthstone/destructive_command_guard/main/install.sh" | bash -s -- --no-configure --easy-mode`

func newInitCmd(opts *options, in io.Reader) *cobra.Command {
	var hosts []string
	cmd := &cobra.Command{
		Use:   "init [source]",
		Short: "Choose a source, hosts and sets and save them to the config",
		Long: `Choose a source, hosts and sets, and save them to
~/.config/skillset/config.toml.

The source is a local directory, a git URL or owner/repo on GitHub, with an
optional @ref: a branch, a tag or a full commit hash. Without @ref it is
main. Without a source, init keeps the saved one, or uses
zenkoding/skillset@main. A git source is fetched to list its sets.

Hosts and sets are chosen in a picker. With --host or --set, no picker is
shown: a flag replaces its part of the saved config and the other part is
kept (everything when nothing is saved). A new source drops the saved sets
it lacks. The config lists the chosen sets;
the sets they require are added at install time.`,
		Example: `  skillset init
  skillset init zenkoding/skillset@v1.2.0
  skillset init https://github.com/zenkoding/skillset@feature/x --host claude --set base
  skillset init ~/src/skillset`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Flags().Changed("catalog") || cmd.Flags().Changed("url") {
				return errors.New("init takes its source as an argument: skillset init <source>")
			}
			if opts.home == "" {
				return errors.New("no home directory: set HOME or --home")
			}
			cfg, _, err := catalog.ReadConfig(opts.home)
			if err != nil {
				return err
			}
			var saved catalog.Source
			if len(cfg.Sources) == 1 {
				saved = cfg.Sources[0]
			}
			src, assumed := DefaultSource, false
			switch {
			case len(args) == 1:
				if src, assumed, err = parseSource(args[0]); err != nil {
					return err
				}
			case saved != catalog.Source{}:
				src = saved
			}
			c, err := opts.openSource(cmd, src)
			if assumed && errors.Is(err, source.ErrRefNotFound) {
				return fmt.Errorf("%w (no @ref given, so main was assumed: add @<branch>)", err)
			}
			if err != nil {
				return err
			}
			if src != saved {
				// The saved sets came from another source: keep the ones
				// this source has, and say which were dropped.
				var kept, dropped []string
				for _, n := range cfg.Sets {
					if _, ok := c.Sets[n]; ok {
						kept = append(kept, n)
					} else {
						dropped = append(dropped, n)
					}
				}
				if len(dropped) > 0 {
					cmd.Printf("dropped saved sets this source lacks: %s\n", strings.Join(dropped, ", "))
				}
				// With none left, an empty list would read as "nothing
				// saved" and pick every set: make the user choose.
				if len(dropped) > 0 && len(kept) == 0 && len(hosts) > 0 && len(opts.sets) == 0 {
					return errors.New("none of the saved sets are in this source: pass --set")
				}
				cfg.Sets = kept
			}
			switch {
			case len(hosts) > 0 || len(opts.sets) > 0:
				// A flag replaces its part; the other keeps the saved
				// value, or everything when nothing is saved.
				if len(hosts) > 0 {
					cfg.Hosts = hosts
				}
				if len(opts.sets) > 0 {
					cfg.Sets = opts.sets
				}
				if len(cfg.Hosts) == 0 {
					cfg.Hosts = catalog.Hosts
				}
				if len(cfg.Sets) == 0 {
					cfg.Sets = c.Names()
				}
			case !isTerminal(in):
				return errors.New("no terminal for the picker: pass --host or --set")
			default:
				picked, err := pick(c, cfg, in, cmd.OutOrStdout())
				if err != nil {
					return err
				}
				cfg.Hosts, cfg.Sets = picked.Hosts, picked.Sets
			}
			for _, h := range cfg.Hosts {
				if !contains(catalog.Hosts, h) {
					return fmt.Errorf("unknown host %q", h)
				}
			}
			// Check the names; the config keeps only the chosen ones.
			if _, err := c.Resolve(cfg.Sets); err != nil {
				return err
			}
			cfg.Sources = []catalog.Source{src}
			if err := catalog.WriteConfig(opts.home, cfg); err != nil {
				return err
			}
			cmd.Printf("wrote %s\n", catalog.ConfigPath(opts.home))
			if _, err := exec.LookPath("dcg"); err != nil {
				cmd.Printf("warning: dcg is not installed, so destructive commands are not guarded. Install it with:\n  %s\n", DcgInstall)
			}
			// Until install (#50), plan and apply ignore the source.
			if src.Path != "" {
				cmd.Printf("next: skillset --catalog %s plan\n", filepath.Join(src.Path, CatalogDir))
			} else {
				cmd.Println("note: plan and apply do not read this source yet; they use the pinned release")
			}
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&hosts, "host", nil, "hosts to install for: "+strings.Join(catalog.Hosts, ", "))
	return cmd
}

// parseSource turns init's argument into a source: an existing directory
// is a local source, anything else a git source at its ref or main.
// assumed reports that the argument had no ref, so main was used.
func parseSource(arg string) (src catalog.Source, assumed bool, err error) {
	if fi, err := os.Stat(arg); err == nil && fi.IsDir() {
		abs, err := filepath.Abs(arg)
		if err != nil {
			return catalog.Source{}, false, err
		}
		return catalog.Source{Name: filepath.Base(abs), Path: abs}, false, nil
	}
	url, ref, err := source.Parse(arg)
	if err != nil {
		return catalog.Source{}, false, err
	}
	if ref == "" {
		ref, assumed = DefaultSource.Ref, true
	}
	return catalog.Source{Name: source.Name(url), URL: url, Ref: ref}, assumed, nil
}

// openSource loads the catalog in src. A git source is resolved and
// fetched into the cache by commit, so a branch name with a slash works.
func (o *options) openSource(cmd *cobra.Command, src catalog.Source) (*catalog.Catalog, error) {
	if src.Path != "" {
		cmd.Printf("source: %s\n", src.Path)
		return catalog.Load(filepath.Join(src.Path, CatalogDir))
	}
	ctx := cmd.Context()
	commit, err := source.Resolve(ctx, src.URL, src.Ref)
	if err != nil {
		return nil, err
	}
	dir, err := source.FetchCommit(ctx, src.URL, commit, filepath.Join(o.home, CacheDir))
	if err != nil {
		return nil, err
	}
	cmd.Printf("source: %s at %s (%s)\n", src.URL, src.Ref, short(commit))
	c, err := catalog.Load(filepath.Join(dir, CatalogDir))
	if err != nil {
		return nil, fmt.Errorf("source %s at %s: %w", src.URL, src.Ref, err)
	}
	return c, nil
}

// pick shows the picker, starting from cur (or everything when cur is
// empty). Sets show what they require; init adds those after the picker.
func pick(c *catalog.Catalog, cur catalog.Config, in io.Reader, out io.Writer) (catalog.Config, error) {
	if len(cur.Hosts) == 0 {
		cur.Hosts = catalog.Hosts
	}
	if len(cur.Sets) == 0 {
		cur.Sets = c.Names()
	}
	// huh ticks the options already in the bound values.
	cfg := catalog.Config{Hosts: append([]string(nil), cur.Hosts...), Sets: append([]string(nil), cur.Sets...)}
	var hostOpts []huh.Option[string]
	for _, h := range catalog.Hosts {
		hostOpts = append(hostOpts, huh.NewOption(h, h).Selected(contains(cur.Hosts, h)))
	}
	var setOpts []huh.Option[string]
	for _, n := range c.Names() {
		s := c.Sets[n]
		label := n
		if s.Description != "" {
			label += ": " + s.Description
		}
		if len(s.Requires) > 0 {
			label += " (requires " + strings.Join(s.Requires, ", ") + ")"
		}
		setOpts = append(setOpts, huh.NewOption(label, n).Selected(contains(cur.Sets, n)))
	}
	nonEmpty := func(what string) func([]string) error {
		return func(v []string) error {
			if len(v) == 0 {
				return fmt.Errorf("pick at least one %s", what)
			}
			return nil
		}
	}
	form := huh.NewForm(
		huh.NewGroup(huh.NewMultiSelect[string]().
			Title("Hosts").Description("Which agents do you run?").
			Options(hostOpts...).Validate(nonEmpty("host")).Value(&cfg.Hosts)),
		huh.NewGroup(huh.NewMultiSelect[string]().
			Title("Sets").Description("Required sets are added for you.").
			Options(setOpts...).Validate(nonEmpty("set")).Value(&cfg.Sets)),
	).WithInput(in).WithOutput(out)
	if err := form.Run(); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func contains(list []string, v string) bool {
	for _, s := range list {
		if s == v {
			return true
		}
	}
	return false
}
