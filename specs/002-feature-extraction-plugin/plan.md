# Implementation Plan: Feature Extraction Plugin

**Branch**: `002-feature-extraction-plugin` | **Date**: 2026-08-31 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/002-feature-extraction-plugin/spec.md`

**Note**: This template is filled in by the `/speckit-plan` command; its definition describes the execution workflow.

## Summary

A Claude Code plugin that closes the loop described in the SDD (§21): after a
Spec-Kit feature is specified and planned (`tasks.md` exists), Claude — inside
the same conversation, following a documented rubric — derives the 7-dimension
`EstimationFeatures` vector from `spec.md`/`plan.md`/`tasks.md`, then shells
out to the existing, unmodified `estimatorctl` CLI (specs/001) to obtain the
actual prediction. The plugin never computes Delivery Effort itself (Principle
IX) — it only produces the features *in*, and persists both the derivation and
the resulting record next to the spec, append-only.

## Technical Context

**Language/Version**: POSIX `sh` for the hook and build scripts (matches this
project's existing `.specify/init-options.json` `"script": "sh"` convention);
Markdown+YAML frontmatter for the plugin manifest and skill; the estimation
logic itself is unchanged Go 1.27 (specs/001), packaged, not rewritten.

**Primary Dependencies**: Claude Code's plugin runtime (hooks, skills,
`bin/`-on-PATH convention); a local Go toolchain, required only to build the
bundled engine on first use in a target project (per the user's confirmed
"compile on install" decision); this repo's own `cmd/estimatorctl`,
`internal/estimation`, `internal/record`, `internal/storage/sqlite` as the
engine being packaged (reused as-is, not duplicated logic).

**Storage**: No new database. The plugin writes plain files into the target
feature's own `specs/<id>-*/` directory (`derived-features*.json`,
`estimation-record*.json`); the actual `EstimationRecord`/`OutcomeRecord`
persistence continues to be `estimatorctl`'s embedded SQLite store
(specs/001), untouched by this feature.

**Testing**: A `sh`-based test harness (`plugin/tests/`) for the two
deterministic pieces — the hook's path-matching/feature-directory resolution,
and the engine build/reuse script — run via `make test-plugin`. The
LLM-driven derivation step is not unit-testable; it is validated the same way
specs/001 was validated beyond `go test`: a manual, scripted quickstart run
against a real Claude Code session (quickstart.md).

**Target Platform**: Same as specs/001 (macOS/Linux; Windows via WSL) —
anywhere Claude Code and a Go toolchain run.

**Project Type**: Claude Code plugin (single self-contained unit at
`plugin/`), installable into any Spec-Kit-enabled project via
`claude --plugin-dir plugin/`. Not a library/service in its own right.

**Performance Goals**: Engine binary build is a one-time cost per plugin
version (a few seconds); each subsequent extraction+estimation completes
within the same conversation turn it was triggered in.

**Constraints**:
- MUST NOT duplicate or reimplement the deterministic estimation model
  (Estimation Engine Constraints) — the plugin only ever calls
  `estimatorctl`, never recomputes Delivery Effort/Confidence/etc. itself.
- MUST NOT block or break the user's normal workflow when Spec Kit isn't
  initialized or the Go toolchain is missing (SC-005) — failures degrade to a
  skipped/explained no-op, never a hard stop.
- MUST follow the append-only rule already enforced by the engine: a second
  derivation/estimation for the same feature directory is an additional file,
  never an overwrite (FR-007).
- The engine source packaged inside `plugin/engine/` MUST stay in sync with
  the canonical source at the repo root (`cmd/`, `internal/`, `go.mod`,
  `go.sum`) — see Project Structure below; drift is a build-time defect, not
  an acceptable variance.

**Scale/Scope**: One feature directory per invocation (the one whose
`tasks.md` just changed). Cross-feature analytics/aggregation is explicitly
out of scope here (spec.md Assumptions) — a future feature, not this one.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Per Principle X, the graphify dependency graph was consulted before writing
this plan (`graphify explain internal_record_service`,
`graphify explain internal_estimation_features_go_estimation_estimationfeatures`)
and cross-checked against the actual files: `EstimationFeatures` is consumed
only by `internal/estimation` and `internal/record` today, with no existing
Claude-Code-aware caller — confirming this plugin is purely additive glue and
requires zero changes to the existing engine packages.

| Principle / Constraint | Assessment | Status |
|---|---|---|
| I. Evidence Over Intuition | Derivation is grounded in the actual `spec.md`/`plan.md`/`tasks.md` content via a documented rubric (research.md), not free intuition. | PASS |
| II. Effort Is Not Duration | Unaffected — the plugin does not touch the engine's output shape. | PASS |
| III. Effort Is Not Cost | Unaffected, same reason. | PASS |
| IV. Business Value Is Not Delivery Effort | Plugin reads spec/plan/tasks content only for the 7 defined dimensions; it does not feed business-value language into the vector. | PASS |
| V. Human/AI Effort Independent | Unaffected — engine formulas untouched. | PASS |
| VI. Prediction, Not Commitment | The skill's user-facing summary (FR-015) MUST present the result as a prediction, matching the engine's own framing — a plan requirement, not just an assumption. | PASS (tracked as a task) |
| VII. Predictions Must Be Measurable | The plugin persists the estimation via the existing `EstimationRecord` mechanism, so it remains comparable to a future `OutcomeRecord` exactly as specs/001 designed. | PASS |
| VIII. Model Must Be Calibratable | Unaffected — no model logic here. | PASS |
| IX. LLMs Are Not the Estimation Engine | **Central design constraint of this feature.** Claude derives the *features* (non-reproducible LLM reasoning, clearly labeled and justified); the *quantitative* prediction is always produced by the deterministic `estimatorctl` binary, never by the LLM. The plugin's skill instructions MUST NOT ask Claude to state a Delivery Effort/Confidence number itself. | PASS (tracked as a task) |
| X. Consult the Dependency Graph | Done above. | PASS |
| Engine independence (Estimation Engine Constraints) | The plugin calls `estimatorctl` through its existing CLI contract (specs/001/contracts/cli.md) as an external process; `internal/estimation`/`internal/record` are not imported into any Claude-Code-aware code. | PASS |
| Spec-Kit as first specification integration | This feature literally is that integration (docs/sdd-estimate.md §6). | PASS |
| Planning is a primary evidence source | The rubric explicitly reads `plan.md` (not just `spec.md`) for integration/verification/uncertainty signals. | PASS |

No violations. Complexity Tracking is empty.

## Project Structure

### Documentation (this feature)

```text
specs/002-feature-extraction-plugin/
├── plan.md              # This file (/speckit-plan command output)
├── research.md          # Phase 0 output (/speckit-plan command)
├── data-model.md        # Phase 1 output (/speckit-plan command)
├── quickstart.md        # Phase 1 output (/speckit-plan command)
├── contracts/           # Phase 1 output (/speckit-plan command)
└── tasks.md             # Phase 2 output (/speckit-tasks command - NOT created by /speckit-plan)
```

### Source Code (repository root)

```text
plugin/                          # the installable unit (claude --plugin-dir plugin/)
├── .claude-plugin/
│   └── plugin.json               # manifest: name, description, version, author
├── hooks/
│   └── hooks.json                 # PostToolUse hook, matches specs/*/tasks.md writes
├── skills/
│   └── estimate-feature/
│       └── SKILL.md               # the rubric + steps Claude follows (this is the "code"
│                                   # for FR-002/003/013 — reasoning happens here, not in a script)
├── scripts/
│   ├── locate-feature.sh          # resolves the owning specs/<id>-*/ dir from a changed path (FR-014)
│   ├── on-tasks-md-change.sh      # PostToolUse hook entrypoint; calls locate-feature.sh,
│   │                               # emits additionalContext telling Claude to run the skill
│   └── ensure-engine.sh           # builds plugin/bin/estimatorctl from plugin/engine/ if missing
│                                   # or stale; clear error if `go` is unavailable (FR-011/012)
├── engine/                        # synced copy of cmd/estimatorctl + internal/* + go.mod/go.sum —
│                                   # see "Engine sync" below; never hand-edited
└── tests/
    ├── test_locate_feature.sh
    └── test_ensure_engine.sh

