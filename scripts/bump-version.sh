#!/bin/bash
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

usage() {
  echo "Usage: $0 [patch|minor|major]"
  echo "  patch  - 0.0.X (default)"
  echo "  minor  - 0.X.0"
  echo "  major  - X.0.0"
  exit 1
}

LEVEL="${1:-patch}"

if [[ ! "$LEVEL" =~ ^(patch|minor|major)$ ]]; then
  usage
fi

cd "$REPO_ROOT" || exit 1

CURRENT=$(cat VERSION)
IFS='.' read -r MAJOR MINOR PATCH <<< "$CURRENT"

case "$LEVEL" in
  major) MAJOR=$((MAJOR + 1)); MINOR=0; PATCH=0 ;;
  minor) MINOR=$((MINOR + 1)); PATCH=0 ;;
  patch) PATCH=$((PATCH + 1)) ;;
esac

NEW_VERSION="$MAJOR.$MINOR.$PATCH"

echo "$NEW_VERSION" > VERSION

git add VERSION
git commit -m "chore(release): v${NEW_VERSION}"
git tag "v${NEW_VERSION}"
git push origin HEAD --follow-tags

echo "Bumped to v${NEW_VERSION}, committed, and pushed."
