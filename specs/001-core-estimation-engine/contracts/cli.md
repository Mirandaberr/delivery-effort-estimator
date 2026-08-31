# Contract: `estimatorctl` CLI

Text in/out: flags in, JSON on stdout, errors on stderr with non-zero exit code.

## `estimatorctl estimate`

```text
estimatorctl estimate --work-item <id> --features <path-to-json-or-'-'-for-stdin>
```

Input JSON (`EstimationFeatures`, all 7 fields required — FR-001/FR-002):

```json
{
  "context_complexity": 0.72,
  "domain_complexity": 0.41,
  "integration_complexity": 0.83,
  "verification_complexity": 0.67,
  "human_decision_load": 0.54,
  "ai_execution_complexity": 0.38,
  "uncertainty": 0.61
}
```

Success (exit 0), stdout = the full `EstimationRecord` as JSON (see data-model.md).

Failure (exit 1), stderr = `{"error": "missing_feature", "field": "uncertainty"}` or
similar structured error — one of `missing_feature`, `invalid_json`.

## `estimatorctl record-outcome`

```text
estimatorctl record-outcome --estimation-id <id> --outcome <path-to-json-or-'-'>
```

Input JSON (`OutcomeRecord` fields, all required except `rework`/`incidents` default 0
— FR-011):

```json
{
  "actual_human_effort": 6.0,
  "actual_ai_usage": 42000,
  "actual_ai_cost_usd": 2.10,
  "actual_lead_time_days": 1.5,
  "actual_verification_effort": 3.0,
  "actual_integration_effort": 2.0,
  "rework": 1,
  "incidents": 0,
  "completion_timestamp": "2026-09-02T10:00:00Z"
}
```

Success (exit 0): stdout = the persisted `OutcomeRecord` JSON.

Failure (exit 1): `{"error": "unknown_estimation_id", "estimation_id": "..."}` (FR-012)
or `{"error": "outcome_already_recorded", "estimation_id": "..."}` (FR-013).

## `estimatorctl error-report`

```text
estimatorctl error-report --estimation-id <id>
```

Success (exit 0): stdout = `ErrorReport` JSON (see data-model.md).

Failure (exit 1): `{"error": "unknown_estimation_id", ...}` or
`{"error": "no_outcome_recorded", "estimation_id": "..."}` (no outcome yet to compare
against).
