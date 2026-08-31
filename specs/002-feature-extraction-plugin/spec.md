# Feature Specification: Feature Extraction Plugin

**Feature Branch**: `002-feature-extraction-plugin`

**Created**: 2026-08-31

**Status**: Draft

**Input**: User description: "Framework integrado a Claude Code: cuando se inicie un desarrollo y se cree un SDD (via Spec Kit), el framework debe usar el motor de estimación más la información del SDD para crear la estimación y estadísticas automáticamente, sin que el usuario tenga que armar el JSON de features a mano."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Automatic prediction right after planning (Priority: P1)

A developer working inside Claude Code on a Spec-Kit-enabled project finishes
specifying and planning a new feature (`spec.md`, `plan.md`, `tasks.md` exist).
Without running any extra command, they see a Delivery Effort prediction for
that feature in the same conversation — derived from what was just written,
not from a hand-built feature-vector JSON.

**Why this priority**: This is the entire point of the integration described
in the SDD (§21 "Integración inicial"): the estimation engine is only useful
in practice if it's fed automatically from the artifacts a Spec-Kit workflow
already produces. Without this, users are back to manually authoring JSON,
which the SDD explicitly says should not be necessary (§10).

**Independent Test**: In a Spec-Kit project with this plugin installed, run
through `/speckit.specify` → `/speckit.plan` → `/speckit.tasks` for a real
feature. Confirm a Delivery Effort prediction appears in the same session
after `tasks.md` is produced, without the user supplying a features JSON.

**Acceptance Scenarios**:

1. **Given** a feature directory with `spec.md`, `plan.md`, and `tasks.md`
   just completed, **When** `tasks.md` is created, **Then** the system
   derives a normalized feature vector from those three files and produces a
   Delivery Effort prediction without further user input.
2. **Given** a project without Spec Kit initialized (no `.specify/`),
   **When** any file is created or edited, **Then** the system does not
   attempt estimation and does not interrupt the user's normal workflow.
3. **Given** a feature whose `tasks.md` is edited again after an estimation
   already exists for it, **When** the edit is saved, **Then** a new,
   independently timestamped estimation is produced and the earlier one
   remains unchanged and retrievable.

---

### User Story 2 - Auditable derivation, not a black box (Priority: P2)

A developer reviews *why* the system scored a feature's complexity, human
decision load, or uncertainty the way it did, so they can sanity-check or
challenge the derivation instead of trusting an opaque number.

**Why this priority**: The SDD's foundational principle (§3, "Evidence over
intuition") only holds if the evidence trail is visible. An automatically
derived feature vector that can't be inspected is worse than a manual one,
because it *looks* authoritative while hiding its reasoning.

**Independent Test**: After an automatic estimation runs, open the feature
directory and find a file with the 7 derived dimension values, each paired
with a short justification referencing the spec/plan/tasks content it came
from.

**Acceptance Scenarios**:

1. **Given** an automatically derived feature vector, **When** a developer
   inspects the feature directory, **Then** each of the 7 dimensions has an
   accompanying human-readable justification.
2. **Given** a derived feature vector a developer disagrees with, **When**
   they read the justification, **Then** they can trace the score back to a
   specific statement in `spec.md`, `plan.md`, or `tasks.md`.

---

### User Story 3 - Drop the plugin into any Spec-Kit project (Priority: P3)

A team lead installs this integration into a different Spec-Kit-enabled
Claude Code project (not this repository) and it works without a separate,
manual setup of the estimation engine.

**Why this priority**: The SDD (§21) requires the first integration to work
with Claude Code and Spec Kit generally, not only inside this repository.
This is what turns the engine from "a tool in this repo" into "a framework
teams install."

**Independent Test**: Install the plugin from a local directory into a
freshly scaffolded Spec-Kit project (`claude --plugin-dir <path>`), complete
one feature's specify→plan→tasks flow, and confirm an estimation is produced
with no manual engine build/config step beyond having Go installed.

**Acceptance Scenarios**:

1. **Given** a fresh Spec-Kit project with the plugin installed and Go
   available on the machine, **When** the plugin is used for the first time,
   **Then** it builds the estimation engine binary itself and reuses it on
   subsequent runs.
2. **Given** a machine without a Go toolchain, **When** the plugin attempts
   its first build, **Then** the developer sees a clear, actionable error
   naming the missing dependency, and no partial or fabricated estimation is
   produced.

---

### Edge Cases

- What happens when `tasks.md` is created for a feature whose `spec.md` or
  `plan.md` is missing or clearly incomplete? → Derivation proceeds
  best-effort from what exists; every dimension not clearly supported by
  content is scored conservatively and flagged as an assumption, consistent
  with the existing engine's `Clamp()`/assumption-note behavior for
  out-of-range or default input.
- What happens when a project has several Spec-Kit features in progress at
  once? → The system must resolve estimation to the exact feature directory
  the changed `tasks.md` belongs to, never to "whichever feature was worked
  on last" in conversation.
- What happens when the estimation engine binary needs to be rebuilt (e.g.,
  engine source inside the plugin was updated to a new version)? → The
  system detects the version mismatch and rebuilds rather than silently
  running a stale binary.
- What happens if the automatic derivation is clearly wrong or absurd (e.g.
  Claude misreads the spec)? → Because User Story 2 makes the derivation
  visible and re-running is non-destructive (append-only), a developer can
  request a redo; this feature does not need a built-in "correct/override"
  mechanism beyond what a manual re-trigger already provides.
