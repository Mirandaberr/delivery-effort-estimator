# Phase 0 Research: Feature Extraction Plugin

Two decisions were already made with the user before this plan was written
(trigger mechanism, engine distribution strategy) and are recorded here as
resolved decisions rather than open questions. The remaining decisions —
the derivation rubric, file layout, and engine-freshness check — are new to
this feature and resolved below.

## Decision: Automatic trigger via a PostToolUse hook that injects context, not a hook that reasons itself

**Rationale**: A Claude Code hook is a deterministic shell script; it cannot
perform the LLM reasoning FR-002/FR-003 require (mapping spec/plan/tasks
prose to 7 justified scores). The hook's only job is *detection*: fire when
a `tasks.md` write lands inside a `specs/<id>-*/` directory that also has
`spec.md` and `plan.md`, resolve which feature directory that is, and hand
control back to Claude by injecting an instruction to run the
`estimate-feature` skill for that directory. The actual derivation happens
in a normal Claude turn, inside the skill — which is also why the skill
remains independently, manually invocable (FR-009): the hook is just one way
to reach the same skill.

**Verified contract** (Claude Code hooks reference, confirmed before writing
this decision, not assumed):
- Event: `PostToolUse`, matcher `Write` and a separate matcher `Edit` (two
  hook entries, not one combined matcher — see contracts/hook.md for why).
- Hook stdin JSON includes `tool_name`, `tool_input.file_path`,
  `hook_event_name`, `cwd`. `on-tasks-md-change.sh` reads
  `tool_input.file_path` from this payload.
- Hook stdout JSON to inject context:
  `{"hookSpecificOutput": {"hookEventName": "PostToolUse", "additionalContext": "..."}}`.
  PostToolUse cannot block (the tool already ran); `additionalContext` is
  the only mechanism this feature needs.
- Injected context reaches Claude on its next model-generation step, which
  in practice is still within the same assistant turn that just wrote
  `tasks.md` (the turn continues generating before yielding to the user) —
  so the automatic estimation visibly happens right after planning
  completes, satisfying FR-015, not on some later, disconnected turn.
- Plugin `hooks/hooks.json` supports an `if` field with permission-rule
  syntax (e.g. `"if": "Write(specs/*/tasks.md)"`) for coarse path
  pre-filtering, but its exact matching semantics are not fully documented.
  This feature does not rely on `if` for correctness — `on-tasks-md-change.sh`
  independently re-validates the path itself (FR-014), so a mismatch in
  `if` semantics degrades to "hook runs slightly more often, script still
  filters correctly," never to a false positive reaching the skill.

**Alternatives considered**:
- A `Stop` hook — rejected (this was the initial, incorrect research result):
  fires after *every* assistant turn regardless of what changed, and still
  can't reason about content itself.
- Having the hook shell out to `claude -p "extract features..."` as a
  recursive headless call — rejected: adds a second, separate LLM
  invocation and cost when the conversation already in progress can do the
  same reasoning for free in its next turn; also loses the user-visible
  trail (User Story 2) of *when* and *why* the estimation ran.
- A `PreToolUse` hook — rejected: firing before the tool runs means
  `tasks.md` would not exist on disk yet for the skill to read.
- Relying solely on the `if` field for path filtering, with no in-script
  check — rejected: its glob/multi-tool semantics are documented
  ambiguously; correctness of FR-014 must not depend on an unverified
  detail, so the script re-checks regardless of what `if` already filtered.

## Decision: Feature-directory resolution is path-convention-based, not conversation-memory-based

**Rationale**: FR-014 requires correctness even with multiple in-flight
Spec-Kit features. `locate-feature.sh` resolves strictly from the file path
Claude Code reports for the write (`specs/<id>-<slug>/tasks.md`) plus a
same-directory check for sibling `spec.md`/`plan.md` — never from "the
feature we were just discussing," which conversation context can get wrong
across topic switches. A `tasks.md` write anywhere that doesn't match
`specs/*/tasks.md` with those two siblings present is ignored (Edge Case:
stray file named `tasks.md`).

**Alternatives considered**:
- Trusting `.specify/feature.json`'s `feature_directory` — rejected as the
  sole source: it reflects the *last* feature `/speckit-specify` touched,
  which can be stale if the user edits an older feature's `tasks.md` by
  hand; the write path itself is ground truth and is used instead.

## Decision: Derivation rubric — concrete anchors per dimension, not open-ended judgment

**Rationale**: Constitution Principle I requires evidence, not intuition, and
FR-013 requires the rubric to be *documented and consistent*. Each of the 7
`EstimationFeatures` dimensions (specs/001) is given a 5-point anchor scale
(0.0 / 0.25 / 0.5 / 0.75 / 1.0) tied to concrete, checkable signals in
`spec.md`/`plan.md`/`tasks.md`, mirroring which output each dimension drives
most in the specs/001 formulas (research.md there) — e.g.
`integration_complexity` is anchored on *count of external touchpoints*
because it dominates `IntegrationEffort` (weight 0.70) in the engine.

