#!/bin/sh
# Resolves the owning specs/<id>-*/ directory for a tasks.md path, per
# specs/002-feature-extraction-plugin/contracts/hook.md steps 2-3.
#
# Usage: locate-feature.sh <path/to/tasks.md>
# On success: prints the feature directory to stdout, exits 0.
# On any mismatch (not a tasks.md, not under specs/, missing spec.md or
# plan.md sibling): prints nothing, exits 1.
set -eu

path="${1:-}"

if [ "$(basename "$path")" != "tasks.md" ]; then
  exit 1
fi

dir=$(dirname "$path")
parent=$(dirname "$dir")
parent_base=$(basename "$parent")

if [ "$parent_base" != "specs" ]; then
  exit 1
fi

if [ ! -f "$dir/spec.md" ] || [ ! -f "$dir/plan.md" ]; then
  exit 1
fi

echo "$dir"