- What happens when `tasks.md` is created outside of a `specs/*` Spec-Kit
  layout (e.g., an unrelated file happens to be named `tasks.md`)? → The
  system only acts on paths that match the Spec-Kit feature-directory
  convention recorded via `.specify/feature.json` or the `specs/<id>-*/`
  layout; anything else is ignored.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST detect, without requiring a user-issued command,
  when a Spec-Kit feature's `tasks.md` is created or modified in a project.
- **FR-002**: Upon detection, system MUST derive a normalized 7-dimension
  feature vector (context_complexity, domain_complexity,
  integration_complexity, verification_complexity, human_decision_load,
  ai_execution_complexity, uncertainty) by reasoning over that feature's
  `spec.md`, `plan.md`, and `tasks.md` content.
- **FR-003**: System MUST record, alongside the derived feature vector, a
  human-readable justification for each of the 7 dimension scores.
- **FR-004**: System MUST persist the derived feature vector (values +
  justifications) as a file inside the feature's own spec directory.
- **FR-005**: System MUST obtain the Delivery Effort prediction by calling
  into the existing, independent Estimation Engine (`specs/001-core-estimation-engine`)
  with the derived feature vector — it MUST NOT reimplement or duplicate the
  engine's scoring model.
- **FR-006**: System MUST persist the resulting estimation record inside the
  feature's spec directory, associated with that feature.
- **FR-007**: System MUST NOT overwrite or delete a previously generated
  feature-vector derivation or estimation for a feature directory when
  `tasks.md` changes again; each recomputation MUST produce an additional,
  independently timestamped derivation and estimation.
- **FR-008**: If the project has no Spec Kit initialized, or the estimation
  engine cannot be made available, system MUST skip automatic estimation
  without interrupting the user's normal workflow, and MUST make clear why it
  was skipped.
- **FR-009**: System MUST also allow a user to explicitly (manually) trigger
  extraction and estimation for a chosen feature directory, independent of
  the automatic trigger.
- **FR-010**: System MUST be installable into any Spec-Kit-enabled Claude
  Code project as a self-contained unit, without a separate manual setup
  procedure for the estimation engine beyond having a Go toolchain present.
- **FR-011**: On first use in a project, if the estimation engine binary is
  not already built (or is stale relative to the bundled engine version),
  system MUST build it automatically from bundled source using the local Go
  toolchain, and reuse it on subsequent runs without rebuilding.
- **FR-012**: If the local Go toolchain is unavailable, system MUST report a
  clear, actionable error naming the missing dependency, rather than failing
  silently or producing an estimation from incomplete data.
- **FR-013**: The feature-vector derivation MUST follow a documented,
  consistent rubric mapping specification/plan/tasks content to each of the
  7 dimensions, so scoring is explainable and comparable across features
  rather than ad hoc per invocation.
- **FR-014**: System MUST correctly identify which feature directory a given
  `tasks.md` change belongs to, even when a project has multiple Spec-Kit
  features in progress, and MUST NOT estimate the wrong feature.
- **FR-015**: System MUST surface the newly produced prediction to the user
  within the same conversation turn, not only save it silently to disk.

### Key Entities

- **DerivedFeatureVector**: The 7 normalized [0,1] dimensions defined by the
  Estimation Engine (`EstimationFeatures`), plus a per-dimension
  justification string and a reference to which feature directory / files
  it was derived from and when. This is new to this feature; it feeds the
  existing engine, it does not replace `EstimationFeatures`.
- **EstimationRecord / Prediction**: Reused as-is from
  `specs/001-core-estimation-engine`; this feature is a producer of their
  input, not a redefinition of them.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: After completing specification and planning for a feature in
  Claude Code, a developer receives a Delivery Effort prediction in the same
  session without manually authoring a features JSON.
- **SC-002**: 100% of automatically generated estimations include a
  human-readable justification for each of the 7 scored dimensions.
- **SC-003**: Re-planning a feature never destroys a previously generated
  derivation or estimation — every historical one remains retrievable.
- **SC-004**: Installing the plugin into a new Spec-Kit project and
  completing one feature's specify→plan→tasks flow produces a working
  estimation with no manual engine build/install step beyond having Go
  installed.
- **SC-005**: A project without Spec Kit, or without a working engine build,
  never has its normal workflow blocked or broken by this feature.

## Assumptions

- The Estimation Engine (`specs/001-core-estimation-engine`) is a stable,
  versioned dependency this feature calls into as-is; this feature does not
  change the engine's model, formulas, or interfaces.
- Feature-vector derivation is performed by the same Claude Code
  conversation already producing the spec/plan/tasks, following a documented
  rubric — it is not a separate deterministic text parser. Derivation
  quality is bounded by the quality of the underlying spec/plan/tasks
  content, consistent with the SDD's own framing (§10: values can be derived
  from "especificación, planning, conversación LLM-humano").
- The automatic trigger fires specifically on `tasks.md` creation/update
  (the point at which a feature is considered specified *and* planned), not
  on `spec.md` or `plan.md` alone.
- Requiring a local Go toolchain to build the engine binary on first use is
  acceptable for this feature's target users; distributing precompiled
  binaries per OS/architecture is explicitly deferred (a decision already
  made when scoping this feature).
- Per SDD §21 ("La integración con otros agentes debe ser posible sin
  modificar el núcleo del framework"), this feature only adds a
  Claude-Code-specific glue layer; it must not require changes to the
  engine's own independence from any specific agent or LLM provider.
