#!/usr/bin/env bash
# release.sh: tag the next release with svu and push the tag. The tag starts
# .github/workflows/release.yml, which runs GoReleaser.
#
#   scripts/release.sh [-n] [-y] [next|patch|minor|major]
#
#   next   (default) from commit titles: feat: = minor, fix: or anything else = patch
#   -n     dry run: print the tag, change nothing
#   -y     don't ask before pushing
set -euo pipefail

SVU_VERSION=v3.4.1

dry_run=0
yes=0
while getopts "ny" opt; do
  case "$opt" in
    n) dry_run=1 ;;
    y) yes=1 ;;
    *) sed -n '2,10p' "$0"; exit 2 ;;
  esac
done
shift $((OPTIND - 1))
bump=${1:-next}
case "$bump" in
  next | patch | minor | major) ;;
  *) echo "error: unknown bump '$bump' (next, patch, minor or major)" >&2; exit 2 ;;
esac

svu() {
  if command -v svu >/dev/null 2>&1; then
    command svu "$@"
  else
    go run "github.com/caarlos0/svu/v3@${SVU_VERSION}" "$@"
  fi
}

[ "$(git branch --show-current)" = main ] || { echo "error: not on main" >&2; exit 1; }
[ -z "$(git status --porcelain)" ] || { echo "error: working tree is not clean" >&2; exit 1; }
git fetch --quiet --tags origin main
[ "$(git rev-parse HEAD)" = "$(git rev-parse origin/main)" ] || { echo "error: main is not in sync with origin/main" >&2; exit 1; }

current=$(svu current)
if [ "$bump" = next ]; then
  # --always: commits without a conventional prefix still bump the patch.
  # --v0: stay on v0 until a major bump is asked for.
  tag=$(svu next --always --v0)
else
  tag=$(svu "$bump")
fi

if git rev-parse --quiet --verify "refs/tags/$tag" >/dev/null; then
  echo "error: tag $tag already exists" >&2
  exit 1
fi

echo "current: $current"
echo "next:    $tag ($bump) at $(git rev-parse --short HEAD)"
[ "$dry_run" -eq 0 ] || exit 0

if [ "$yes" -eq 0 ]; then
  read -r -p "Tag and push $tag? [y/N] " answer
  [ "$answer" = y ] || [ "$answer" = Y ] || { echo "aborted"; exit 1; }
fi

git tag -a "$tag" -m "$tag"
git push origin "$tag"
echo "pushed $tag: https://github.com/zenkoding/skillset/actions/workflows/release.yml"
