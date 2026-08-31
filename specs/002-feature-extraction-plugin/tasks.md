---

description: "Task list for Feature Extraction Plugin"
---

# Tasks: Feature Extraction Plugin

**Input**: Design documents from `/specs/002-feature-extraction-plugin/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Included for the deterministic shell scripts (hook, feature-location,
engine build) since they're the only pieces that are meaningfully unit-testable
per plan.md's Testing strategy; the LLM-driven skill itself is validated via
quickstart.md's manual scenarios, not automated tests.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)

## Phase 1: Setup

**Purpose**: Scaffold the plugin as a self-contained, installable unit.

- [ ] T001 Create the plugin directory skeleton per plan.md Project Structure: `plugin/.claude-plugin/`, `plugin/hooks/`, `plugin/skills/estimate-feature/`, `plugin/scripts/`, `plugin/engine/`, `plugin/tests/`
- [ ] T002 [P] Write `plugin/.claude-plugin/plugin.json` (name, description, version `0.1.0`, author)
- [ ] T003 [P] Write `scripts/sync-plugin-engine.sh` at the repo root: copies `cmd/`, `internal/`, `go.mod`, `go.sum` into `plugin/engine/`, then writes `plugin/engine/.sync-checksum` (`cksum` over the copied files, per research.md)

**Checkpoint**: Plugin skeleton and the engine-sync mechanism exist; nothing functional yet.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: The shared machinery every user story depends on — engine packaging, feature-directory resolution, and their tests. No user story can be demoed until this phase passes.

- [ ] T004 Add `sync-plugin-engine` and `test-plugin` targets to the root `Makefile`; make `test-plugin` run every script in `plugin/tests/`; make `test` depend on `sync-plugin-engine` followed by a diff check that fails if `plugin/engine/` drifted from `cmd/`/`internal/`/`go.mod`/`go.sum`
- [ ] T005 Run `make sync-plugin-engine` to produce the first `plugin/engine/` snapshot and `plugin/engine/.sync-checksum`
- [ ] T006 [P] Implement `plugin/scripts/ensure-engine.sh` per contracts/ensure-engine.md (missing-`go` check, checksum-based rebuild skip, `go build -o plugin/bin/estimatorctl ./cmd/estimatorctl`, `plugin/bin/.built-checksum` marker)
- [ ] T007 [P] Implement `plugin/scripts/locate-feature.sh` per contracts/hook.md steps 2-3: given a file path, print the owning `specs/<id>-*/` directory if `spec.md` and `plan.md` are both present alongside it, otherwise exit non-zero with no output
- [ ] T008 [P] Write `plugin/tests/test_ensure_engine.sh`: first run builds the binary; second run with no source change does not rebuild (assert `plugin/bin/estimatorctl` mtime unchanged); a changed `.sync-checksum` triggers a rebuild; `PATH` stripped of `go` exits 1 with a message naming the missing toolchain
- [ ] T009 [P] Write `plugin/tests/test_locate_feature.sh`: accepts a real `specs/002-feature-extraction-plugin/tasks.md`-style path, rejects a `tasks.md` outside any `specs/*/` directory, rejects one missing a sibling `spec.md` or `plan.md`

**Checkpoint**: `make test-plugin` passes; the engine builds and is reused correctly; feature-directory resolution is correct and tested. All user stories can now build on this.

---

## Phase 3: User Story 1 - Automatic prediction right after planning (Priority: P1) 🎯 MVP

**Goal**: `tasks.md` being written triggers, in the same turn, a derived feature vector and a Delivery Effort prediction — no manual JSON authoring.

**Independent Test**: quickstart.md Scenario 1 (and Scenario 3 for the skip-safely edge).

### Implementation for User Story 1

- [ ] T010 [US1] Implement `plugin/scripts/on-tasks-md-change.sh` per contracts/hook.md: read stdin, extract `tool_input.file_path` via `grep`/`sed` (no `jq`, per research.md), call `locate-feature.sh`, on success print the `hookSpecificOutput.additionalContext` JSON, otherwise exit 0 silently
- [ ] T011 [US1] Write `plugin/hooks/hooks.json` with two `PostToolUse` entries (matcher `Write`, matcher `Edit`), both invoking `${CLAUDE_PLUGIN_ROOT}/scripts/on-tasks-md-change.sh`, per contracts/hook.md
- [ ] T012 [US1] Write `plugin/skills/estimate-feature/SKILL.md`: frontmatter (`user-invocable: true`, manual + hook-triggered per FR-009), the full research.md rubric table (all 7 dimensions with their anchors and primary signal sources), and the exact steps from contracts/skill.md (read files → score+justify → `ensure-engine.sh` → write `derived-features.json` → run `estimatorctl estimate` → save `estimation-record.json` → report as a prediction, never a commitment, per Constitution Principle VI)
- [ ] T013 [P] [US1] Write `plugin/tests/test_on_tasks_md_change.sh`: feeds synthetic `PostToolUse` stdin JSON for a valid `tasks.md` path and asserts the exact `additionalContext` JSON shape; feeds a non-`tasks.md` path and a `tasks.md` missing siblings, asserts silent no-op (empty stdout, exit 0) in both
- [ ] T014 [US1] Manually run quickstart.md Scenario 1 and Scenario 3 inside a real Claude Code session (`claude --plugin-dir plugin/`) against a throwaway Spec-Kit project; confirm the prediction appears unprompted after `tasks.md` is written, a second edit produces a second, non-overwriting `estimation/<timestamp>/` directory, and a non-Spec-Kit project / stray `tasks.md` produces no effect

**Checkpoint**: User Story 1 fully working end-to-end — this is the MVP.

---

## Phase 4: User Story 2 - Auditable derivation, not a black box (Priority: P2)

**Goal**: Every automatically derived feature vector carries a per-dimension justification a developer can trace back to spec/plan/tasks content.

**Independent Test**: quickstart.md Scenario 2.

### Implementation for User Story 2

- [ ] T015 [US2] In `plugin/skills/estimate-feature/SKILL.md` (extends T012), make explicit and non-skippable: every key in `features` must have a matching, non-empty `justifications` entry before `derived-features.json` is written (data-model.md validation rule) — including the conservative-score-plus-assumption-note path for missing/thin source sections (Edge Cases)
- [ ] T016 [US2] Manually run quickstart.md Scenario 2 against the output of T014's run: open `derived-features.json`, confirm all 7 justifications are present and each traces to a specific statement in that feature's `spec.md`/`plan.md`/`tasks.md`

**Checkpoint**: Derivations are auditable, not opaque numbers.

---

## Phase 5: User Story 3 - Drop the plugin into any Spec-Kit project (Priority: P3)

**Goal**: `plugin/` alone (no rest of this repo) works in a fresh Spec-Kit project with just a Go toolchain present.

**Independent Test**: quickstart.md Scenario 4.

### Implementation for User Story 3

- [ ] T017 [US3] Verify `plugin/engine/go.mod`'s module path and all imports resolve correctly when `plugin/engine/` is built from outside this repository (no accidental relative-path or repo-root assumptions leaking in from the sync in T005)
- [ ] T018 [P] [US3] Write `plugin/README.md`: how to install (`claude --plugin-dir <path>/plugin`), prerequisites (Go toolchain), and what happens on first use (engine build)
- [ ] T019 [US3] Manually run quickstart.md Scenario 4: copy only `plugin/` to a location with no other part of this repo, install it into a fresh Spec-Kit project, repeat Scenario 1, confirm identical behavior

**Checkpoint**: The plugin is genuinely portable, not repo-coupled.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [ ] T020 [P] Update root `README.md`: add a section describing the plugin, linking to `specs/002-feature-extraction-plugin/` and `plugin/README.md`
- [ ] T021 Run `make sync-plugin-engine && make test-plugin && make test`; fix anything failing
- [ ] T022 Run quickstart.md's full "Automated coverage" section end-to-end and confirm every listed expectation holds

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies.
- **Foundational (Phase 2)**: Depends on Setup — BLOCKS all user stories (US1's hook needs `locate-feature.sh`; US1's skill needs `ensure-engine.sh`).
- **User Story 1 (Phase 3)**: Depends on Foundational only. This is the MVP.
- **User Story 2 (Phase 4)**: Depends on US1's skill existing (T012) — it refines/hardens the same file (T015) rather than adding a new component.
- **User Story 3 (Phase 5)**: Depends on Foundational (the engine-sync mechanism) but is otherwise independent of US1/US2's runtime behavior — it can proceed as soon as T005 exists, though verifying it (T019) is most meaningful once US1 works end-to-end.
- **Polish (Phase 6)**: Depends on all three user stories being complete.

### Parallel Opportunities

- T002, T003 in Setup.
- T006, T007, T008, T009 in Foundational (distinct files).
- T013 in US1 can be written in parallel with T012 (different files), but should be run only after T010 exists.
- T018 in US3 can proceed in parallel with T017/T019.
- T020 in Polish is independent of T021/T022.

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1 (Setup) and Phase 2 (Foundational).
2. Complete Phase 3 (US1) — this alone delivers the feature's entire point
   (SDD §21 automatic integration).
3. **STOP and VALIDATE**: run quickstart.md Scenarios 1 and 3 for real.

### Incremental Delivery

1. Setup + Foundational → engine builds, feature-directory resolution works.
2. US1 → automatic estimation works end-to-end (MVP).
3. US2 → derivation is provably auditable, not just functional.
4. US3 → plugin is portable beyond this repository.
5. Polish → documented and drift-checked.
