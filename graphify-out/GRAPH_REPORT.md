# Graph Report - delivery-effort-estimator  (2026-08-31)

## Corpus Check
- 99 files · ~61,976 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 678 nodes · 963 edges · 53 communities (41 shown, 12 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 83 edges (avg confidence: 0.83)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `a5f501a5`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- common.sh
- Tasks: [FEATURE NAME]
- Execution Steps
- speckit-analyze/SKILL.md
- testing.T
- Feature Specification: [FEATURE NAME]
- speckit-plan/SKILL.md
- speckit-specify/SKILL.md
- speckit-tasks/SKILL.md
- Core Principles
- Core Principles
- Implementation Plan: Core Estimation Engine
- Implementation Plan: [FEATURE]
- speckit-checklist/SKILL.md
- speckit-clarify/SKILL.md
- speckit-implement/SKILL.md
- speckit-constitution/SKILL.md
- speckit-taskstoissues/SKILL.md
- [CHECKLIST TYPE] Checklist: [FEATURE NAME]
- Quickstart: Core Estimation Engine
- Predict
- Tasks: Core Estimation Engine
- Store
- github.com/jmirandev/delivery-effort-estimator/internal/record.Service
- estimatorctl/main.go
- Phase 1 Data Model: Core Estimation Engine
- Phase 0 Research: Core Estimation Engine
- Contract: `estimatorctl` CLI
- NewService
- Implementation Plan: Feature Extraction Plugin
- Service
- Predict
- Tasks: Feature Extraction Plugin
- engine/cmd/estimatorctl/main.go
- Phase 0 Research: Feature Extraction Plugin
- Quickstart: Feature Extraction Plugin
- Contract: `on-tasks-md-change.sh` (PostToolUse hook)
- Contract: `estimate-feature` skill
- Data Model: Feature Extraction Plugin
- estimate-feature/SKILL.md
- Contract: `plugin/scripts/ensure-engine.sh`
- internal/record/repository.go
- engine/internal/record/repository.go
- test_ensure_engine.sh
- test_locate_feature.sh
- test_on_tasks_md_change.sh
- ensure-engine.sh
- locate-feature.sh
- on-tasks-md-change.sh
- check-plugin-engine-sync.sh
- sync-plugin-engine.sh
- github.com/jmirandev/delivery-effort-estimator
- github.com/jmirandev/delivery-effort-estimator

## God Nodes (most connected - your core abstractions)
1. `Tasks: [FEATURE NAME]` - 13 edges
2. `Tasks: Core Estimation Engine` - 13 edges
3. `NewService()` - 12 edges
4. `Service` - 11 edges
5. `Service` - 11 edges
6. `Core Principles` - 11 edges
7. `create-new-feature.sh script` - 10 edges
8. `Predict()` - 10 edges
9. `NewService()` - 10 edges
10. `Predict()` - 10 edges

## Surprising Connections (you probably didn't know these)
- `minimalEstimationRecord()` --calls--> `Predict()`  [INFERRED]
  internal/storage/sqlite/sqlite_test.go → internal/estimation/model.go
- `openTestStore()` --calls--> `Open()`  [INFERRED]
  internal/storage/sqlite/sqlite_test.go → internal/storage/sqlite/sqlite.go
- `main()` --calls--> `NewService()`  [INFERRED]
  plugin/engine/cmd/estimatorctl/main.go → plugin/engine/internal/record/service.go
- `main()` --calls--> `Open()`  [INFERRED]
  plugin/engine/cmd/estimatorctl/main.go → plugin/engine/internal/storage/sqlite/sqlite.go
- `handleEstimateError()` --calls--> `IsMissingFeatureError()`  [INFERRED]
  plugin/engine/cmd/estimatorctl/main.go → plugin/engine/internal/estimation/features.go

## Import Cycles
- None detected.

## Communities (53 total, 12 thin omitted)

### Community 0 - "common.sh"
Cohesion: 0.13
Nodes (29): check-prerequisites.sh script, check_dir(), check_file(), find_specify_root(), format_speckit_command(), get_current_branch(), get_feature_paths(), get_invoke_separator() (+21 more)

### Community 1 - "Tasks: [FEATURE NAME]"
Cohesion: 0.07
Nodes (26): Dependencies & Execution Order, Format: `[ID] [P?] [Story] Description`, Implementation for User Story 1, Implementation for User Story 2, Implementation for User Story 3, Implementation Strategy, Incremental Delivery, MVP First (User Story 1 Only) (+18 more)

### Community 2 - "Execution Steps"
Cohesion: 0.12
Nodes (15): 1. Initialize Convergence Context, 2. Load Artifacts (Progressive Disclosure), 3. Build the Intent Inventory, 4. Assess the Codebase and Classify Findings, 5. Assign Severity, 6. Present the In-Session Findings Summary, 7. Append Convergence Tasks (or report converged), 8. Provide Next Actions (Handoff) (+7 more)

### Community 3 - "speckit-analyze/SKILL.md"
Cohesion: 0.08
Nodes (25): 1. Initialize Analysis Context, 2. Load Artifacts (Progressive Disclosure), 3. Build Semantic Models, 4. Detection Passes (Token-Efficient Analysis), 5. Severity Assignment, 6. Produce Compact Analysis Report, 7. Provide Next Actions, 8. Offer Remediation (+17 more)

### Community 4 - "testing.T"
Cohesion: 0.07
Nodes (41): testing.T, EstimationFeatures, MissingFeatureError, IsMissingFeatureError(), ParseFeatures(), TestClampOutOfRangeValues(), TestClampWithinRangeProducesNoAssumptions(), TestParseFeaturesCompleteVector() (+33 more)

### Community 5 - "Feature Specification: [FEATURE NAME]"
Cohesion: 0.15
Nodes (12): Assumptions, Edge Cases, Feature Specification: [FEATURE NAME], Functional Requirements, Key Entities *(include if feature involves data)*, Measurable Outcomes, Requirements *(mandatory)*, Success Criteria *(mandatory)* (+4 more)

### Community 6 - "speckit-plan/SKILL.md"
Cohesion: 0.18
Nodes (10): Completion Report, Done When, Key rules, Mandatory Post-Execution Hooks, Outline, Phase 0: Outline & Research, Phase 1: Design & Contracts, Phases (+2 more)

### Community 7 - "speckit-specify/SKILL.md"
Cohesion: 0.18
Nodes (10): Completion Report, Done When, For AI Generation, Mandatory Post-Execution Hooks, Outline, Pre-Execution Checks, Quick Guidelines, Section Requirements (+2 more)

### Community 8 - "speckit-tasks/SKILL.md"
Cohesion: 0.18
Nodes (10): Checklist Format (REQUIRED), Completion Report, Done When, Mandatory Post-Execution Hooks, Outline, Phase Structure, Pre-Execution Checks, Task Generation Rules (+2 more)

### Community 9 - "Core Principles"
Cohesion: 0.12
Nodes (15): Core Principles, Delivery Effort Estimator Constitution, Estimation Engine Constraints, Governance, I. Evidence Over Intuition, II. Effort Is Not Duration, III. Effort Is Not Cost, IV. Business Value Is Not Delivery Effort (+7 more)

### Community 10 - "Core Principles"
Cohesion: 0.18
Nodes (10): Core Principles, Governance, [PRINCIPLE_1_NAME], [PRINCIPLE_2_NAME], [PRINCIPLE_3_NAME], [PRINCIPLE_4_NAME], [PRINCIPLE_5_NAME], [PROJECT_NAME] Constitution (+2 more)

### Community 11 - "Implementation Plan: Core Estimation Engine"
Cohesion: 0.07
Nodes (25): Content Quality, Feature Readiness, Notes, Requirement Completeness, Specification Quality Checklist: Core Estimation Engine, Complexity Tracking, Constitution Check, Documentation (this feature) (+17 more)

### Community 12 - "Implementation Plan: [FEATURE]"
Cohesion: 0.22
Nodes (8): Complexity Tracking, Constitution Check, Documentation (this feature), Implementation Plan: [FEATURE], Project Structure, Source Code (repository root), Summary, Technical Context

### Community 13 - "speckit-checklist/SKILL.md"
Cohesion: 0.25
Nodes (7): Anti-Examples: What NOT To Do, Checklist Purpose: "Unit Tests for English", Example Checklist Types & Sample Items, Execution Steps, Post-Execution Checks, Pre-Execution Checks, User Input

### Community 14 - "speckit-clarify/SKILL.md"
Cohesion: 0.29
Nodes (6): Completion Report, Done When, Mandatory Post-Execution Hooks, Outline, Pre-Execution Checks, User Input

### Community 15 - "speckit-implement/SKILL.md"
Cohesion: 0.29
Nodes (6): Completion Report, Done When, Mandatory Post-Execution Hooks, Outline, Pre-Execution Checks, User Input

### Community 16 - "speckit-constitution/SKILL.md"
Cohesion: 0.33
Nodes (5): Outline, Post-Execution Checks, Pre-Execution Checks, Scope Guard, User Input

### Community 17 - "speckit-taskstoissues/SKILL.md"
Cohesion: 0.40
Nodes (4): Outline, Post-Execution Checks, Pre-Execution Checks, User Input

### Community 18 - "[CHECKLIST TYPE] Checklist: [FEATURE NAME]"
Cohesion: 0.40
Nodes (4): [Category 1], [Category 2], [CHECKLIST TYPE] Checklist: [FEATURE NAME], Notes

### Community 19 - "Quickstart: Core Estimation Engine"
Cohesion: 0.07
Nodes (24): Delivery Effort Estimator — Claude Code Plugin, Directory layout, How it works, Install, What happens on first use, Build & test, Claude Code plugin, Delivery Effort Estimator (+16 more)

### Community 20 - "Predict"
Cohesion: 0.15
Nodes (18): Compute(), ActualOutcome, DimensionError, ErrorReport, Prediction, TestComputeDimensionErrors(), clampRange(), DeriveRisks() (+10 more)

### Community 21 - "Tasks: Core Estimation Engine"
Cohesion: 0.08
Nodes (24): Dependencies & Execution Order, Format: `[ID] [P?] [Story] Description`, Implementation for User Story 1, Implementation for User Story 2, Implementation for User Story 3, Implementation Strategy, Incremental Delivery, MVP First (User Story 1 Only) (+16 more)

### Community 22 - "Store"
Cohesion: 0.11
Nodes (16): database/sql.DB, github.com/jmirandev/delivery-effort-estimator/internal/record.EstimationRecord, github.com/jmirandev/delivery-effort-estimator/internal/record.OutcomeRecord, migrate(), Store, Open(), migrate(), Store (+8 more)

### Community 23 - "github.com/jmirandev/delivery-effort-estimator/internal/record.Service"
Cohesion: 0.35
Nodes (17): handleErrorReport(), handleEstimate(), handleGetEstimation(), handleRecordOutcome(), main(), writeError(), writeJSON(), github.com/jmirandev/delivery-effort-estimator/internal/record.Service (+9 more)

### Community 24 - "estimatorctl/main.go"
Cohesion: 0.38
Nodes (13): dbPath(), emitError(), fatal(), handleEstimateError(), handleOutcomeError(), handleReportError(), main(), printJSON() (+5 more)

### Community 25 - "Phase 1 Data Model: Core Estimation Engine"
Cohesion: 0.22
Nodes (8): DimensionError, ErrorReport, EstimationFeatures, EstimationRecord, OutcomeRecord, Phase 1 Data Model: Core Estimation Engine, Prediction, State Transitions

### Community 26 - "Phase 0 Research: Core Estimation Engine"
Cohesion: 0.25
Nodes (7): Decision: Append-only schema enforces Constitution Principle VIII in the storage layer, not just in application logic, Decision: Cross-dimension error comparison uses raw numeric deltas, not unit conversion, Decision: Deterministic v1 model — fixed linear weights over the 7 SDD features, Decision: Language & runtime — Go 1.27, standard library only + one pure-Go driver, Decision: Storage — embedded SQLite behind a repository interface, Graphify dependency-graph check (Constitution Principle X), Phase 0 Research: Core Estimation Engine

### Community 27 - "Contract: `estimatorctl` CLI"
Cohesion: 0.40
Nodes (4): Contract: `estimatorctl` CLI, `estimatorctl error-report`, `estimatorctl estimate`, `estimatorctl record-outcome`

### Community 28 - "NewService"
Cohesion: 0.13
Nodes (19): EstimationRecord, OutcomeRecord, fakeEstimationRepo, fakeOutcomeRepo, newFakeEstimationRepo(), newFakeOutcomeRepo(), EstimationRecord, EstimationRepository (+11 more)

### Community 29 - "Implementation Plan: Feature Extraction Plugin"
Cohesion: 0.07
Nodes (26): Content Quality, Feature Readiness, Notes, Requirement Completeness, Specification Quality Checklist: Feature Extraction Plugin, Complexity Tracking, Constitution Check, Documentation (this feature) (+18 more)

### Community 30 - "Service"
Cohesion: 0.14
Nodes (14): github.com/jmirandev/delivery-effort-estimator/internal/estimation.ErrorReport, github.com/jmirandev/delivery-effort-estimator/internal/estimation.EstimationFeatures, github.com/jmirandev/delivery-effort-estimator/internal/estimation.Prediction, time.Time, EstimationRecord, EstimationRepository, OutcomeRecord, OutcomeRepository (+6 more)

### Community 31 - "Predict"
Cohesion: 0.15
Nodes (18): Compute(), ActualOutcome, DimensionError, ErrorReport, Prediction, TestComputeDimensionErrors(), clampRange(), DeriveRisks() (+10 more)

### Community 32 - "Tasks: Feature Extraction Plugin"
Cohesion: 0.11
Nodes (17): Dependencies & Execution Order, Format: `[ID] [P?] [Story] Description`, Implementation for User Story 1, Implementation for User Story 2, Implementation for User Story 3, Implementation Strategy, Incremental Delivery, MVP First (User Story 1 Only) (+9 more)

### Community 33 - "engine/cmd/estimatorctl/main.go"
Cohesion: 0.38
Nodes (13): dbPath(), emitError(), fatal(), handleEstimateError(), handleOutcomeError(), handleReportError(), main(), printJSON() (+5 more)

### Community 34 - "Phase 0 Research: Feature Extraction Plugin"
Cohesion: 0.22
Nodes (8): Decision: Automatic trigger via a PostToolUse hook that injects context, not a hook that reasons itself, Decision: Derivation rubric — concrete anchors per dimension, not open-ended judgment, Decision: Engine freshness check via a checksum marker, not a semantic version bump per sync, Decision: Feature-directory resolution is path-convention-based, not conversation-memory-based, Decision: Hook stdin parsing uses `grep`/`sed`, not `jq`, Decision: Output files live under a per-feature `estimation/` subdirectory, one timestamped pair per run, Graphify dependency-graph check (Constitution Principle X), Phase 0 Research: Feature Extraction Plugin

### Community 35 - "Quickstart: Feature Extraction Plugin"
Cohesion: 0.25
Nodes (7): Automated coverage (scriptable), Prerequisites, Quickstart: Feature Extraction Plugin, Scenario 1 — Automatic estimation after planning (User Story 1), Scenario 2 — Derivation is auditable (User Story 2), Scenario 3 — Skipped, not broken, when preconditions aren't met (SC-005), Scenario 4 — Portable install (User Story 3)

### Community 36 - "Contract: `on-tasks-md-change.sh` (PostToolUse hook)"
Cohesion: 0.40
Nodes (4): Behavior, Contract: `on-tasks-md-change.sh` (PostToolUse hook), Input (stdin), Output contract

### Community 37 - "Contract: `estimate-feature` skill"
Cohesion: 0.40
Nodes (4): Contract: `estimate-feature` skill, Postconditions, Preconditions, Steps (contract, not prose — see SKILL.md for the full rubric text)

### Community 38 - "Data Model: Feature Extraction Plugin"
Cohesion: 0.40
Nodes (4): Data Model: Feature Extraction Plugin, DerivedFeatureVector, PluginEngineBuild (operational state, not a domain entity), Reused entities (specs/001-core-estimation-engine, unchanged)

### Community 39 - "estimate-feature/SKILL.md"
Cohesion: 0.50
Nodes (3): Never do this, Purpose, Steps

### Community 40 - "Contract: `plugin/scripts/ensure-engine.sh`"
Cohesion: 0.50
Nodes (3): Behavior, Contract: `plugin/scripts/ensure-engine.sh`, Output contract

## Knowledge Gaps
- **287 isolated node(s):** `common.sh script`, `github.com/jmirandev/delivery-effort-estimator`, `EstimationRepository`, `OutcomeRepository`, `github.com/jmirandev/delivery-effort-estimator` (+282 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 325 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **12 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `NewService()` connect `NewService` to `engine/cmd/estimatorctl/main.go`, `github.com/jmirandev/delivery-effort-estimator/internal/record.Service`?**
  _High betweenness centrality (0.024) - this node is a cross-community bridge._
- **Why does `Service` connect `NewService` to `Service`?**
  _High betweenness centrality (0.016) - this node is a cross-community bridge._
- **Are the 8 inferred relationships involving `NewService()` (e.g. with `main()` and `main()`) actually correct?**
  _`NewService()` has 8 INFERRED edges - model-reasoned connections that need verification._
- **What connects `common.sh script`, `github.com/jmirandev/delivery-effort-estimator`, `EstimationRepository` to the rest of the system?**
  _287 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `common.sh` be split into smaller, more focused modules?**
  _Cohesion score 0.12698412698412698 - nodes in this community are weakly interconnected._
- **Should `Tasks: [FEATURE NAME]` be split into smaller, more focused modules?**
  _Cohesion score 0.07407407407407407 - nodes in this community are weakly interconnected._
- **Should `Execution Steps` be split into smaller, more focused modules?**
  _Cohesion score 0.125 - nodes in this community are weakly interconnected._