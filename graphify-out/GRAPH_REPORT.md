# Graph Report - delivery-effort-estimator  (2026-08-31)

## Corpus Check
- 29 files · ~30,982 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 221 nodes · 248 edges · 20 communities
- Extraction: 98% EXTRACTED · 2% INFERRED · 0% AMBIGUOUS · INFERRED: 5 edges (avg confidence: 0.5)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `f42d1d85`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- common.sh
- Tasks: [FEATURE NAME]
- Execution Steps
- speckit-analyze/SKILL.md
- 4. Detection Passes (Token-Efficient Analysis)
- Feature Specification: [FEATURE NAME]
- speckit-plan/SKILL.md
- speckit-specify/SKILL.md
- speckit-tasks/SKILL.md
- Core Principles
- Core Principles
- create-new-feature.sh script
- Implementation Plan: [FEATURE]
- speckit-checklist/SKILL.md
- speckit-clarify/SKILL.md
- speckit-implement/SKILL.md
- speckit-constitution/SKILL.md
- speckit-taskstoissues/SKILL.md
- [CHECKLIST TYPE] Checklist: [FEATURE NAME]

## God Nodes (most connected - your core abstractions)
1. `Tasks: [FEATURE NAME]` - 13 edges
2. `create-new-feature.sh script` - 10 edges
3. `get_repo_root()` - 8 edges
4. `get_feature_paths()` - 8 edges
5. `resolve_template_content()` - 8 edges
6. `setup-tasks.sh script` - 8 edges
7. `check-prerequisites.sh script` - 7 edges
8. `Execution Steps` - 7 edges
9. `4. Detection Passes (Token-Efficient Analysis)` - 7 edges
10. `Execution Steps` - 7 edges

## Surprising Connections (you probably didn't know these)
- `create-new-feature.sh script` --calls--> `get_repo_root()`  [EXTRACTED]
  .specify/scripts/bash/create-new-feature.sh → .specify/scripts/bash/common.sh
- `create-new-feature.sh script` --calls--> `_persist_feature_json()`  [EXTRACTED]
  .specify/scripts/bash/create-new-feature.sh → .specify/scripts/bash/common.sh
- `create-new-feature.sh script` --calls--> `resolve_template_content()`  [EXTRACTED]
  .specify/scripts/bash/create-new-feature.sh → .specify/scripts/bash/common.sh
- `check-prerequisites.sh script` --calls--> `check_dir()`  [EXTRACTED]
  .specify/scripts/bash/check-prerequisites.sh → .specify/scripts/bash/common.sh
- `check-prerequisites.sh script` --calls--> `check_file()`  [EXTRACTED]
  .specify/scripts/bash/check-prerequisites.sh → .specify/scripts/bash/common.sh

## Import Cycles
- None detected.

## Communities (20 total, 0 thin omitted)

### Community 0 - "common.sh"
Cohesion: 0.17
Nodes (22): check-prerequisites.sh script, check_dir(), check_file(), find_specify_root(), format_speckit_command(), get_current_branch(), get_feature_paths(), get_invoke_separator() (+14 more)

### Community 1 - "Tasks: [FEATURE NAME]"
Cohesion: 0.07
Nodes (26): Dependencies & Execution Order, Format: `[ID] [P?] [Story] Description`, Implementation for User Story 1, Implementation for User Story 2, Implementation for User Story 3, Implementation Strategy, Incremental Delivery, MVP First (User Story 1 Only) (+18 more)

### Community 2 - "Execution Steps"
Cohesion: 0.12
Nodes (15): 1. Initialize Convergence Context, 2. Load Artifacts (Progressive Disclosure), 3. Build the Intent Inventory, 4. Assess the Codebase and Classify Findings, 5. Assign Severity, 6. Present the In-Session Findings Summary, 7. Append Convergence Tasks (or report converged), 8. Provide Next Actions (Handoff) (+7 more)

### Community 3 - "speckit-analyze/SKILL.md"
Cohesion: 0.15
Nodes (12): 7. Provide Next Actions, 8. Offer Remediation, 9. Check for extension hooks, Analysis Guidelines, Context, Context Efficiency, Goal, Operating Constraints (+4 more)

### Community 4 - "4. Detection Passes (Token-Efficient Analysis)"
Cohesion: 0.15
Nodes (13): 1. Initialize Analysis Context, 2. Load Artifacts (Progressive Disclosure), 3. Build Semantic Models, 4. Detection Passes (Token-Efficient Analysis), 5. Severity Assignment, 6. Produce Compact Analysis Report, A. Duplication Detection, B. Ambiguity Detection (+5 more)

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
Cohesion: 0.18
Nodes (10): Core Principles, Governance, [PRINCIPLE_1_NAME], [PRINCIPLE_2_NAME], [PRINCIPLE_3_NAME], [PRINCIPLE_4_NAME], [PRINCIPLE_5_NAME], [PROJECT_NAME] Constitution (+2 more)

### Community 10 - "Core Principles"
Cohesion: 0.18
Nodes (10): Core Principles, Governance, [PRINCIPLE_1_NAME], [PRINCIPLE_2_NAME], [PRINCIPLE_3_NAME], [PRINCIPLE_4_NAME], [PRINCIPLE_5_NAME], [PROJECT_NAME] Constitution (+2 more)

### Community 11 - "create-new-feature.sh script"
Cohesion: 0.44
Nodes (7): clean_branch_name(), fit_branch_name(), generate_branch_name(), get_highest_from_specs(), is_feature_number_in_range(), create-new-feature.sh script, spec_prefix_exists()

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

## Knowledge Gaps
- **143 isolated node(s):** `common.sh script`, `User Input`, `Pre-Execution Checks`, `Goal`, `Operating Constraints` (+138 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 151 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Execution Steps` connect `4. Detection Passes (Token-Efficient Analysis)` to `speckit-analyze/SKILL.md`?**
  _High betweenness centrality (0.008) - this node is a cross-community bridge._
- **What connects `common.sh script`, `User Input`, `Pre-Execution Checks` to the rest of the system?**
  _143 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Tasks: [FEATURE NAME]` be split into smaller, more focused modules?**
  _Cohesion score 0.07407407407407407 - nodes in this community are weakly interconnected._
- **Should `Execution Steps` be split into smaller, more focused modules?**
  _Cohesion score 0.125 - nodes in this community are weakly interconnected._