#!/bin/sh
# Tests plugin/scripts/on-tasks-md-change.sh per
# specs/002-feature-extraction-plugin/contracts/hook.md.
set -eu

PLUGIN_ROOT=$(cd "$(dirname "$0")/.." && pwd)
REPO_ROOT=$(cd "$PLUGIN_ROOT/.." && pwd)
HOOK="$PLUGIN_ROOT/scripts/on-tasks-md-change.sh"

fail() { echo "FAIL: $1" >&2; exit 1; }

# 1. A valid tasks.md write (real feature dir) produces additionalContext
# naming that feature directory.
input=$(printf '{"hook_event_name":"PostToolUse","tool_name":"Write","tool_input":{"file_path":"specs/002-feature-extraction-plugin/tasks.md"},"cwd":"%s"}' "$REPO_ROOT")
out=$(printf '%s' "$input" | "$HOOK")
case "$out" in
  *'"additionalContext"'*'specs/002-feature-extraction-plugin'*) ;;
  *) fail "expected additionalContext naming the feature dir, got: $out" ;;
esac

# 2. A non-tasks.md write is a silent no-op.
input=$(printf '{"hook_event_name":"PostToolUse","tool_name":"Write","tool_input":{"file_path":"specs/002-feature-extraction-plugin/spec.md"},"cwd":"%s"}' "$REPO_ROOT")
out=$(printf '%s' "$input" | "$HOOK")
[ -z "$out" ] || fail "expected no output for a non-tasks.md write, got: $out"

# 3. A tasks.md outside specs/ (no siblings) is a silent no-op.
tmp=$(mktemp -d)
mkdir -p "$tmp/scratch"
touch "$tmp/scratch/tasks.md"
input=$(printf '{"hook_event_name":"PostToolUse","tool_name":"Write","tool_input":{"file_path":"scratch/tasks.md"},"cwd":"%s"}' "$tmp")
out=$(printf '%s' "$input" | "$HOOK")
[ -z "$out" ] || fail "expected no output for a stray tasks.md, got: $out"
rm -rf "$tmp"

echo "PASS: test_on_tasks_md_change.sh"
