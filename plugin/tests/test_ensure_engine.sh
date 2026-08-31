#!/bin/sh
# Tests plugin/scripts/ensure-engine.sh: build, reuse, rebuild-on-change,
# and missing-toolchain error paths.
set -eu

PLUGIN_ROOT=$(cd "$(dirname "$0")/.." && pwd)
ENSURE="$PLUGIN_ROOT/scripts/ensure-engine.sh"
BIN="$PLUGIN_ROOT/bin/estimatorctl"
BUILT_CHECKSUM="$PLUGIN_ROOT/bin/.built-checksum"

fail() { echo "FAIL: $1" >&2; exit 1; }

# Clean slate.
rm -f "$BIN" "$BUILT_CHECKSUM"

# 1. First run builds the binary.
"$ENSURE"
[ -x "$BIN" ] || fail "expected estimatorctl to be built"
[ -f "$BUILT_CHECKSUM" ] || fail "expected .built-checksum to be written"

# 2. Second run with no source change does not rebuild.
mtime_before=$(stat -f '%m' "$BIN" 2>/dev/null || stat -c '%Y' "$BIN")
sleep 1
"$ENSURE"
mtime_after=$(stat -f '%m' "$BIN" 2>/dev/null || stat -c '%Y' "$BIN")
[ "$mtime_before" = "$mtime_after" ] || fail "expected no rebuild when checksum unchanged"

# 3. A changed .sync-checksum triggers a rebuild.
echo "different-checksum-value" > "$PLUGIN_ROOT/engine/.sync-checksum"
sleep 1
"$ENSURE"
mtime_rebuilt=$(stat -f '%m' "$BIN" 2>/dev/null || stat -c '%Y' "$BIN")
[ "$mtime_rebuilt" != "$mtime_after" ] || fail "expected rebuild when checksum changed"
[ "$(cat "$BUILT_CHECKSUM")" = "different-checksum-value" ] || fail "expected built-checksum to be updated"

# Restore the real checksum so other tests / builds aren't left broken.
( cd "$PLUGIN_ROOT/.." && ./scripts/sync-plugin-engine.sh >/dev/null )
"$ENSURE"

# 4. Missing go toolchain fails clearly. Strip only go's directory from
# PATH so dirname/cat/mkdir/etc. used by the script itself still resolve.
go_dir=$(dirname "$(command -v go)")
stripped_path=$(echo "$PATH" | tr ':' '\n' | grep -v -x "$go_dir" | tr '\n' ':')
stderr_file=$(mktemp)
if PATH="$stripped_path" "$ENSURE" 2>"$stderr_file"; then
  fail "expected failure with no go on PATH"
fi
grep -q "go toolchain not found" "$stderr_file" \
  || fail "expected a clear missing-toolchain message"
rm -f "$stderr_file"

echo "PASS: test_ensure_engine.sh"
