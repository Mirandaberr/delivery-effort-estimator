<!--
Sync Impact Report
- Version change: none → 1.0.0 (initial ratification)
- Principles established: I–X (I–IX derived verbatim from docs/sdd-estimate.md §3
  "Principios"; X added by this command per sdd-start skill convention)
- Added sections: Core Principles (10), Estimation Engine Constraints,
  Specification & Development Workflow, Governance
- Removed sections: none (initial document)
- Deferred/TODO items: none
-->

# Delivery Effort Estimator Constitution

## Core Principles

### I. Evidence Over Intuition
Every estimation MUST be grounded in available evidence: specification, architecture,
repository state, historical outcomes, dependencies, prior results, and information
supplied directly by the engineer. Estimates MUST NOT be produced from unexamined
intuition when evidence sources are available.

### II. Effort Is Not Duration
Required effort and elapsed time are distinct metrics and MUST be modeled and reported
separately. A prediction MUST NOT collapse "how much work" and "how long it takes" into
a single number.

### III. Effort Is Not Cost
Economic cost MUST be calculated independently of effort. AI cost (tokens, model,
infrastructure) is a separate economic dimension from Human Effort, AI Effort,
Verification Effort, and Integration Effort, and MUST NOT be conflated with them.

### IV. Business Value Is Not Delivery Effort
The business value of an initiative MUST NOT artificially alter its estimated effort.
Business Value MAY be used downstream for prioritization and ROI, but it is not an
input to the Estimation Engine's effort computation.

### V. Human and AI Effort Are Independent Dimensions
Human participation and AI agent participation MUST be measurable independently of
each other. A low-human-effort task MAY carry high AI effort, and vice versa; the
model MUST be able to represent both simultaneously.

### VI. Estimation Is a Prediction, Not a Commitment
An estimate represents a prediction based on information available at a specific
point in time. It MUST NOT be treated, stored, or presented as a truth, guarantee, or
commitment.

### VII. Predictions Must Be Measurable
Every prediction MUST be comparable against a real outcome after the fact. A
prediction that cannot later be matched to an Outcome Record MUST NOT be considered
complete.

### VIII. The Model Must Be Calibratable (NON-NEGOTIABLE)
The system MUST allow its parameters and methodology to be revised based on historical
evidence (Prediction vs. Actual). Calibration MUST NOT react to a single deviation; it
MUST require statistically sufficient evidence before adjusting the model, and MUST
produce a new, distinct model version rather than mutating a prior one (see §Estimation
Engine Constraints).

### IX. LLMs Are Not the Estimation Engine
LLMs and AI agents MAY analyze context and produce features (context complexity,
domain complexity, uncertainty, etc.), but the quantitative estimation engine MUST
remain decoupled from any specific LLM provider. No estimation output may depend on
non-reproducible LLM inference at prediction time.

### X. Consult the Dependency Graph Before Planning
Before planning a change to this codebase, the graphify dependency graph MUST be
consulted (`graphify explain <node>`, `graphify path <A> <B>`) to identify affected
modules/services, and its output MUST be verified against the actual files it points
to rather than trusted as ground truth on its own.

## Estimation Engine Constraints

- The Estimation Engine MUST remain independent of Claude Code, GitHub Spec Kit, any
  specific LLM provider, programming language, and version control system; its only
  contract is a normalized set of features plus historical context in, and a
  structured prediction (Delivery Effort, Confidence, Prediction Interval, Human
  Effort, AI Effort, Verification Effort, Integration Effort, Risk, Expected
  Duration, Expected AI Cost) out.
- The first version of the engine MUST use a deterministic/simple model. Machine
  Learning MUST NOT be introduced until sufficient historical observation data exists
  to justify it (see docs/sdd-estimate.md §11, §22).
- The Delivery Effort Unit (DEU) is a relative, framework-owned unit. It MUST NOT be
  presented as equivalent to hours or Story Points; any mapping to hours, lead time,
  or cost is itself a prediction, never a definition.
- The model MUST be versioned. Historical estimations MUST NEVER be overwritten or
  silently recomputed: an estimation produced under Model vN MUST remain reproducible
  under Model vN after later calibration produces vN+1.
- Every EstimationRecord MUST persist its modelVersion and calibrationVersion so a
  historical prediction can be reconstructed exactly as it was produced.

## Specification & Development Workflow

- The framework MUST NOT impose a single specification format. It MUST be able to
  consume SDD, BDD, DDD, TDD, Markdown, and future structured formats without
  confusing specification content with estimation output.
- The first specification integration is GitHub Spec Kit (this repository's
  `.specify/` workflow); other formats are added without modifying the estimation
  core.
- Planning artifacts (components, modules, services, APIs, dependencies, testing and
  deployment strategy, migrations, uncertainties, pending decisions) are a primary
  evidence source for the Estimation Engine and MUST be captured structurally, not
  freeform-only, wherever the plan intends to derive features from them.
- Calibration analysis (absolute/relative error, bias, prediction-interval coverage,
  error by category/project/change-type) MUST be reproducible from stored
  EstimationRecord and OutcomeRecord data alone — not from ad hoc recomputation.

## Governance

This constitution supersedes ad hoc practice for this repository. Every
`/speckit-plan` and `/speckit-tasks` output MUST be checked against these principles
before implementation begins; a plan that violates a principle MUST either be revised
or carry an explicit, justified exception recorded in the plan itself.

Amendments follow semantic versioning:
- MAJOR: backward-incompatible removal or redefinition of a principle.
- MINOR: a new principle or materially expanded section is added.
- PATCH: wording, clarification, or non-semantic fixes.

Amendments are made via this same file, with a Sync Impact Report prepended to the
change describing the version delta and rationale. Compliance review happens at
`/speckit-plan` time by cross-checking planned work against Core Principles and
Estimation Engine Constraints above.

**Version**: 1.0.0 | **Ratified**: 2026-08-31 | **Last Amended**: 2026-08-31
