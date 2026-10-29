---
name: cut-release
description: Cut a new release of a repository that releases from semver tags: pick the next version, tag the default branch and push the tag. Use when the user asks to release, cut, tag or publish a new version.
compatibility: Requires Git, a writable checkout with push access, and the repository's release tooling (for example svu).
---

# Cut a release

A release is a semver tag (`vX.Y.Z`) on the default branch. Pushing the tag is the release: CI builds and publishes from it, and it cannot be quietly undone once others have fetched it. Get the user's go-ahead for the exact tag before pushing.

## Find the repository's way

Read `CONTRIBUTING.md`, `README.md` and the `Makefile` for a release target or script, and `.github/workflows/` for the workflow that tags trigger. Prefer the repository's own tooling over raw commands; it encodes its checks and version rules. In skillset that is:

```bash
make release/plan          # show the next tag, change nothing
make release               # tag main and push (asks first; scripts/release.sh)
make release BUMP=minor    # force patch, minor or major
```

Without repository tooling, use [svu](https://github.com/caarlos0/svu): `svu current`, then `svu next --always --v0` or `svu patch|minor|major`.

## Check before tagging

- On the default branch, clean tree, in sync with the remote after fetching tags.
- CI is green on that commit.
- The tag does not exist yet. The bump matches what changed since the last tag (`git log <last-tag>..HEAD --oneline`): new features mean minor, fixes or internal changes mean patch, breaking changes mean major (minor while on v0). Override svu's guess when commit titles mislead it.

## Tag and follow through

Show the user the current tag, the proposed tag, the commit and the reason for the bump, and wait for approval. Then run the release command non-interactively only if the user approved that exact tag (for skillset: `scripts/release.sh -y <bump>`).

After the push, watch the release workflow run for that tag. Report the tag, the workflow result and the release URL. If the workflow fails, report the failing step; don't delete or move a pushed tag without the user's explicit say-so.
