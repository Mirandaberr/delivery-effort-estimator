# Delivery Effort Estimator

An AI-native Delivery Effort Estimation Framework: turns the information
available during specification/planning of a work item into a quantitative
Delivery Effort prediction, with confidence, uncertainty, and risk — see
[`docs/sdd-estimate.md`](docs/sdd-estimate.md) for the full design document.

This repository currently implements its first vertical slice: the
**Core Estimation Engine**. See
[`specs/001-core-estimation-engine/`](specs/001-core-estimation-engine/) for
the full spec, plan, and task breakdown (built with
[GitHub Spec Kit](https://github.com/github/spec-kit)).

## What's implemented

Given a normalized 7-dimension feature vector for a work item, the engine:

- computes a deterministic `Prediction` — Delivery Effort (DEU), Confidence,
  a Prediction Interval (P50/P80), Human/AI/Verification/Integration Effort,
  Risk, Expected Duration, and Expected AI Cost;
- persists it as an immutable, versioned `EstimationRecord`;
- accepts an `OutcomeRecord` once a work item is delivered;
- computes an on-demand `ErrorReport` comparing prediction vs. actual.

No feature extraction, no automatic calibration, and no Specification/Planning
integration yet — those are follow-up features (see spec.md "Assumptions").

## Build & test

```bash
make build   # go build ./...
make test    # go test ./...
make vet     # go vet ./...
make fmt     # gofmt -l . (lists any unformatted files)
```

## Try it

See [`specs/001-core-estimation-engine/quickstart.md`](specs/001-core-estimation-engine/quickstart.md)
for the full walkthrough. Short version:

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

Or run the HTTP server (`go run ./cmd/estimatord`, default `:8080`) — see
[`specs/001-core-estimation-engine/contracts/http.md`](specs/001-core-estimation-engine/contracts/http.md)
for the routes.

By default both binaries store data in `./data/estimator.db` (embedded
SQLite, pure-Go driver, no CGO); override with `ESTIMATOR_DB_PATH`.

## Project layout

```text
cmd/estimatorctl/   CLI entrypoint
cmd/estimatord/     HTTP JSON entrypoint
internal/estimation/  pure deterministic model (no I/O)
internal/record/      orchestration + storage interfaces
internal/storage/sqlite/  SQLite-backed repository implementation
docs/sdd-estimate.md   the original design document
specs/                 GitHub Spec Kit artifacts (spec/plan/tasks per feature)
```
