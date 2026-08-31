# Quickstart: Core Estimation Engine

Validates User Stories 1–3 end-to-end (spec.md) once implementation exists.

## Prerequisites

- Go 1.27+ installed (`go version`).
- From repo root: `go build ./...` succeeds.

## Setup

```bash
go run ./cmd/estimatorctl --help   # confirms the binary builds and lists commands
```

The default SQLite file is created automatically at `./data/estimator.db` on first
write (directory created if missing).

## Scenario 1 — Get a structured prediction (User Story 1)

```bash
cat > /tmp/features.json <<'EOF'
{
  "context_complexity": 0.72,
  "domain_complexity": 0.41,
  "integration_complexity": 0.83,
  "verification_complexity": 0.67,
  "human_decision_load": 0.54,
  "ai_execution_complexity": 0.38,
  "uncertainty": 0.61
}
EOF

go run ./cmd/estimatorctl estimate --work-item WI-1 --features /tmp/features.json
```

**Expected**: JSON `EstimationRecord` printed with all `Prediction` fields populated
(`delivery_effort_deu`, `confidence`, both interval bounds, the four effort
sub-scores, `risk_score`/`risk_label`, `expected_duration_days`,
`expected_ai_cost_usd`). Running the exact same command again produces byte-identical
`prediction` content (Acceptance Scenario 1.2 — determinism), differing only in `id`
and `timestamp`.

## Scenario 2 — Reproducibility across model versions (User Story 2)

```bash
go run ./cmd/estimatorctl estimate --work-item WI-1 --features /tmp/features.json \
  | tee /tmp/first-estimate.json
```

Note the printed `id`. After any future model change ships as `v2-*`, re-fetching this
same `id` (`GET /estimations/{id}` or a future `estimatorctl get` command) must return
the identical `prediction` and `model_version: "v1-linear"` recorded at creation time —
this is the manual check for Acceptance Scenario 2.2 until an automated regression
test covers it in `internal/record/service_test.go`.

## Scenario 3 — Compare prediction against a real outcome (User Story 3)

```bash
ESTIMATION_ID=$(python3 -c "import json;print(json.load(open('/tmp/first-estimate.json'))['id'])")

cat > /tmp/outcome.json <<'EOF'
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
EOF

go run ./cmd/estimatorctl record-outcome --estimation-id "$ESTIMATION_ID" --outcome /tmp/outcome.json
go run ./cmd/estimatorctl error-report --estimation-id "$ESTIMATION_ID"
```

**Expected**: `error-report` prints an `ErrorReport` JSON with one `DimensionError`
entry per comparable dimension (see data-model.md), each carrying
`absolute_error`/`relative_error`/`bias`.

**Edge case check**: re-running `record-outcome` for the same `$ESTIMATION_ID` a
second time must fail with `{"error": "outcome_already_recorded", ...}` and a non-zero
exit code (FR-013).

## Automated coverage

- `go test ./internal/estimation/...` — table-driven tests for the formulas in
  research.md, including the clamping and zero-cost edge cases from spec.md.
- `go test ./internal/record/...` — service-level tests for duplicate-outcome and
  unknown-id rejection (FR-012/FR-013), run against a temp-file SQLite DB.
- `go test ./...` from repo root must pass before a feature is considered done.
