#!/bin/sh
# Syncs the canonical Estimation Engine source (cmd/, internal/, go.mod,
# go.sum) into plugin/engine/ so the Claude Code plugin can build a
# standalone estimatorctl binary without this whole repository being
# present (specs/002-feature-extraction-plugin/research.md "Engine sync").
#
# plugin/engine/ is generated output: never hand-edit it, always re-run
# this script after changing cmd/ or internal/.
set -eu

REPO_ROOT=$(cd "$(dirname "$0")/.." && pwd)
DEST="$REPO_ROOT/plugin/engine"

rm -rf "$DEST"
mkdir -p "$DEST"

cp -R "$REPO_ROOT/cmd" "$DEST/cmd"
cp -R "$REPO_ROOT/internal" "$DEST/internal"
cp "$REPO_ROOT/go.mod" "$DEST/go.mod"
cp "$REPO_ROOT/go.sum" "$DEST/go.sum"

find "$DEST" -type f \( -name '*.go' -o -name 'go.mod' -o -name 'go.sum' \) \
  | LC_ALL=C sort \
  | while IFS= read -r f; do cat "$f"; done \
  | cksum > "$DEST/.sync-checksum"

echo "Synced engine source into $DEST ($(cat "$DEST/.sync-checksum"))"
