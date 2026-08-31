# Graph Report - delivery-effort-estimator  (2026-08-31)

## Corpus Check
- 57 files · ~45,334 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 434 nodes · 607 edges · 29 communities (28 shown, 1 thin omitted)
- Extraction: 94% EXTRACTED · 6% INFERRED · 0% AMBIGUOUS · INFERRED: 36 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `93e8cfd6`
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
- EstimationRecord
- Service
- estimatorctl/main.go
- Phase 1 Data Model: Core Estimation Engine
- Phase 0 Research: Core Estimation Engine
- Contract: `estimatorctl` CLI
- github.com/jmirandev/delivery-effort-estimator

## God Nodes (most connected - your core abstractions)
1. `Service` - 19 edges
2. `EstimationRecord` - 13 edges
3. `Tasks: [FEATURE NAME]` - 13 edges
4. `Tasks: Core Estimation Engine` - 13 edges
5. `NewService()` - 12 edges
6. `Core Principles` - 11 edges
7. `create-new-feature.sh script` - 10 edges
8. `Predict()` - 10 edges
9. `main()` - 9 edges
10. `OutcomeRecord` - 9 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `NewService()`  [EXTRACTED]
  cmd/estimatorctl/main.go → internal/record/service.go
- `main()` --calls--> `Open()`  [EXTRACTED]
  cmd/estimatorctl/main.go → internal/storage/sqlite/sqlite.go
- `runEstimate()` --references--> `Service`  [EXTRACTED]
  cmd/estimatorctl/main.go → internal/record/service.go
- `runGet()` --references--> `Service`  [EXTRACTED]
  cmd/estimatorctl/main.go → internal/record/service.go
- `runRecordOutcome()` --references--> `Service`  [EXTRACTED]
  cmd/estimatorctl/main.go → internal/record/service.go

## Import Cycles
- None detected.

## Communities (29 total, 1 thin omitted)

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
Cohesion: 0.14
Nodes (25): MissingFeatureError, testing.T, IsMissingFeatureError(), ParseFeatures(), TestClampOutOfRangeValues(), TestClampWithinRangeProducesNoAssumptions(), TestParseFeaturesCompleteVector(), TestParseFeaturesInvalidJSON() (+17 more)

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
Cohesion: 0.09
Nodes (18): Build & test, Delivery Effort Estimator, Project layout, Try it, What's implemented, Contract: `estimatord` HTTP JSON API, `GET /estimations/{estimationId}`, `GET /estimations/{estimationId}/error-report` (+10 more)

### Community 20 - "Predict"
Cohesion: 0.14
Nodes (17): ActualOutcome, DimensionError, Compute(), ErrorReport, TestComputeDimensionErrors(), EstimationFeatures, clampRange(), DeriveRisks() (+9 more)

### Community 21 - "Tasks: Core Estimation Engine"
Cohesion: 0.08
Nodes (24): Dependencies & Execution Order, Format: `[ID] [P?] [Story] Description`, Implementation for User Story 1, Implementation for User Story 2, Implementation for User Story 3, Implementation Strategy, Incremental Delivery, MVP First (User Story 1 Only) (+16 more)

### Community 22 - "EstimationRecord"
Cohesion: 0.15
Nodes (9): database/sql.DB, time.Time, EstimationRecord, OutcomeRecord, migrate(), Open(), fakeEstimationRepo, fakeOutcomeRepo (+1 more)

### Community 23 - "Service"
Cohesion: 0.25
Nodes (13): handleErrorReport(), handleEstimate(), handleGetEstimation(), handleRecordOutcome(), main(), writeError(), writeJSON(), net/http.HandlerFunc (+5 more)

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

## Knowledge Gaps
- **215 isolated node(s):** `common.sh script`, `github.com/jmirandev/delivery-effort-estimator`, `User Input`, `Pre-Execution Checks`, `Goal` (+210 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 232 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **1 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Service` connect `Service` to `estimatorctl/main.go`, `Predict`, `testing.T`, `EstimationRecord`?**
  _High betweenness centrality (0.018) - this node is a cross-community bridge._
- **Why does `EstimationRecord` connect `EstimationRecord` to `Predict`, `testing.T`, `Service`?**
  _High betweenness centrality (0.009) - this node is a cross-community bridge._
- **Why does `NewService()` connect `testing.T` to `estimatorctl/main.go`, `Service`?**
  _High betweenness centrality (0.008) - this node is a cross-community bridge._
- **What connects `common.sh script`, `github.com/jmirandev/delivery-effort-estimator`, `User Input` to the rest of the system?**
  _215 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `common.sh` be split into smaller, more focused modules?**
  _Cohesion score 0.12698412698412698 - nodes in this community are weakly interconnected._
- **Should `Tasks: [FEATURE NAME]` be split into smaller, more focused modules?**
  _Cohesion score 0.07407407407407407 - nodes in this community are weakly interconnected._
- **Should `Execution Steps` be split into smaller, more focused modules?**
  _Cohesion score 0.125 - nodes in this community are weakly interconnected._