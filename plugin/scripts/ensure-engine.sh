#!/bin/sh
# Builds plugin/bin/estimatorctl from plugin/engine/ if it's missing or
# stale, and reuses it otherwise. See
# specs/002-feature-extraction-plugin/contracts/ensure-engine.md.
set -eu

PLUGIN_ROOT=$(cd "$(dirname "$0")/.." && pwd)
ENGINE_DIR="$PLUGIN_ROOT/engine"
BIN="$PLUGIN_ROOT/bin/estimatorctl"
BUILT_CHECKSUM="$PLUGIN_ROOT/bin/.built-checksum"
SYNC_CHECKSUM="$ENGINE_DIR/.sync-checksum"

if ! command -v go >/dev/null 2>&1; then
  echo "ENGINE_ERROR: go toolchain not found on PATH" >&2
  exit 1
fi

mkdir -p "$PLUGIN_ROOT/bin"

if [ -x "$BIN" ] && [ -f "$BUILT_CHECKSUM" ] && [ -f "$SYNC_CHECKSUM" ]; then
  if [ "$(cat "$BUILT_CHECKSUM")" = "$(cat "$SYNC_CHECKSUM")" ]; then
    exit 0
  fi
fi

if ! (cd "$ENGINE_DIR" && go build -o "$BIN" ./cmd/estimatorctl) 1>&2; then
  rm -f "$BIN" "$BUILT_CHECKSUM"
  echo "ENGINE_ERROR: go build failed for plugin/engine/cmd/estimatorctl" >&2
  exit 1
fi

cp "$SYNC_CHECKSUM" "$BUILT_CHECKSUM"
exit 0
