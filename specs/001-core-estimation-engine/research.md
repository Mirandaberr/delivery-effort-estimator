# Phase 0 Research: Core Estimation Engine

No `[NEEDS CLARIFICATION]` markers were left in `spec.md`, but the spec deliberately
leaves the *quantitative* model, storage schema strategy, and comparison semantics
unspecified (the SDD explicitly says the v1 model can be "deterministic/simple" and
that DEU's meaning is "learned empirically"). Those are technical decisions this plan
must make explicit before Phase 1 design. Each is recorded below.

## Decision: Language & runtime — Go 1.27, standard library only + one pure-Go driver

**Rationale**: User explicitly asked for "lo más liviano y eficiente posible" over
matching the existing Java/Spring Boot stack used by sibling projects. A compiled Go
binary starts in single-digit milliseconds, needs no JVM/interpreter, and the domain
logic here is pure arithmetic — it does not benefit from a web framework or ORM.
`modernc.org/sqlite` is chosen over `mattn/go-sqlite3` specifically because it is
CGO-free, keeping the binary a true single static artifact (cross-compilable without a
C toolchain).

**Alternatives considered**:
- Java/Spring Boot — rejected: consistent with sibling projects, but explicitly heavier
  than what the user asked for (JVM startup, DI container, larger footprint).
- Python (stdlib only) — viable and simpler to write, but a compiled Go binary is
  lighter at runtime (no interpreter, lower memory) and the CLI/HTTP surface needed
  here is small enough that Go's verbosity cost is negligible.
- Rust — lightest possible runtime, rejected only because it adds meaningfully more
  implementation time for the same deterministic-arithmetic payoff; revisit if a future
  feature needs Rust-level performance guarantees.

## Decision: Storage — embedded SQLite behind a repository interface

**Rationale**: User confirmed "H2/SQLite embebido" — zero external infrastructure to
start, and the constitution ("Estimation Engine Constraints") already requires the
engine be decoupled from any specific storage. Modeling `EstimationRepository` /
`OutcomeRepository` as interfaces in `internal/record` means a future PostgreSQL
implementation is a second adapter, not a rewrite.

**Alternatives considered**:
- Plain JSON files on disk — rejected: no query capability for the future
  Calibration Engine (out of scope here, but the schema must not block it).
- In-memory only — rejected per spec FR-008/FR-009 (records must persist and remain
  reproducible across process restarts); user also explicitly picked the embedded-DB
  option over in-memory.

## Decision: Deterministic v1 model — fixed linear weights over the 7 SDD features

**Rationale**: SDD §11 permits (requires, for v1) a deterministic/simple model, and
Constitution Principle IX forbids the estimation math itself from calling an LLM. A
transparent, fixed-weight linear combination is the simplest model that is (a)
reproducible, (b) explainable per-dimension (Constitution Principle I — evidence, not
a black box), and (c) trivially versionable (`ModelVersion = "v1-linear"`) so later
calibration can replace the weights wholesale as `v2` without touching v1's stored
records.

Sub-effort scores (each computed in `[0,1]` before DEU scaling):

| Output | Formula |
|---|---|
| `HumanEffort` | `0.45·human_decision_load + 0.30·context_complexity + 0.15·domain_complexity + 0.10·uncertainty` |
| `AIEffort` | `0.50·ai_execution_complexity + 0.25·context_complexity + 0.15·domain_complexity + 0.10·uncertainty` |
| `VerificationEffort` | `0.55·verification_complexity + 0.20·integration_complexity + 0.15·uncertainty + 0.10·domain_complexity` |
| `IntegrationEffort` | `0.70·integration_complexity + 0.15·domain_complexity + 0.15·uncertainty` |

`DeliveryEffort` (DEU, scale 0–10): `10 · (0.35·HumanEffort + 0.20·AIEffort + 0.25·VerificationEffort + 0.20·IntegrationEffort)`.

`Risk` (0–1, plus a Low/Medium/High label at 0.34/0.67 thresholds): equal to
`uncertainty` in v1 — the simplest defensible mapping until calibration data suggests
a richer risk formula.

`Confidence`: `clamp(1 - uncertainty, 0.05, 0.70)`. The upper cap of `0.70` for model
`v1-linear` is deliberate: with zero calibration history, the model must not claim
high confidence even when input uncertainty is reported as low (this directly
implements the Edge Case requirement that the system "must still return an interval
... rather than omit uncertainty" and Constitution Principle VII).

`PredictionInterval` (P50/P80): `P50 = DeliveryEffort`;
`spread = 0.25 + 0.75·(1 - Confidence)`; `P80 = DeliveryEffort · (1 + spread)`. The
`0.25` floor guarantees P80 is never equal to P50 — the interval always expresses some
uncertainty, satisfying the same edge case even in the best-confidence scenario.

`ExpectedDuration` (calendar days): `0.2 · DeliveryEffort`. Chosen to land close to the
SDD's own worked example (§12: "8 DEU → 1.5 days"; this formula gives 1.6 days),
while being explicitly labeled a *prediction*, not a defined conversion factor
(Constitution: DEU section) — recalibration will change this constant, not the
DeliveryEffort computation itself.

`ExpectedAICost` (USD): `0.3 · AIEffort · DeliveryEffort`, with `ExpectedAICostEstimated
= false` and cost forced to `0` whenever `AIEffort == 0` exactly — this satisfies the
Edge Case requirement to distinguish "zero predicted cost" from "cost could not be
estimated" instead of returning an indistinguishable `0`.

**Alternatives considered**:
- Equal weighting across all 7 features for every output — rejected: fails Principle
  I (not evidence-based; SDD explicitly separates Human/AI/Verification/Integration as
  distinct concerns with different drivers, e.g. `integration_complexity` should
  dominate `IntegrationEffort`, not be diluted equally with unrelated features).
- Leaving Confidence uncapped (`1 - uncertainty` directly) — rejected: an SDD-input of
  `uncertainty = 0` would then claim 100% confidence from a fixed-weight linear model
  with zero historical validation, directly contradicting Principle VI (a prediction,
  not a truth).

## Decision: Cross-dimension error comparison uses raw numeric deltas, not unit conversion

**Rationale**: `HumanEffort`/`AIEffort`/`VerificationEffort`/`IntegrationEffort` are
reported on the framework's own 0–10-ish scale, while `OutcomeRecord`'s
`actualHumanEffort`/`actualAIUsage`/etc. are real-world units (hours, tokens). SDD §12
is explicit that "DEU's meaning will be learned empirically" — i.e., establishing the
DEU-to-real-unit mapping is the *Calibration Engine's* job, not this feature's.
Therefore `ErrorReport` computes `absoluteError = |actual - predicted|`,
`relativeError = absoluteError / predicted` (flagged undefined when `predicted == 0`),
and `bias = actual - predicted` per dimension as raw numeric deltas. `ExpectedDuration`
vs `actualLeadTime` (days) and `ExpectedAICost` vs `actualAICost` (USD) are the two
dimensions that are already unit-comparable in v1.

**Alternatives considered**:
- Deferring `ErrorReport` entirely until a unit-conversion model exists — rejected:
  spec FR-014/User Story 3 require it now; the raw-delta numbers are exactly the input
  the future Calibration Engine needs to learn that conversion (SDD §17).

## Decision: Append-only schema enforces Constitution Principle VIII in the storage layer, not just in application logic

**Rationale**: Rather than trusting every caller of `EstimationRepository` to "not
update" a record, the SQLite schema exposes only `INSERT` operations from
`internal/record/service.go`; there is no `Update*` method on the repository
interfaces at all. Re-estimating a work item creates a new row with a new `id` and the
same `workItemId` (FR-010). A duplicate `OutcomeRecord` for an `estimationId` that
already has one is rejected by a uniqueness constraint on `outcome_records.estimation_id`
(FR-013), enforced at the DB level so it holds even if application logic has a bug.

## Graphify dependency-graph check (Constitution Principle X)

`graphify update .` was run at repo bootstrap (`graphify-out/graph.json`,
`GRAPH_REPORT.md`): 221 nodes / 248 edges, entirely spec-kit scaffolding and templates
— this repository has no pre-existing application code or services for this feature to
conflict with. No `graphify explain`/`graphify path` queries were needed for this
plan; this check will be re-run before `/speckit-tasks` for any *future* feature once
`internal/estimation` and `internal/record` exist, to confirm new work doesn't
silently duplicate or bypass them.
