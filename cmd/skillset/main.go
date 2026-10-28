// Command skillset installs sets from a catalog into ~/.claude, ~/.codex
// and ~/.agents.
//
//	skillset init      choose hosts and sets
//	skillset plan      show what would change
//	skillset apply     make the changes
//	skillset update    move the pin to the newest release tag, then apply
//	skillset uninstall remove what apply installed
//
// The catalog is this repository's config/ directory, fetched at the ref
// pinned in ~/.config/skillset/lock.toml into ~/.cache/skillset/<ref>/. With
// --catalog DIR, links point into DIR instead, so edits there are live.
// Choices live in ~/.config/skillset/config.toml; without one, every set
// and host is installed. ~/.config/skillset/local/ is a set applied last.
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"

	"github.com/zenkoding/skillset/internal/catalog"
	"github.com/zenkoding/skillset/internal/install"
	"github.com/zenkoding/skillset/internal/source"
)

// DefaultURL is the catalog repository.
const DefaultURL = "https://github.com/zenkoding/skillset.git"

// CatalogDir is the catalog inside the repository.
const CatalogDir = "config"

// CacheDir holds fetched catalogs, one directory per ref, relative to HOME.
const CacheDir = ".cache/skillset"

// LockPath is the pinned ref, relative to HOME.
const LockPath = catalog.ConfigDir + "/lock.toml"

