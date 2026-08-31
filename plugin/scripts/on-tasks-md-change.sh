#!/bin/sh
# PostToolUse hook entrypoint. Reads the tool-call JSON from stdin, checks
# whether it's a tasks.md write for a real Spec-Kit feature, and if so
# injects context telling Claude to run the estimate-feature skill.
# See specs/002-feature-extraction-plugin/contracts/hook.md.
set -eu

PLUGIN_ROOT=$(cd "$(dirname "$0")/.." && pwd)
LOCATE="$PLUGIN_ROOT/scripts/locate-feature.sh"

input=$(cat)

file_path=$(printf '%s' "$input" \
  | grep -o '"file_path"[[:space:]]*:[[:space:]]*"[^"]*"' \
  | head -n1 \
  | sed -E 's/.*:[[:space:]]*"([^"]*)"/\1/')

if [ -z "$file_path" ]; then
  exit 0
fi

cwd=$(printf '%s' "$input" \
  | grep -o '"cwd"[[:space:]]*:[[:space:]]*"[^"]*"' \
  | head -n1 \
  | sed -E 's/.*:[[:space:]]*"([^"]*)"/\1/')

if [ -n "$cwd" ] && [ -d "$cwd" ]; then
  cd "$cwd"
fi

feature_dir=$("$LOCATE" "$file_path" 2>/dev/null) || exit 0

printf '{"hookSpecificOutput": {"hookEventName": "PostToolUse", "additionalContext": "A Spec-Kit feature just completed planning: %s. Run the estimate-feature skill for this feature directory now."}}\n' "$feature_dir"

exit 0
