#!/bin/sh
# Fails if plugin/engine/ has drifted from the canonical cmd/, internal/,
# go.mod, go.sum. Run `make sync-plugin-engine` to fix. Wired into `make
# test` so drift is a build-time failure, not a trusted assumption
# (specs/002-feature-extraction-plugin/plan.md Constraints).
set -eu

REPO_ROOT=$(cd "$(dirname "$0")/.." && pwd)
TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT

cp -R "$REPO_ROOT/cmd" "$TMP/cmd"
cp -R "$REPO_ROOT/internal" "$TMP/internal"
cp "$REPO_ROOT/go.mod" "$TMP/go.mod"
cp "$REPO_ROOT/go.sum" "$TMP/go.sum"

if ! diff -rq -x '.sync-checksum' "$TMP" "$REPO_ROOT/plugin/engine" >/dev/null 2>&1; then
  echo "plugin/engine/ is out of sync with cmd/, internal/, go.mod, go.sum." >&2
  echo "Run: make sync-plugin-engine" >&2
  diff -rq -x '.sync-checksum' "$TMP" "$REPO_ROOT/plugin/engine" >&2 || true
  exit 1
fi

echo "plugin/engine/ is in sync."
