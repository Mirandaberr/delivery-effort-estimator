# Tasks: Core Estimation Engine

**Input**: Design documents from `/specs/001-core-estimation-engine/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Included — quickstart.md already commits to specific `_test.go` files, and
Constitution Principle VII/VIII (predictions must be measurable, model must be
calibratable) make correctness of the deterministic formulas load-bearing.

**Organization**: Grouped by user story (spec.md) so each is independently testable.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependency on an incomplete task)
- **[Story]**: US1 / US2 / US3, per spec.md priorities

## Path Conventions

Idiomatic Go layout per plan.md: `cmd/`, `internal/` at repository root (no `src/`;
tests colocated as `_test.go` next to the code they cover, no separate `tests/` tree).

---

## Phase 1: Setup

**Purpose**: Project initialization

- [ ] T001 Initialize Go module at repo root: `go mod init github.com/jmirandev/delivery-effort-estimator` (creates `go.mod`, pinned to `go 1.27`)
- [ ] T002 [P] Add SQLite dependency: `go get modernc.org/sqlite` (updates `go.mod`/`go.sum`)
- [ ] T003 [P] Add a `Makefile` at repo root with `fmt`, `vet`, `test`, `build` targets wrapping `gofmt -l .`, `go vet ./...`, `go test ./...`, `go build ./...`

**Checkpoint**: `go build ./...` runs (with no source files yet, trivially succeeds) and dependency is resolvable.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Shared types, storage schema, and entrypoint skeletons every user story builds on

**⚠️ CRITICAL**: No user story task may start until this phase is complete

- [ ] T004 Create `EstimationFeatures` struct (7 fields per data-model.md) in `internal/estimation/features.go`, with a `Validate() (adjustments []string, err error)` method skeleton (real clamping/missing-field logic lands in Phase 3, T014)
- [ ] T005 [P] Create `Prediction`, `EstimationRecord`, `OutcomeRecord`, `ErrorReport`, `DimensionError` structs (data-model.md) in `internal/record/types.go`
- [ ] T006 Create `EstimationRepository` and `OutcomeRepository` interfaces in `internal/record/repository.go` — `EstimationRepository` exposes only `Save(EstimationRecord) error` and `Get(id string) (EstimationRecord, bool, error)` (no `Update`, per Constitution Principle VIII); `OutcomeRepository` exposes only `Save(OutcomeRecord) error` and `GetByEstimationID(id string) (OutcomeRecord, bool, error)`
- [ ] T007 Create append-only schema migration in `internal/storage/sqlite/migrations.go`: `estimation_records` and `outcome_records` tables per data-model.md, with `UNIQUE(estimation_id)` on `outcome_records` (FR-013)
- [ ] T008 Implement `internal/storage/sqlite/sqlite.go`: SQLite-backed implementation of both repository interfaces from T006, opening the DB file (default `./data/estimator.db`, creating the parent directory if missing) and running the T007 migration on startup
- [ ] T009 [P] Create `cmd/estimatorctl/main.go` skeleton: argument parsing that dispatches to `estimate` / `record-outcome` / `error-report` subcommands (each returning a placeholder "not implemented" error for now)
- [ ] T010 [P] Create `cmd/estimatord/main.go` skeleton: `net/http` server registering the four routes from `contracts/http.md` (each a placeholder handler returning 501 for now)

**Checkpoint**: Foundation ready — shared types, storage, and entrypoints exist; user story phases can now proceed.

---

## Phase 3: User Story 1 - Get a structured prediction from a feature set (Priority: P1) 🎯 MVP

**Goal**: Given a complete `EstimationFeatures` vector, return a fully-populated,
deterministic `Prediction` (FR-001–FR-007, FR-016–FR-018).

**Independent Test**: Run quickstart.md Scenario 1 — `estimatorctl estimate` twice
with the same input and confirm identical `prediction` content both times.

### Tests for User Story 1 ⚠️ write first, confirm they fail

- [ ] T011 [P] [US1] Table-driven tests for the deterministic formulas (HumanEffort, AIEffort, VerificationEffort, IntegrationEffort, DeliveryEffort, Confidence cap, Prediction Interval floor, ExpectedDuration, ExpectedAICost + zero-cost flag) from research.md, in `internal/estimation/model_test.go`
- [ ] T012 [P] [US1] Tests for `EstimationFeatures.Validate()`: missing-dimension rejection (FR-002) and out-of-range clamping with adjustment message (FR-003), in `internal/estimation/features_test.go`

### Implementation for User Story 1

- [ ] T013 [US1] Implement `Predict(EstimationFeatures) Prediction` in `internal/estimation/model.go` per research.md formulas, satisfying T011 (depends on T004, T005, T011)
- [ ] T014 [US1] Implement real clamping/missing-field logic in `EstimationFeatures.Validate()` in `internal/estimation/features.go`, satisfying T012 (depends on T012)
- [ ] T015 [US1] Implement `Service.Estimate(workItemID string, features EstimationFeatures) (EstimationRecord, error)` in `internal/record/service.go`: calls `Validate`, calls `Predict`, assembles an `EstimationRecord` (`model_version = "v1-linear"`, `calibration_version = "uncalibrated"`, `assumptions` populated from any clamps), persists via `EstimationRepository.Save` (depends on T013, T014, T006, T008)
- [ ] T016 [US1] Wire `estimatorctl estimate --work-item --features` to `Service.Estimate`, printing the resulting `EstimationRecord` as JSON on success or a structured error (`missing_feature`/`invalid_json`) on stderr with exit 1, in `cmd/estimatorctl/main.go` (depends on T015, T009; contracts/cli.md)
- [ ] T017 [US1] Wire `POST /work-items/{workItemId}/estimations` to `Service.Estimate`, returning `201` + `EstimationRecord` JSON or `400` + structured error, in `cmd/estimatord/main.go` (depends on T015, T010; contracts/http.md)

**Checkpoint**: User Story 1 fully functional — quickstart.md Scenario 1 passes.

---

## Phase 4: User Story 2 - Every prediction is a reproducible, versioned record (Priority: P2)

**Goal**: A stored `EstimationRecord` can always be fetched back unchanged, and
re-estimating a work item never overwrites a prior record (FR-008–FR-010).

**Independent Test**: Run quickstart.md Scenario 2 — fetch a record by id twice and
confirm byte-identical `prediction`/`model_version`; re-estimate the same work item
and confirm two distinct ids exist.

### Tests for User Story 2 ⚠️ write first, confirm they fail

- [ ] T018 [P] [US2] Test in `internal/record/service_test.go`: two `Estimate` calls for the same `workItemID` produce two `EstimationRecord`s with different `id`s and both remain independently retrievable
- [ ] T019 [P] [US2] Test in `internal/storage/sqlite/sqlite_test.go`: `EstimationRepository` has no method capable of mutating a saved row (compile-time interface check) and `Get` after `Save` returns byte-identical content

### Implementation for User Story 2

- [ ] T020 [US2] Implement `Service.GetEstimation(id string) (EstimationRecord, error)` in `internal/record/service.go`, returning a structured `unknown_estimation_id` error when absent (depends on T006, T008)
- [ ] T021 [US2] Wire `GET /estimations/{estimationId}` to `Service.GetEstimation` (`200` + record, or `404`) in `cmd/estimatord/main.go` (depends on T020, T017)
- [ ] T022 [US2] Add `estimatorctl get --estimation-id` command wired to `Service.GetEstimation`, for the manual reproducibility check in quickstart.md Scenario 2, in `cmd/estimatorctl/main.go` (depends on T020, T016)

**Checkpoint**: User Stories 1 and 2 both work independently — quickstart.md Scenario 2 passes.

---

## Phase 5: User Story 3 - Compare a prediction against what actually happened (Priority: P3)

**Goal**: Accept an `OutcomeRecord` linked to an existing `EstimationRecord` and
compute an on-demand `ErrorReport` (FR-011–FR-015).

**Independent Test**: Run quickstart.md Scenario 3 — record an outcome, fetch its
error report, then confirm a second `record-outcome` for the same id is rejected.

### Tests for User Story 3 ⚠️ write first, confirm they fail

- [ ] T023 [P] [US3] Tests in `internal/storage/sqlite/sqlite_test.go`: `OutcomeRepository.Save` rejects a second outcome for the same `estimation_id` (unique constraint from T007) and `Service`-level unknown-`estimation_id` handling
- [ ] T024 [P] [US3] Tests in `internal/estimation/errorreport_test.go`: `Compute(record, outcome)` returns correct `absolute_error`/`relative_error`/`bias` per dimension listed in data-model.md, including `relative_error == nil` when `predicted == 0`

### Implementation for User Story 3

- [ ] T025 [US3] Implement `Compute(EstimationRecord, OutcomeRecord) ErrorReport` in `internal/estimation/errorreport.go` per research.md's raw-delta approach, satisfying T024 (depends on T024, T005)
- [ ] T026 [US3] Implement `Service.RecordOutcome(estimationID string, outcome OutcomeRecord) (OutcomeRecord, error)` in `internal/record/service.go`: rejects unknown `estimationID` (`unknown_estimation_id`) and an existing outcome (`outcome_already_recorded`), else persists via `OutcomeRepository.Save` (depends on T023, T006, T008)
- [ ] T027 [US3] Implement `Service.ErrorReport(estimationID string) (ErrorReport, error)` combining the stored `EstimationRecord`, its `OutcomeRecord` (error if none — `no_outcome_recorded`), and `estimation.Compute` (depends on T025, T026)
- [ ] T028 [US3] Wire `estimatorctl record-outcome` and `estimatorctl error-report` subcommands to `Service.RecordOutcome`/`Service.ErrorReport` in `cmd/estimatorctl/main.go` (depends on T026, T027, T009; contracts/cli.md)
- [ ] T029 [US3] Wire `POST /estimations/{id}/outcome` (`201`/`404`/`409`) and `GET /estimations/{id}/error-report` (`200`/`404`) in `cmd/estimatord/main.go` (depends on T026, T027, T010; contracts/http.md)

**Checkpoint**: All three user stories independently functional — full quickstart.md passes end-to-end.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [ ] T030 [P] Write `README.md` at repo root: what this is, how to build (`make build`), how to run quickstart.md, link to `docs/sdd-estimate.md` and `specs/001-core-estimation-engine/`
- [ ] T031 Run `make fmt vet test` from repo root; fix any findings until clean
- [ ] T032 Manually execute all three quickstart.md scenarios end-to-end against a fresh `./data/estimator.db` and fix any drift found between docs and behavior
- [ ] T033 [P] Add a top-level `internal/record/repository_test.go` asserting (via `go vet`-checkable interface satisfaction, not reflection) that no repository interface exposes an update/delete method — a lightweight guard for Constitution Principle VIII surviving future edits

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: no dependencies
- **Foundational (Phase 2)**: depends on Setup — blocks all user stories
- **User Story 1 (Phase 3)**: depends on Foundational only
- **User Story 2 (Phase 4)**: depends on Foundational + reuses `Service`/repository from US1 (T006, T008) but adds no new dependency on US1's tasks beyond those shared foundations — independently testable once Foundational is done
- **User Story 3 (Phase 5)**: depends on Foundational + shared `Service`/repository (T006, T008); independently testable the same way
- **Polish (Phase 6)**: depends on whichever user stories are in scope for a given delivery

### Within Each User Story

- Tests (T011/T012, T018/T019, T023/T024) are written first and must fail before their implementation tasks
- `internal/estimation` (pure logic) before `internal/record` (orchestration) before `cmd/*` (delivery)
- Story complete and checkpointed before moving to the next priority

### Parallel Opportunities

- T002, T003 in Setup
- T005, T009, T010 in Foundational (different files from T004/T006/T007/T008)
- T011 + T012 (US1 tests, different files)
- T018 + T019 (US2 tests, different files)
- T023 + T024 (US3 tests, different files)
- T030 + T033 in Polish

---

## Parallel Example: User Story 1

```bash
# Tests, together:
Task: "Table-driven tests for deterministic formulas in internal/estimation/model_test.go"
Task: "Tests for EstimationFeatures.Validate() in internal/estimation/features_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Phase 1 (Setup) → Phase 2 (Foundational) → Phase 3 (US1)
2. **STOP and VALIDATE**: run quickstart.md Scenario 1
3. This alone already satisfies the SDD's core promise (§1): a quantitative Delivery
   Effort prediction from available information — everything else is refinement.

### Incremental Delivery

1. Setup + Foundational → foundation ready
2. + US1 → validate Scenario 1 (MVP)
3. + US2 → validate Scenario 2 (reproducibility)
4. + US3 → validate Scenario 3 (prediction vs. actual)
5. + Polish → `make fmt vet test` clean, README in place

---

## Notes

- [P] tasks touch different files and have no incomplete-task dependency
- Every task lists an exact file path per the SKILL.md format requirement
- Run `go test ./...` after each phase, not just at the end
- Commit after each checkpoint (end of a phase), not after every single task