func main() {
	if err := newRootCmd(os.Stdin, os.Stdout).Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

type options struct {
	catalog string
	url     string
	home    string
	sets    []string
}

func newRootCmd(in io.Reader, out io.Writer) *cobra.Command {
	opts := &options{}
	root := &cobra.Command{
		Use:           "skillset",
		Short:         "Install agent instructions, skills and settings from composable sets",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.SetOut(out)
	root.PersistentFlags().StringVar(&opts.catalog, "catalog", "", "catalog directory to use instead of a release (links point into it)")
	root.PersistentFlags().StringVar(&opts.url, "url", DefaultURL, "catalog git repository")
	root.PersistentFlags().StringVar(&opts.home, "home", os.Getenv("HOME"), "home directory to install into")
	root.PersistentFlags().StringSliceVar(&opts.sets, "set", nil, "sets to install, with what they require (default: config, else all)")

	// apply installs from ref ("" means the pinned ref), fetching it again
	// first when fresh is set.
	apply := func(cmd *cobra.Command, ref string, fresh bool) error {
		p, lock, err := opts.plan(cmd, ref, fresh)
		if err != nil {
			return err
		}
		cmd.Print(p)
		if isTerminal(in) {
			if err := p.Ask(in, out); err != nil {
				return err
			}
		}
		if err := p.Apply(); err != nil {
			return err
		}
		if lock != nil {
			if err := source.WriteLock(filepath.Join(opts.home, LockPath), *lock); err != nil {
				return err
			}
			cmd.Printf("installed %s (%s)\n", lock.Ref, short(lock.Commit))
			return nil
		}
		cmd.Println("installed")
		return nil
	}

	root.AddCommand(newInitCmd(opts, in))
	root.AddCommand(&cobra.Command{
		Use:   "plan",
		Short: "Show what apply would change",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			p, _, err := opts.plan(cmd, "", false)
			if err != nil {
				return err
			}
			cmd.Print(p)
			return p.Err()
		},
	})
	root.AddCommand(&cobra.Command{
		Use:   "apply",
		Short: "Install the chosen sets at the pinned ref",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return apply(cmd, "", false)
		},
	})
	var ref string
	update := &cobra.Command{
		Use:   "update",
		Short: "Pin the newest release tag (or --ref), then apply",
		Long: `Pin the newest release tag, or --ref, then apply. Older refs stay
cached: update --ref <old tag> switches back. A branch ref such as main is
fetched fresh every time.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if opts.catalog != "" {
				return errors.New("update pins a release; with --catalog, pull the checkout and apply")
			}
			to := ref
			if to == "" {
				tag, err := source.LatestTag(cmd.Context(), opts.url)
				if err != nil {
					return err
				}
				to = tag
			}
			// Fetch reuses a cached ref; a branch must be fetched again.
			return apply(cmd, to, !release.MatchString(to))
		},
	}
	update.Flags().StringVar(&ref, "ref", "", "tag or branch to pin instead of the newest release tag")
	root.AddCommand(update)
	var dryRun bool
	uninstall := &cobra.Command{
		Use:   "uninstall",
		Short: "Remove the links and unedited copies apply installed",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if opts.home == "" {
				return errors.New("no home directory: set HOME or --home")
			}
			// Without --catalog, links point into the release cache. The
			// local set's links are ours too.
			root := filepath.Join(opts.home, CacheDir)
			if opts.catalog != "" {
				root = opts.catalog
			}
			removed, kept, err := install.Uninstall(opts.home, !dryRun, root, filepath.Join(opts.home, catalog.LocalDir))
			if err != nil {
				return err
			}
			verb := "removed"
			if dryRun {
				verb = "would remove"
			}
			for _, p := range removed {
				cmd.Printf("%s %s\n", verb, p)
			}
			for _, p := range kept {
				cmd.Printf("kept %s (changed or not ours)\n", p)
			}
			return nil
		},
	}
	uninstall.Flags().BoolVar(&dryRun, "dry-run", false, "only list what would be removed")
	root.AddCommand(uninstall)
	return root
}

// release matches release tags, which never change once fetched.
var release = regexp.MustCompile(`^v\d+\.\d+\.\d+$`)

// load opens the catalog: --catalog, or the remote one at ref, else the
// pinned ref, else the newest tag; fresh fetches ref again. It prints which
// catalog it uses. lock is what to pin after an apply; it is nil with
// --catalog.
func (o *options) load(cmd *cobra.Command, ref string, fresh bool) (c *catalog.Catalog, lock *source.Lock, err error) {
	if o.home == "" {
		return nil, nil, errors.New("no home directory: set HOME or --home")
	}
	if o.catalog != "" {
		if c, err = catalog.Load(o.catalog); err != nil {
			return nil, nil, err
		}
		cmd.Printf("catalog: %s\n", c.Root)
		return c, nil, nil
	}
	ctx := cmd.Context()
	if ref == "" {
		l, err := source.ReadLock(filepath.Join(o.home, LockPath))
		if err != nil {
			return nil, nil, err
		}
		ref = l.Ref
	}
	if ref == "" {
		if ref, err = source.LatestTag(ctx, o.url); err != nil {
			return nil, nil, err
		}
	}
	dir, commit, err := o.fetch(ctx, ref, fresh)
	if err != nil {
		return nil, nil, err
	}
	cmd.Printf("catalog: %s at %s (%s)\n", o.url, ref, short(commit))
	if c, err = catalog.Load(filepath.Join(dir, CatalogDir)); err != nil {
		return nil, nil, fmt.Errorf("catalog %s: %w", ref, err)
	}
	return c, &source.Lock{Ref: ref, Commit: commit}, nil
}

// fetch fetches ref into the cache. With fresh, the cached checkout moves
// aside first and comes back if the fetch fails, so links into it keep
// working when offline.
func (o *options) fetch(ctx context.Context, ref string, fresh bool) (dir, commit string, err error) {
	// The ref names a cache directory: check it before touching any path.
	if err := source.ValidRef(ref); err != nil {
		return "", "", err
	}
	cache := filepath.Join(o.home, CacheDir)
	cached := filepath.Join(cache, ref)
	if _, err := os.Stat(cached); !fresh || err != nil {
		return source.Fetch(ctx, o.url, ref, cache)
	}
	aside, err := os.MkdirTemp(cache, ".stale-") // cache exists: it holds cached
	if err != nil {
		return "", "", err
	}
	old := filepath.Join(aside, ref)
	if err := os.Rename(cached, old); err != nil {
		return "", "", err
	}
	dir, commit, err = source.Fetch(ctx, o.url, ref, cache)
	if err != nil {
		if rerr := os.Rename(old, cached); rerr != nil {
			return "", "", errors.Join(err, rerr)
		}
		_ = os.Remove(aside)
		return "", "", err
	}
	return dir, commit, os.RemoveAll(aside)
}

// plan installs the sets and hosts from --set, else the config, else
// everything, plus the local set.
func (o *options) plan(cmd *cobra.Command, ref string, fresh bool) (*install.Plan, *source.Lock, error) {
	c, lock, err := o.load(cmd, ref, fresh)
	if err != nil {
		return nil, nil, err
	}
	cfg, _, err := catalog.ReadConfig(o.home)
	if err != nil {
		return nil, nil, err
	}
	names := cfg.Sets
	if len(o.sets) > 0 {
		names = o.sets
	}
	local, err := c.LoadLocal(o.home)
	if err != nil {
		return nil, nil, err
	}
	sel, err := c.Select(names, cfg.Hosts, local)
	if err != nil && len(o.sets) == 0 && len(cfg.Sets) > 0 {
		// The names came from the config, which may predate the catalog.
		return nil, nil, fmt.Errorf("%w (sets from %s; run skillset init)", err, catalog.ConfigPath(o.home))
	}
	if err != nil {
		return nil, nil, err
	}
	src, err := c.Source(sel.Sets, sel.Hosts)
	if err != nil {
		return nil, nil, err
	}
	if lock != nil {
		// Links into any cached ref are ours, so switching refs relinks.
		src.Roots = append(src.Roots, filepath.Join(o.home, CacheDir))
	}
	p, err := install.MakeSourcePlan(src, o.home)
	return p, lock, err
}

func short(commit string) string {
	if len(commit) > 12 {
		return commit[:12]
	}
	return commit
}

func isTerminal(r io.Reader) bool {
	f, ok := r.(*os.File)
	if !ok {
		return false
	}
	fi, err := f.Stat()
	return err == nil && fi.Mode()&os.ModeCharDevice != 0
}