| Dimension | Primary signal source | 0.0 anchor | 0.5 anchor | 1.0 anchor |
|---|---|---|---|---|
| `context_complexity` | plan.md "components/modules affected"; spec.md scenario breadth | Change confined to one already-understood file/module | Spans a handful of modules within one service | Requires understanding large or unfamiliar parts of the system across services |
| `domain_complexity` | spec.md problem description, Key Entities | CRUD-like, no domain rules beyond basic validation | Moderate business rules or a non-trivial state machine | Deep domain modeling (compliance, financial/legal rules, complex multi-actor workflows) |
| `integration_complexity` | plan.md components/APIs/events/DBs/infra/external systems list | No external touchpoints (pure internal logic) | 2-3 touchpoints (e.g., one API + one DB) | 4+ touchpoints or any external/third-party system dependency |
| `verification_complexity` | plan.md testing strategy; spec.md acceptance scenarios & edge cases | 1-2 straightforward acceptance scenarios, no regression risk | Several scenarios plus at least one edge case needing dedicated test design | Broad scenario coverage, regression-sensitive area, and/or manual validation required |
| `human_decision_load` | spec.md Assumptions/[NEEDS CLARIFICATION] history; plan.md pending decisions | Fully unambiguous requirements, no judgment calls left | A few reasonable defaults were chosen among real alternatives | Multiple open architectural/product trade-offs the human must actively resolve |
| `ai_execution_complexity` | tasks.md task count/breadth; plan.md structure decision | A handful of small, mechanical tasks | Moderate task count spanning several files with some non-trivial logic | Large or intricate task set requiring iteration, multi-file reasoning, or nontrivial tool/build loops |
| `uncertainty` | spec.md Assumptions section; plan.md unresolved dependencies/unknowns | No assumptions recorded, no unresolved dependencies | A few named assumptions, none blocking | Several unresolved unknowns or assumptions the plan itself flags as risky |

Scores between anchors are interpolated by judgment; every score MUST be
paired with a one- to two-sentence justification quoting or paraphrasing the
specific spec/plan/tasks content it came from (FR-003), never left as a bare
number. When a section the rubric depends on is missing or too thin to
support a judgment, the dimension is scored conservatively (toward the 0.0
anchor for effort-type dimensions, toward the 1.0 anchor for `uncertainty`)
and the justification says so explicitly — this is the "assumption" path
the Edge Cases section describes, deliberately mirroring how the engine's
own `Clamp()` documents out-of-range input (specs/001).

**Alternatives considered**:
- A single free-text prompt ("estimate these 7 dimensions from 0 to 1") with
  no anchors — rejected: produces inconsistent scoring across features and
  fails Principle I's evidence requirement; nothing to audit against.
- Deterministic keyword/heuristic scoring (e.g., counting words, regex over
  plan.md) instead of LLM judgment — rejected: the SDD (§10) explicitly
  allows LLM-human conversation as a feature source and a fixed heuristic
  would badly misjudge domain/context complexity that requires actually
  understanding the text, not just its shape.

## Decision: Output files live under a per-feature `estimation/` subdirectory, one timestamped pair per run

**Rationale**: FR-004/FR-006/FR-007 require persistence inside the feature's
own spec directory, without ever overwriting a prior run. Rather than
inventing a versioned-filename scheme at the root of `specs/<id>-*/`
(cluttering it and risking collisions with Spec-Kit's own files), each run
writes into `specs/<id>-*/estimation/<UTC-timestamp>/`:
- `derived-features.json` — the 7 scores + justifications + source file
  list (User Story 2's audit trail).
- `estimation-record.json` — the exact JSON `estimatorctl estimate` printed
  (specs/001/contracts/cli.md), saved verbatim.

A directory-per-run makes "never overwrite" structurally true (a new
directory, not a new field in a shared file) and keeps all Spec-Kit's own
files untouched.

**Alternatives considered**:
- Storing only in `estimatorctl`'s SQLite database, with nothing written
  next to the spec — rejected: User Story 2 requires the derivation to be
  discoverable by *reading the feature directory*, not by knowing to query a
  separate database file.

## Decision: Hook stdin parsing uses `grep`/`sed`, not `jq`

**Rationale**: `jq` is not guaranteed present on a fresh macOS machine (not
bundled with the OS), and requiring the user to install it before the
plugin works at all would contradict "lightest and most efficient" and
FR-010's self-contained-install goal. The stdin payload's shape is fixed and
documented (previous decision), so `on-tasks-md-change.sh` extracts
`tool_input.file_path` with a single `grep -o '"file_path"[[:space:]]*:[[:space:]]*"[^"]*"'`
plus a trailing `sed` to strip the key — sufficient for a flat, known field,
without adding a JSON-parser dependency.

**Alternatives considered**:
- Depending on `jq` — rejected for the portability reason above; revisit if
  a future hook needs to parse deeply nested or variable-shape JSON, where
  `grep`/`sed` would stop being reliable.

## Decision: Engine freshness check via a checksum marker, not a semantic version bump per sync

**Rationale**: `plugin/engine/` (Project Structure) is a mechanical copy of
`cmd/estimatorctl`, `internal/*`, `go.mod`, `go.sum`. `scripts/sync-plugin-engine.sh`
writes a `plugin/engine/.sync-checksum` file (`cksum` over the concatenated
synced files — POSIX, no extra dependency, portable across macOS/Linux).
`plugin/scripts/ensure-engine.sh` rebuilds `plugin/bin/estimatorctl` only
when that checksum differs from the one recorded at the binary's last build
(`plugin/bin/.built-checksum`), keeping FR-011's "reuse on subsequent runs"
cheap while still catching real engine changes (FR-011's staleness clause).

**Alternatives considered**:
- Rebuilding on every invocation — rejected: defeats FR-011's "reuse without
  rebuilding" and adds unnecessary latency to every estimation.
- Comparing file modification times — rejected: not reliably preserved by
  `git clone`/plugin distribution, unlike content-based checksums.

## Graphify dependency-graph check (Constitution Principle X)

`graphify explain internal_record_service` and
`graphify explain internal_estimation_features_go_estimation_estimationfeatures`
were run before writing this plan (see plan.md's Constitution Check) and
confirmed `EstimationFeatures` has no existing Claude-Code-aware caller —
this feature is additive glue with zero required changes to
`internal/estimation`, `internal/record`, or `internal/storage/sqlite`.
