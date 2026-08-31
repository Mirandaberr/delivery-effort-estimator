# Data Model: Feature Extraction Plugin

This feature introduces one new entity, `DerivedFeatureVector`. It produces
`EstimationFeatures` as input to, and reuses `Prediction`/`EstimationRecord`
as output from, the existing Estimation Engine (specs/001-core-estimation-engine/data-model.md)
without modification.

## DerivedFeatureVector

Written to `specs/<id>-*/estimation/<UTC-timestamp>/derived-features.json` by
the `estimate-feature` skill (FR-002, FR-003, FR-004).

| Field | Type | Notes |
|---|---|---|
| `context_complexity`, `domain_complexity`, `integration_complexity`, `verification_complexity`, `human_decision_load`, `ai_execution_complexity`, `uncertainty` | float, each `[0,1]` | The 7 `EstimationFeatures` dimensions (specs/001), at the **top level** of the file, not nested — see below |
| `feature_directory` | string | Repo-relative path, e.g. `specs/002-feature-extraction-plugin` |
| `derived_at` | string (RFC3339, UTC) | When the derivation ran; also the containing directory's timestamp |
| `source_files` | array of strings | Which of `spec.md`/`plan.md`/`tasks.md` were actually present and read |
| `justifications` | object | One key per dimension above, each a 1-2 sentence string explaining the score, per the research.md rubric |
| `assumptions` | array of strings | Notes for any dimension scored conservatively due to missing/thin source content (mirrors the engine's own `Clamp()` assumption notes from specs/001) |

**Why the 7 dimensions are flat, not nested under a `features` key**:
`estimatorctl estimate --features <path>` (specs/001/contracts/cli.md) is
passed this exact file. Its `ParseFeatures` (internal/estimation/features.go)
reads the 7 dimension names directly off the top-level JSON object and
ignores any other top-level keys it doesn't recognize — it does **not**
look inside a nested `features` object. So `derived-features.json` doubles
as the CLI's `--features` input verbatim: `feature_directory`,
`justifications`, etc. are the "harmless extra keys" `ParseFeatures` ignores
(contracts/skill.md step 5).

**Validation rules**:
- Every one of the 7 dimension keys MUST have a matching key in
  `justifications` — a score with no justification is invalid output
  (FR-003 is not optional).
- The 7 dimension values MUST already be in `[0,1]`; the skill is
  responsible for producing valid input, but `estimatorctl` still clamps
  defensively (specs/001 `Clamp()`) since it is a separate process boundary.

**State transitions**: None — immutable once written, like `EstimationRecord`.
A new derivation for the same feature directory is a new directory under
`estimation/`, never an edit to a prior one (FR-007).

## Reused entities (specs/001-core-estimation-engine, unchanged)

- **EstimationFeatures** — the 7 top-level dimension fields above, once
  validated, are read directly as this type by `estimatorctl estimate`.
- **Prediction / EstimationRecord** — the unmodified output of
  `estimatorctl estimate`, saved verbatim as
  `specs/<id>-*/estimation/<UTC-timestamp>/estimation-record.json`.

## PluginEngineBuild (operational state, not a domain entity)

Not part of the estimation domain model, but recorded here since it's
user-visible file state: `plugin/bin/.built-checksum` holds the `cksum`
value of `plugin/engine/.sync-checksum` at the time `plugin/bin/estimatorctl`
was last built (research.md "Engine freshness check"). It exists purely to
make `ensure-engine.sh` idempotent and is not versioned in git (build
output, like `plugin/bin/estimatorctl` itself).
