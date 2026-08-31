# Quickstart: Feature Extraction Plugin

Two layers of validation, matching what's actually scriptable: the
deterministic pieces (hook detection, engine build/reuse, engine sync) are
validated with plain shell commands; the end-to-end automatic-trigger
behavior requires a real Claude Code session, since it depends on the
plugin runtime and an LLM turn (User Story 1/2).

## Prerequisites

- This repo built once (`make build`) so `internal/estimation`,
  `internal/record`, etc. exist to sync from.
- Go on `PATH` (for the "happy path" scenarios below).

## Automated coverage (scriptable)

```bash
make sync-plugin-engine   # copies cmd/, internal/, go.mod, go.sum into plugin/engine/
make test-plugin          # runs plugin/tests/*.sh
```

Expected: `plugin/engine/.sync-checksum` is written; all `plugin/tests/`
scripts pass, covering:
- `locate-feature.sh` accepts `specs/002-feature-extraction-plugin/tasks.md`
  (spec.md and plan.md exist there) and rejects a `tasks.md` created outside
  any `specs/*/` directory, or one whose siblings are missing.
- `ensure-engine.sh` builds `plugin/bin/estimatorctl` on first run, and does
  **not** rebuild (verified by checking the binary's mtime is unchanged) on
  a second run with no source changes.
- `ensure-engine.sh` rebuilds after `make sync-plugin-engine` runs again
  following a real change under `internal/` (checksum differs).
- `ensure-engine.sh` fails with a clear message when `PATH` is temporarily
  stripped of `go` (`PATH=/usr/bin:/bin plugin/scripts/ensure-engine.sh`
  should exit 1 and name the missing toolchain, not hang or crash).

## Scenario 1 — Automatic estimation after planning (User Story 1)

**Manual, inside Claude Code**, since this exercises the real plugin
runtime:

1. In a separate, throwaway Spec-Kit-initialized project, start Claude Code
   with `claude --plugin-dir <path-to-this-repo>/plugin`.
2. Run `/speckit.specify`, `/speckit.plan`, `/speckit.tasks` for a small
   feature.
3. **Expected**: within the same turn `tasks.md` is created, Claude reports
   a Delivery Effort prediction for that feature, unprompted, and
   `specs/<id>-*/estimation/<timestamp>/derived-features.json` +
   `estimation-record.json` exist on disk.
4. Edit `tasks.md` again (e.g., add a task) and save.
   **Expected**: a *second* `estimation/<timestamp>/` directory appears;
   the first one is untouched (FR-007 / SC-003).

## Scenario 2 — Derivation is auditable (User Story 2)

1. Open the `derived-features.json` written in Scenario 1.
   **Expected**: all 7 dimensions present, each with a non-empty
   `justifications` entry referencing something concrete from that feature's
   `spec.md`/`plan.md`/`tasks.md`.

## Scenario 3 — Skipped, not broken, when preconditions aren't met (SC-005)

1. In a project with no `.specify/` directory at all, create any `tasks.md`
   file. **Expected**: nothing happens — no error, no unexpected file
   writes, normal Claude Code usage is unaffected.
2. Temporarily rename `plugin/engine/` (simulating a broken install) and
   trigger the skill manually on a real feature directory. **Expected**: a
   clear error naming the problem, no `estimation/<timestamp>/` directory is
   created (contracts/skill.md postconditions), and the user's conversation
   is otherwise unaffected.

## Scenario 4 — Portable install (User Story 3)

1. Copy just the `plugin/` directory to a machine with Go installed but no
   copy of this repository.
2. `claude --plugin-dir <copied-path>` in a fresh Spec-Kit project.
3. Repeat Scenario 1. **Expected**: identical result — `plugin/engine/`
   being a synced, self-contained copy means the target project needs
   nothing beyond the `plugin/` directory itself and a Go toolchain.
