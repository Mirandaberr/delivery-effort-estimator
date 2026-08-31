# Phase 1 Data Model: Core Estimation Engine

All records are append-only (Constitution Principle VIII — see research.md "Append-only
schema"). Field types are given as Go types; the SQLite column types follow directly.

## EstimationFeatures

Input to the engine. Not persisted on its own — persisted embedded inside an
`EstimationRecord`.

| Field | Type | Validation |
|---|---|---|
| `context_complexity` | `float64` | clamp to `[0,1]`; required (FR-002) |
| `domain_complexity` | `float64` | clamp to `[0,1]`; required |
| `integration_complexity` | `float64` | clamp to `[0,1]`; required |
| `verification_complexity` | `float64` | clamp to `[0,1]`; required |
| `human_decision_load` | `float64` | clamp to `[0,1]`; required |
| `ai_execution_complexity` | `float64` | clamp to `[0,1]`; required |
| `uncertainty` | `float64` | clamp to `[0,1]`; required |

Any clamp applied is recorded as a human-readable string in `EstimationRecord.assumptions`
(e.g. `"clamped ai_execution_complexity from 1.40 to 1.00"`) — FR-003.

## Prediction

Computed output, embedded inside `EstimationRecord` (never persisted standalone).

| Field | Type | Notes |
|---|---|---|
| `delivery_effort_deu` | `float64` | see research.md formula |
| `confidence` | `float64` | `[0.05, 0.70]` for `v1-linear` |
| `prediction_interval_p50` | `float64` | == `delivery_effort_deu` |
| `prediction_interval_p80` | `float64` | always `> p50` |
| `human_effort` | `float64` | |
| `ai_effort` | `float64` | |
| `verification_effort` | `float64` | |
| `integration_effort` | `float64` | |
| `risk_score` | `float64` | `[0,1]` |
| `risk_label` | `string` | `"low" \| "medium" \| "high"` |
| `expected_duration_days` | `float64` | |
| `expected_ai_cost_usd` | `float64` | `0` when `ExpectedAICostEstimated == false` |
| `expected_ai_cost_estimated` | `bool` | `false` only when `ai_effort == 0` exactly |

## EstimationRecord

Persisted table `estimation_records`, one row per estimation run (never updated).

| Field | Type | Notes |
|---|---|---|
| `id` | `string` (UUID) | primary key |
| `work_item_id` | `string` | caller-supplied identifier; not unique — re-estimation creates a new row (FR-010) |
| `timestamp` | `string` (RFC3339) | set server-side at creation |
| `specification_version` | `string`, nullable | out of scope this feature; always `NULL` for now |
| `planning_version` | `string`, nullable | out of scope this feature; always `NULL` for now |
| `repository_revision` | `string`, nullable | out of scope this feature; always `NULL` for now |
| `features` | JSON blob (`EstimationFeatures`) | as clamped/validated, not the raw input |
| `prediction` | JSON blob (`Prediction`) | full computed output |
| `confidence` | `float64` | denormalized copy of `prediction.confidence` for querying |
| `prediction_interval_p50` / `_p80` | `float64` | denormalized for querying |
| `risks` | JSON array of `string` | human-readable risk notes derived from feature values above defined thresholds |
| `assumptions` | JSON array of `string` | includes clamp adjustments (FR-003) and any other defaults applied |
| `model_version` | `string` | `"v1-linear"` |
| `calibration_version` | `string` | `"uncalibrated"` until a Calibration Engine (out of scope) exists |

**Invariant**: no code path may `UPDATE` or `DELETE` a row in this table (FR-009).

## OutcomeRecord

Persisted table `outcome_records`, at most one row per `estimation_id`
(`UNIQUE(estimation_id)` — FR-013).

| Field | Type | Notes |
|---|---|---|
| `id` | `string` (UUID) | primary key |
| `estimation_id` | `string` | foreign key → `estimation_records.id`; must exist (FR-012) |
| `actual_human_effort` | `float64` | caller-supplied unit (documented as hours by convention) |
| `actual_ai_usage` | `float64` | caller-supplied unit (e.g. tokens or agent-minutes; framework does not enforce which) |
| `actual_ai_cost_usd` | `float64` | |
| `actual_lead_time_days` | `float64` | |
| `actual_verification_effort` | `float64` | |
| `actual_integration_effort` | `float64` | |
| `rework` | `float64` | caller-defined count/measure of rework iterations |
| `incidents` | `int` | count |
| `completion_timestamp` | `string` (RFC3339) | caller-supplied |

## ErrorReport

Computed on demand from an `EstimationRecord` + its `OutcomeRecord` — **not persisted**
(FR-015: computing it must not mutate anything). Returned as a response object only.

| Field | Type | Notes |
|---|---|---|
| `estimation_id` | `string` | |
| `dimensions` | array of `DimensionError` | one entry per dimension present in both records |

### DimensionError

| Field | Type | Notes |
|---|---|---|
| `name` | `string` | e.g. `"human_effort"`, `"ai_cost_usd"`, `"lead_time_days"` |
| `predicted` | `float64` | |
| `actual` | `float64` | |
| `absolute_error` | `float64` | `\|actual - predicted\|` |
| `relative_error` | `float64`, nullable | `null` when `predicted == 0` |
| `bias` | `float64` | `actual - predicted` (signed) |

Dimensions compared in v1: `human_effort`, `ai_effort` (vs `actual_ai_usage`,
raw numeric — see research.md), `verification_effort`, `integration_effort`,
`ai_cost_usd`, `lead_time_days`.

## State Transitions

Both `EstimationRecord` and `OutcomeRecord` are immutable once created — there is no
state machine. The only two operations mutate storage: `INSERT estimation_records` and
`INSERT outcome_records` (rejected if the referenced estimation doesn't exist, or
already has an outcome). `ErrorReport` is pure computation, no storage effect.
