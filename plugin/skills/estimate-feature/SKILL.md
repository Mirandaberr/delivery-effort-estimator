---
name: "estimate-feature"
description: "Derive the 7-dimension EstimationFeatures vector from a Spec-Kit feature's spec.md/plan.md/tasks.md and produce a Delivery Effort prediction via estimatorctl."
argument-hint: "Path to a specs/<id>-*/ feature directory (defaults to the directory named in the triggering context, if any)"
compatibility: "Requires a Spec-Kit project (.specify/) with the feature's spec.md present; plan.md/tasks.md improve accuracy but are not required."
metadata:
  author: "delivery-effort-estimator"
user-invocable: true
disable-model-invocation: false
---

## Purpose

Turn a Spec-Kit feature's own artifacts into a Delivery Effort prediction,
without anyone hand-authoring a features JSON. This skill runs either
because the `PostToolUse` hook (`plugin/hooks/hooks.json`) just told you a
feature's `tasks.md` was written, or because a user explicitly asked you to
estimate a feature directory.

**You derive the features. You never compute the prediction yourself.**
The actual Delivery Effort / Confidence / Prediction Interval numbers come
only from running `estimatorctl`, never from your own arithmetic or
judgment — this is a hard rule (Constitution Principle IX: "LLMs Are Not
the Estimation Engine"). If you find yourself about to state a DEU or
Confidence number you computed by hand, stop — run the CLI instead.

## Steps

1. **Identify the feature directory.** Use the directory named in the
   triggering context (hook or user argument). Confirm `spec.md` exists
   there; if not, tell the user you can't proceed and stop. Read
   `plan.md` and `tasks.md` too if they exist — missing ones just mean
   fewer signals, not a hard failure.

2. **Score the 7 dimensions.** For each dimension below, read the primary
   signal source, pick the closest anchor (interpolating if it's between
   two), and write a 1-2 sentence justification quoting or paraphrasing the
   specific content that drove the score. If the relevant section is
   missing or too thin to judge, score conservatively (toward 0.0 for the
   effort-type dimensions, toward 1.0 for `uncertainty`) and say so in the
   justification — that becomes an entry in `assumptions` too.

   | Dimension | Primary signal source | 0.0 anchor | 0.5 anchor | 1.0 anchor |
   |---|---|---|---|---|
   | `context_complexity` | plan.md "components/modules affected"; spec.md scenario breadth | Change confined to one already-understood file/module | Spans a handful of modules within one service | Requires understanding large or unfamiliar parts of the system across services |
   | `domain_complexity` | spec.md problem description, Key Entities | CRUD-like, no domain rules beyond basic validation | Moderate business rules or a non-trivial state machine | Deep domain modeling (compliance, financial/legal rules, complex multi-actor workflows) |
   | `integration_complexity` | plan.md components/APIs/events/DBs/infra/external systems list | No external touchpoints (pure internal logic) | 2-3 touchpoints (e.g., one API + one DB) | 4+ touchpoints or any external/third-party system dependency |
   | `verification_complexity` | plan.md testing strategy; spec.md acceptance scenarios & edge cases | 1-2 straightforward acceptance scenarios, no regression risk | Several scenarios plus at least one edge case needing dedicated test design | Broad scenario coverage, regression-sensitive area, and/or manual validation required |
   | `human_decision_load` | spec.md Assumptions/[NEEDS CLARIFICATION] history; plan.md pending decisions | Fully unambiguous requirements, no judgment calls left | A few reasonable defaults were chosen among real alternatives | Multiple open architectural/product trade-offs the human must actively resolve |
   | `ai_execution_complexity` | tasks.md task count/breadth; plan.md structure decision | A handful of small, mechanical tasks | Moderate task count spanning several files with some non-trivial logic | Large or intricate task set requiring iteration, multi-file reasoning, or nontrivial tool/build loops |
   | `uncertainty` | spec.md Assumptions section; plan.md unresolved dependencies/unknowns | No assumptions recorded, no unresolved dependencies | A few named assumptions, none blocking | Several unresolved unknowns or assumptions the plan itself flags as risky |

   (Full rationale for these anchors: `specs/002-feature-extraction-plugin/research.md`.)

3. **Ensure the engine is built.** Run:
   ```
   <plugin-root>/scripts/ensure-engine.sh
   ```
   If it exits non-zero, stop here. Tell the user exactly what it printed to
   stderr (e.g. "go toolchain not found on PATH") — do not proceed to a
   partial or fabricated estimation (FR-012).

4. **Write the derivation file.** Create
   `<feature-directory>/estimation/<UTC-timestamp>/` (e.g.
   `2026-08-31T18-04-00Z`, colons replaced with `-` for filesystem safety)
   and write `derived-features.json` there, shaped exactly like this — the
   7 dimension keys are **flat, at the top level**, not nested, because this
   same file is passed straight to `estimatorctl` in the next step
   (data-model.md explains why):

   ```json
   {
     "context_complexity": 0.72,
     "domain_complexity": 0.41,
     "integration_complexity": 0.83,
     "verification_complexity": 0.67,
     "human_decision_load": 0.54,
     "ai_execution_complexity": 0.38,
     "uncertainty": 0.61,
     "feature_directory": "specs/002-feature-extraction-plugin",
     "derived_at": "2026-08-31T18:04:00Z",
     "source_files": ["spec.md", "plan.md", "tasks.md"],
     "justifications": {
       "context_complexity": "...",
       "domain_complexity": "...",
       "integration_complexity": "...",
       "verification_complexity": "...",
       "human_decision_load": "...",
       "ai_execution_complexity": "...",
       "uncertainty": "..."
     },
     "assumptions": []
   }
   ```

   Every one of the 7 dimensions MUST have a matching entry in
   `justifications` — do not write the file otherwise.

5. **Run the engine:**
   ```
   <plugin-root>/bin/estimatorctl estimate \
     --work-item <feature-directory-basename> \
     --features <feature-directory>/estimation/<timestamp>/derived-features.json
   ```

6. **Save the result verbatim** as
   `<feature-directory>/estimation/<timestamp>/estimation-record.json`
   (exactly what step 5 printed to stdout — do not edit or reformat it).

7. **Report to the user**, framed explicitly as a prediction, e.g.:

   > Estimated delivery effort for `<feature-directory>`: **X.X DEU**
   > (confidence Y%, risk: <label>). This is a prediction based on the
   > current spec/plan/tasks, not a commitment. Full derivation and record
   > saved under `<feature-directory>/estimation/<timestamp>/`.

## Never do this

- Never overwrite or delete an existing `estimation/<timestamp>/` directory
  for this feature — a re-run always creates a new one (FR-007).
- Never state a Delivery Effort, Confidence, Risk, or any other engine
  output number that didn't come verbatim from `estimatorctl`'s stdout.
- Never skip the justification for a dimension "to save time" — an
  unjustified score defeats the entire point of User Story 2.
