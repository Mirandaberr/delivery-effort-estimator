# Contract: `plugin/scripts/ensure-engine.sh`

## Behavior

1. If `command -v go` fails, print `ENGINE_ERROR: go toolchain not found on PATH` to
   stderr and exit 1. Caller (the skill) MUST stop, not continue, on exit 1.
2. Compute `cksum plugin/engine/.sync-checksum` (research.md "Engine
   freshness check").
3. If `plugin/bin/estimatorctl` exists and `plugin/bin/.built-checksum`
   matches step 2's value, exit 0 immediately (no rebuild) and print
   nothing.
4. Otherwise: run `go build -o plugin/bin/estimatorctl ./cmd/estimatorctl`
   from within `plugin/engine/`. On success, write step 2's checksum to
   `plugin/bin/.built-checksum` and exit 0. On build failure, print the
   compiler output to stderr and exit 1 — never leave a stale or partial
   binary in place silently.

## Output contract

- Exit `0`: `plugin/bin/estimatorctl` is present, built, and ready to run.
- Exit `1`: nothing usable is ready; stderr names the specific cause (missing
  Go, or the `go build` failure output). Callers must surface this to the
  user verbatim rather than swallowing it (FR-012).
- No output to stdout in the success path (keeps the skill's own reporting
  the single source of user-facing text — Constitution Principle VI framing
  lives in the skill, not scattered across scripts).
