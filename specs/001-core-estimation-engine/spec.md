# Feature Specification: Core Estimation Engine

**Feature Branch**: `001-core-estimation-engine`

**Created**: 2026-08-31

**Status**: Draft

**Input**: User description: "Primer slice del AI-Native Delivery Effort Estimation
Framework descrito en docs/sdd-estimate.md: un Estimation Engine determinístico v1
que recibe un vector normalizado de features, produce una predicción estructurada
(Delivery Effort, Confidence, Prediction Interval, Human/AI/Verification/Integration
Effort, Risk, Expected Duration, Expected AI Cost), la persiste como Estimation
Record reproducible, y permite comparar esa predicción contra un Outcome Record real
calculando error absoluto/relativo. Feature extraction automática, calibración
automática y las capas de Specification/Planning quedan fuera de esta primera
feature."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Get a structured prediction from a feature set (Priority: P1)

An engineer (or a calling tool such as Claude Code) has already derived a normalized
feature vector for a work item — however it was derived — and needs a quantitative,
explainable prediction of delivery effort before work starts.

**Why this priority**: This is the smallest possible vertical slice that delivers the
core value promised by the SDD ("transformar información disponible en una predicción
cuantitativa de Delivery Effort"). Without this, nothing else in the framework
(feature extraction, planning integration, calibration) has anything to call.

**Independent Test**: Submit a complete feature vector for a fictitious work item and
verify a fully-populated structured prediction is returned, with no specification
tooling, repository analysis, or planning layer implemented.

**Acceptance Scenarios**:

1. **Given** a normalized feature vector with all required dimensions populated,
   **When** an estimation is requested for a work item, **Then** the system returns a
   prediction containing Delivery Effort (DEU), Confidence, a Prediction Interval,
   Human Effort, AI Effort, Verification Effort, Integration Effort, Risk, Expected
   Duration, and Expected AI Cost.
2. **Given** the same feature vector and the same model version, **When** the
   estimation is requested twice, **Then** both predictions are identical (deterministic,
   no reliance on live LLM inference).
3. **Given** a feature vector with a value outside the [0,1] domain, **When** an
   estimation is requested, **Then** the system clamps the value to the nearest bound,
   completes the estimation, and records the adjustment as an assumption on the result.

---

### User Story 2 - Every prediction is a reproducible, versioned record (Priority: P2)

A team lead needs to trust that a prediction made months ago can still be explained
and reproduced today, even if the model has since been recalibrated.

**Why this priority**: Reproducibility and auditability are constitutional
requirements (Principle VIII) and a prerequisite for any later calibration work —
without a durable record, Prediction vs. Actual comparison (User Story 3) has nothing
to compare against.

**Independent Test**: Store a prediction, note its `modelVersion`, evolve/replace the
model, and confirm the original record still reports the same prediction and can be
explained by re-running only against the model version it was created with.

**Acceptance Scenarios**:

1. **Given** a completed estimation, **When** it is produced, **Then** an Estimation
   Record is persisted with a unique id, the work item id, a timestamp, the input
   features, the full prediction, and the model/calibration version that produced it.
2. **Given** an existing Estimation Record, **When** the underlying model is later
   replaced by a new version, **Then** the stored record's prediction and
   `modelVersion` remain unchanged.
3. **Given** a work item that is re-estimated, **When** the second estimation
   completes, **Then** a new Estimation Record is created rather than overwriting the
   first one.

---

### User Story 3 - Compare a prediction against what actually happened (Priority: P3)

After a work item is delivered, a team lead records what actually happened and wants
to see, immediately, how far off the original prediction was.

**Why this priority**: This closes the smallest possible feedback loop the SDD
requires ("Prediction vs Actual") and produces the first evidence the framework will
later need for calibration — but it depends on User Stories 1 and 2 already existing,
so it is lower priority than either.

**Independent Test**: Given an existing Estimation Record and a submitted Outcome
Record referencing it, request an error report and confirm it returns absolute and
relative error per comparable dimension without invoking any automatic model
adjustment.

**Acceptance Scenarios**:

1. **Given** a stored Estimation Record and a completed work item, **When** an Outcome
   Record is submitted with actual measurements, **Then** the system accepts and
   persists it linked to that Estimation Record.
2. **Given** a matched Estimation Record and Outcome Record, **When** an error report
   is requested, **Then** the system returns absolute error, relative error, and bias
   per comparable dimension (e.g., Human Effort, AI Cost, Lead Time).
3. **Given** an Outcome Record that references a non-existent Estimation Record,
   **When** it is submitted, **Then** the system rejects it with a clear error instead
   of silently creating an orphaned record.

### Edge Cases

- What happens when a required feature dimension is missing entirely from the request
  (as opposed to out-of-range)? System must reject the request rather than guess a
  default value for a missing dimension.
- How does the system behave when it has no historical data yet to inform a Prediction
  Interval? It must still return an interval (conservatively wide) rather than omit
  uncertainty entirely, per constitutional Principle VII ("Predictions must be
  measurable").
- What happens if an Outcome Record is submitted twice for the same Estimation Record?
  The second submission must be rejected or explicitly versioned — it must not silently
  overwrite the first (consistent with Principle VIII's "never overwrite history").
- What happens when Expected AI Cost cannot be computed (e.g., no cost-relevant
  features supplied)? The system must return a zero/undefined cost with an explicit
  flag rather than a misleading zero indistinguishable from "no cost predicted."

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST accept, per estimation request, a normalized ([0,1]) feature
  vector for a work item covering at least: context_complexity, domain_complexity,
  integration_complexity, verification_complexity, human_decision_load,
  ai_execution_complexity, and uncertainty.
- **FR-002**: System MUST reject an estimation request that omits a required feature
  dimension, identifying which dimension is missing.
- **FR-003**: System MUST clamp any feature value outside [0,1] to the nearest bound
  and record that adjustment as an assumption on the resulting Estimation Record,
  rather than rejecting the request or producing an undefined result.
- **FR-004**: System MUST produce, for every accepted estimation request, a prediction
  containing: Delivery Effort (in Delivery Effort Units, DEU), Confidence, a Prediction
  Interval (at minimum two probability points, e.g. P50/P80), Human Effort, AI Effort,
  Verification Effort, Integration Effort, Risk, Expected Duration, and Expected AI
  Cost.
- **FR-005**: The engine's computation MUST be deterministic for a given model version:
  identical inputs and model version MUST always produce an identical prediction, with
  no dependency on live LLM inference or non-reproducible randomness at prediction
  time.
- **FR-006**: System MUST NOT require any specific LLM provider, programming language,
  or version-control system to produce a prediction — the engine's only inputs are the
  feature vector and stored historical/model context.
- **FR-007**: System MUST NOT require Machine Learning infrastructure in this version;
  a deterministic/rule-based computation is sufficient and required for v1.
- **FR-008**: System MUST persist every produced prediction as an Estimation Record
  containing, at minimum: a unique id, work item id, timestamp, the input features, the
  full prediction, confidence, prediction interval, risks, assumptions, model version,
  and calibration version.
- **FR-009**: System MUST version its model explicitly and MUST NOT mutate or
  overwrite a previously stored Estimation Record when the model or calibration
  version changes.
- **FR-010**: System MUST create a new Estimation Record — never update an existing
  one — when a work item is re-estimated.
- **FR-011**: System MUST allow submitting an Outcome Record that references an
  existing Estimation Record, containing at minimum: actual human effort, actual AI
  usage, actual AI cost, actual lead time, actual verification effort, actual
  integration effort, rework, incidents, and a completion timestamp.
- **FR-012**: System MUST reject an Outcome Record that references an Estimation
  Record id that does not exist.
- **FR-013**: System MUST reject a second Outcome Record submission for an Estimation
  Record that already has one, rather than silently overwriting the first.
- **FR-014**: System MUST, on request, compute and return absolute error, relative
  error, and bias between a matched Estimation Record and its Outcome Record, for each
  dimension present in both (e.g., Human Effort, AI Cost, Lead Time, Verification
  Effort, Integration Effort).
- **FR-015**: System MUST NOT perform any automatic model recalibration as part of
  computing or storing an error report — error computation and model calibration are
  separate concerns; this feature only covers the former.
- **FR-016**: System MUST treat the Delivery Effort Unit (DEU) as an internal relative
  unit; it MUST NOT expose a fixed DEU-to-hours or DEU-to-Story-Points conversion as
  authoritative — any such mapping surfaced to a caller MUST be labeled as a
  prediction, not a definition.
- **FR-017**: System MUST report Human Effort, AI Effort, Verification Effort, and
  Integration Effort as independently readable values within a single prediction, not
  collapsed into one combined number.
- **FR-018**: System MUST report Expected Duration separately from Human Effort and
  from Expected AI Cost (effort, duration, and cost remain distinct dimensions).

### Key Entities

- **EstimationFeatures**: The normalized [0,1] input vector describing a work item at
  estimation time (context_complexity, domain_complexity, integration_complexity,
  verification_complexity, human_decision_load, ai_execution_complexity,
  uncertainty).
- **Prediction**: The structured output of one estimation run — deliveryEffort (DEU),
  confidence, predictionInterval, humanEffort, aiEffort, verificationEffort,
  integrationEffort, risk, expectedDuration, expectedAICost.
- **EstimationRecord**: The persisted, immutable result of one estimation run — id,
  workItemId, timestamp, features, prediction, confidence, predictionInterval, risks,
  assumptions, modelVersion, calibrationVersion. (specificationVersion,
  planningVersion, repositoryRevision fields exist per the SDD data model but are
  optional/unset until the Specification and Planning layers are implemented in a
  later feature.)
- **OutcomeRecord**: Real measurements collected after a work item is delivered —
  estimationId, actualHumanEffort, actualAIUsage, actualAICost, actualLeadTime,
  actualVerificationEffort, actualIntegrationEffort, rework, incidents,
  completionTimestamp.
- **ErrorReport**: The computed, on-demand comparison between one EstimationRecord and
  its OutcomeRecord — absolute error, relative error, and bias per comparable
  dimension.
- **ModelVersion**: Identifies which deterministic model logic produced a given
  prediction, independent of CalibrationVersion, so a historical prediction stays
  reproducible after later model changes.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Given a complete, valid feature vector, a fully-populated structured
  prediction (all required fields present) is returned in under 1 second, with no
  manual steps.
- **SC-002**: 100% of produced predictions are persisted as retrievable Estimation
  Records that reproduce an identical prediction when replayed under the same recorded
  model version.
- **SC-003**: Given an existing Estimation Record, recording an Outcome Record and
  obtaining its error report requires no more than two operations (submit outcome,
  request report) and no manual/spreadsheet computation.
- **SC-004**: Repeating an identical estimation request against an unchanged model
  version yields 0% variance across repeated runs.
- **SC-005**: Introducing a new model version never alters the prediction or fields of
  any previously stored Estimation Record (0% historical record mutation).
- **SC-006**: 100% of rejected requests (missing feature dimension, duplicate outcome,
  unknown estimation id) return a specific, actionable reason rather than a generic
  failure.

## Assumptions

- Feature extraction — deriving `EstimationFeatures` from a repository, specification,
  planning document, or LLM/engineer conversation (SDD §10) — is out of scope for this
  feature. This feature starts from an already-produced feature vector; extraction is
  a follow-up feature that will call into this engine.
- Automatic Calibration — adjusting model parameters based on accumulated error (SDD
  §17–19) — is out of scope for this feature. This feature only computes and exposes
  error metrics on demand; using them to change the model is a follow-up feature.
- The Specification Layer (SDD §6) and Planning Layer (SDD §7) are out of scope for
  this feature; `specificationVersion`, `planningVersion`, and `repositoryRevision` on
  the Estimation Record are accepted as optional/unset for now.
- Callers of this feature are engineers, team leads, and internal tooling (e.g. Claude
  Code, CI) integrating programmatically; no end-user graphical interface is in scope
  for v1.
- Expected AI Cost is expressed in USD unless a caller-supplied context specifies
  otherwise.
- Expected Duration is expressed in calendar time (hours/days), kept independent from
  Human Effort per constitutional Principle II (Effort ≠ Duration).
- "Risk" in the prediction output is a qualitative/quantitative indicator derived from
  the Uncertainty categories in SDD §8, not a separate detailed risk-register feature.