scripts/
└── sync-plugin-engine.sh          # repo-root script: rsyncs cmd/, internal/, go.mod, go.sum
                                    # into plugin/engine/; `make sync-plugin-engine` runs it,
                                    # and `make test` fails if plugin/engine/ has drifted
```

**Structure Decision**: The plugin lives at the repository root as a
self-contained `plugin/` directory so it can be pointed to directly with
`claude --plugin-dir plugin/` from any project, per User Story 3 — it does
not need this whole repository checked out to be usable elsewhere.

**Engine sync, not a rewrite**: `internal/estimation`, `internal/record`,
`internal/storage/sqlite`, and `cmd/estimatorctl` (specs/001) are the single
source of truth. `plugin/engine/` is a mechanically synced copy (via
`scripts/sync-plugin-engine.sh`), never edited by hand, so the plugin can
`go build` a standalone binary without requiring the target project to also
contain this repository. `make test` verifies the copy is byte-identical to
the canonical source, so drift is caught immediately rather than trusted.

## Complexity Tracking

> No violations — this section is intentionally empty.

## Post-Phase 1 Constitution Re-check

Design artifacts (research.md, data-model.md, contracts/) confirm the
Constitution Check table above still holds with no new violations:
- The verified hook contract (research.md) keeps all LLM reasoning inside
  the skill, never the hook — Principle IX intact.
- `contracts/skill.md` step 7 requires the user-facing summary to be framed
  as a prediction — Principle VI now has a concrete, checkable requirement,
  not just an intent.
- `contracts/ensure-engine.md` and `data-model.md`'s postconditions ensure a
  failed/missing engine build produces no partial `estimation/<timestamp>/`
  output — SC-005 (never break the user's workflow) is enforced structurally,
  not just documented.

No changes to Technical Context or Project Structure were needed after
Phase 1 design.
