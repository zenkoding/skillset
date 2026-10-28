package main

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/zenkoding/skillset/internal/catalog"
)

func newInitCmd(opts *options, in io.Reader) *cobra.Command {
	var hosts []string
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Choose hosts and sets and save them to the config",
		Long: `Choose hosts and sets in a picker and save them to
~/.config/skillset/config.toml. Sets the chosen ones require are added.
With --host or --set, no picker is shown: a flag replaces its part of the
saved config and the other part is kept (everything when nothing is saved).`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, _, err := opts.load(cmd, "", false)
			if err != nil {
				return err
			}
			cfg, _, err := catalog.ReadConfig(opts.home)
			if err != nil {
				return err
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
				if cfg, err = pick(c, cfg, in, cmd.OutOrStdout()); err != nil {
					return err
				}
			}
			for _, h := range cfg.Hosts {
				if !contains(catalog.Hosts, h) {
					return fmt.Errorf("unknown host %q", h)
				}
			}
			var added []string
			if cfg.Sets, added, err = withRequires(c, cfg.Sets); err != nil {
				return err
			}
			if len(added) > 0 {
				cmd.Printf("added required sets: %s\n", strings.Join(added, ", "))
			}
			if err := catalog.WriteConfig(opts.home, cfg); err != nil {
				return err
			}
			cmd.Printf("wrote %s\nnext: skillset plan\n", catalog.ConfigPath(opts.home))
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&hosts, "host", nil, "hosts to install for: "+strings.Join(catalog.Hosts, ", "))
	return cmd
}

// withRequires returns names plus every set they require, in dependency
// order, and the sets it added.
func withRequires(c *catalog.Catalog, names []string) (all, added []string, err error) {
	sets, err := c.Resolve(names)
	if err != nil {
		return nil, nil, err
	}
	for _, s := range sets {
		all = append(all, s.Name)
		if !contains(names, s.Name) {
			added = append(added, s.Name)
		}
	}
	return all, added, nil
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
