#!/bin/sh
# Tests plugin/scripts/locate-feature.sh per
# specs/002-feature-extraction-plugin/contracts/hook.md steps 2-3.
set -eu

PLUGIN_ROOT=$(cd "$(dirname "$0")/.." && pwd)
REPO_ROOT=$(cd "$PLUGIN_ROOT/.." && pwd)
LOCATE="$PLUGIN_ROOT/scripts/locate-feature.sh"

fail() { echo "FAIL: $1" >&2; exit 1; }

# 1. A real, valid feature directory is accepted.
got=$("$LOCATE" "$REPO_ROOT/specs/002-feature-extraction-plugin/tasks.md")
[ "$got" = "$REPO_ROOT/specs/002-feature-extraction-plugin" ] \
  || fail "expected the real feature directory, got: $got"

# 2. A tasks.md outside any specs/ directory is rejected.
tmp=$(mktemp -d)
touch "$tmp/tasks.md"
if "$LOCATE" "$tmp/tasks.md" >/tmp/locate-feature-test-out 2>&1; then
  fail "expected rejection for a tasks.md outside specs/"
fi
[ -s /tmp/locate-feature-test-out ] && fail "expected no stdout on rejection"
rm -rf "$tmp" /tmp/locate-feature-test-out

# 3. A tasks.md under specs/ missing spec.md/plan.md siblings is rejected.
tmp=$(mktemp -d)
mkdir -p "$tmp/specs/999-incomplete"
touch "$tmp/specs/999-incomplete/tasks.md"
if "$LOCATE" "$tmp/specs/999-incomplete/tasks.md" >/dev/null 2>&1; then
  fail "expected rejection when spec.md/plan.md are missing"
fi
rm -rf "$tmp"

# 4. A file that merely ends in "tasks.md" as a substring, not the exact
# basename, is rejected.
tmp=$(mktemp -d)
mkdir -p "$tmp/specs/999-incomplete"
touch "$tmp/specs/999-incomplete/spec.md" "$tmp/specs/999-incomplete/plan.md"
touch "$tmp/specs/999-incomplete/mytasks.md"
if "$LOCATE" "$tmp/specs/999-incomplete/mytasks.md" >/dev/null 2>&1; then
  fail "expected rejection for a non-exact tasks.md filename"
fi
rm -rf "$tmp"

echo "PASS: test_locate_feature.sh"
