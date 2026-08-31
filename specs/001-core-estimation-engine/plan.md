# Implementation Plan: Core Estimation Engine

**Branch**: `001-core-estimation-engine` | **Date**: 2026-08-31 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/001-core-estimation-engine/spec.md`

## Summary

Build a deterministic Estimation Engine that accepts a normalized `EstimationFeatures`
vector for a work item and returns a structured `Prediction` (Delivery Effort in DEU,
Confidence, Prediction Interval, Human/AI/Verification/Integration Effort, Risk,
Expected Duration, Expected AI Cost), persists it as an immutable, versioned
`EstimationRecord`, accepts a linked `OutcomeRecord` after delivery, and computes an
`ErrorReport` (absolute/relative error, bias) between the two on demand. No feature
extraction, no automatic calibration, no Specification/Planning layer — this is the
smallest engine that can be called by everything else the SDD describes later.

Technical approach: a small Go module with a pure, side-effect-free `estimation`
package holding the weighted deterministic model, an embedded SQLite store behind a
narrow repository interface for append-only records, and a thin CLI + JSON-over-HTTP
wrapper so both humans and tools (Claude Code, CI) can call it without depending on any
particular LLM provider or language runtime.

## Technical Context

**Language/Version**: Go 1.27 (`go.mod` pinned to `go 1.27`)

**Primary Dependencies**: Go standard library only for the domain/engine logic
(`net/http`, `encoding/json`, `database/sql`); `modernc.org/sqlite` (pure-Go SQLite
driver, no CGO) as the sole third-party dependency, for portability of the single
static binary. No web framework, no ORM — deliberately, per "as light and efficient as
possible."

**Storage**: SQLite, one embedded file (default `./data/estimator.db`, overridable via
config), accessed only through `EstimationRepository` / `OutcomeRepository` interfaces
in `internal/record`. Swapping to PostgreSQL later means writing a second
implementation of those interfaces — the `estimation` engine package never touches
storage directly.

**Testing**: `go test ./...` (standard library `testing`). Table-driven unit tests for
the deterministic engine (`internal/estimation`), covering the domain boundaries in
spec.md's Edge Cases (missing dimension, out-of-range clamping, zero-cost flag).
Repository tests run against a temp-file SQLite DB per test (no mocking of SQL).

**Target Platform**: Single static binary, cross-compiled for Linux/macOS server or
local CLI use; runs with zero external services besides its own SQLite file.

**Project Type**: Single project — library core (`internal/estimation`,
`internal/record`) + thin CLI (`cmd/estimatorctl`) + thin HTTP JSON wrapper
(`cmd/estimatord`). (Template "Option 1: Single project", adapted to idiomatic Go
`cmd/`/`internal/` layout instead of `src/`/`lib/` — see Project Structure below.)

**Performance Goals**: A single estimation (pure arithmetic + one SQLite insert)
completes in well under 1s (SC-001); realistically low-single-digit milliseconds, with
no network calls on the prediction path.

**Constraints**: No LLM call, no network dependency, and no ML runtime on the
prediction path (constitution: Estimation Engine Constraints, FR-005/FR-006/FR-007).
Records are append-only — no UPDATE/DELETE path for a stored prediction or its
features (constitution Principle VIII, FR-009/FR-010).

**Scale/Scope**: Single team, low volume in v1 (tens to low hundreds of
EstimationRecord/OutcomeRecord rows). Schema must not block a later move to a shared
Postgres instance or to multi-team scale.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Check | Result |
|---|---|---|
| I. Evidence Over Intuition | Model weights/logic documented in research.md as derived from SDD §8-§12, not ad hoc | PASS |
| II. Effort Is Not Duration | `Prediction` keeps `HumanEffort`/`AIEffort`/etc. and `ExpectedDuration` as distinct fields, never merged | PASS |
| III. Effort Is Not Cost | `ExpectedAICost` is a separate field computed independently of effort fields | PASS |
| IV. Business Value Is Not Delivery Effort | No business-value input exists anywhere in `EstimationFeatures` or the engine | PASS |
| V. Human and AI Effort Are Independent | `HumanEffort` and `AIEffort` are computed from disjoint feature subsets and reported separately | PASS |
| VI. Estimation Is a Prediction, Not a Commitment | API/CLI output and docs (quickstart.md) label results "prediction"; no commitment language | PASS |
| VII. Predictions Must Be Measurable | `OutcomeRecord` + on-demand `ErrorReport` exist specifically to close this loop | PASS |
| VIII. The Model Must Be Calibratable (NON-NEGOTIABLE) | Schema is append-only (INSERT-only on records); `modelVersion`/`calibrationVersion` stored per record; no update path is exposed | PASS |
| IX. LLMs Are Not the Estimation Engine | Engine is pure Go arithmetic; no LLM/API call anywhere on the prediction path | PASS |
| X. Consult the Dependency Graph Before Planning | `graphify update .` already run this session; repo is greenfield (single new module, no existing services to conflict with) — confirmed in research.md | PASS |
| Estimation Engine Constraints | Deterministic v1, no ML; DEU never asserted as equivalent to hours in any output copy; engine decoupled from LLM/VCS/language via the repository-interface boundary | PASS |

No violations. Complexity Tracking table is intentionally empty.

## Project Structure

### Documentation (this feature)

```text
specs/001-core-estimation-engine/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (CLI + HTTP JSON contracts)
└── tasks.md             # Phase 2 output (/speckit-tasks — not created here)
```

### Source Code (repository root)

```text
go.mod
go.sum

cmd/
├── estimatorctl/        # CLI entrypoint: estimate, record-outcome, error-report
│   └── main.go
└── estimatord/          # HTTP JSON entrypoint (thin net/http mux over the same use cases)
    └── main.go

internal/
├── estimation/          # Pure domain logic — no I/O
│   ├── features.go      # EstimationFeatures type + validation/clamping (FR-001..FR-003)
│   ├── model.go         # Deterministic weighted model -> Prediction (FR-004..FR-007, FR-016..FR-018)
│   ├── model_test.go
│   └── errorreport.go   # Absolute/relative error + bias computation (FR-014)
│
├── record/               # Application/use-case layer + storage interfaces
│   ├── types.go          # EstimationRecord, OutcomeRecord, ErrorReport (FR-008..FR-013)
│   ├── service.go        # Orchestrates estimation.Model + repositories, enforces
│   │                      # append-only / duplicate-outcome / unknown-id rules
│   ├── service_test.go
│   └── repository.go     # EstimationRepository, OutcomeRepository interfaces
│
└── storage/
    └── sqlite/            # database/sql + modernc.org/sqlite implementation
        ├── sqlite.go
        ├── migrations.go  # schema creation (append-only tables, versioned)
        └── sqlite_test.go

data/                      # default SQLite file location (gitignored)
```

**Structure Decision**: Single project, idiomatic Go layout (`cmd/` + `internal/`)
rather than the template's generic `src/`/`lib/` — this is the natural structure for a
Go module and keeps the domain (`internal/estimation`) fully decoupled from storage
(`internal/storage/sqlite`) and from delivery mechanism (`cmd/*`), which is exactly the
boundary the constitution's "Estimation Engine Constraints" require. Tests are
colocated as `_test.go` files next to the code they test, per Go convention, rather
than under a separate top-level `tests/` directory.

## Complexity Tracking

*No Constitution Check violations — table intentionally left empty.*
