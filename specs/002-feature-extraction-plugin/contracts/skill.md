# Contract: `estimate-feature` skill

Located at `plugin/skills/estimate-feature/SKILL.md`. Invocable two ways
(FR-009): manually (`/estimate-feature <feature-directory>`), or
automatically via the hook's injected `additionalContext` (contracts/hook.md),
which names the feature directory to use.

## Preconditions

- `<feature-directory>/spec.md` exists (required).
- `<feature-directory>/plan.md` and `<feature-directory>/tasks.md` exist
  where possible; if either is missing, the skill proceeds best-effort per
  the research.md rubric ("missing/thin source → conservative score +
  assumption note"), it does not abort.

## Steps (contract, not prose — see SKILL.md for the full rubric text)

1. Read whichever of `spec.md`/`plan.md`/`tasks.md` exist in
   `<feature-directory>`.
2. Score the 7 `EstimationFeatures` dimensions per the research.md anchor
   table, each with a justification string.
3. Run `plugin/scripts/ensure-engine.sh` (builds/reuses `plugin/bin/estimatorctl`;
   see contracts/ensure-engine.md). If it reports the Go toolchain is
   missing, stop here and tell the user exactly that (FR-012) — do not
   proceed to a partial estimation.
4. Create `<feature-directory>/estimation/<UTC-timestamp>/` and write
   `derived-features.json` (data-model.md `DerivedFeatureVector` shape).
5. Run:
   ```
   plugin/bin/estimatorctl estimate \
     --work-item <feature-directory basename> \
     --features <feature-directory>/estimation/<UTC-timestamp>/derived-features.json
   ```
   (per specs/001/contracts/cli.md — `estimate` reads only the `features`
   object's keys from a JSON file, so `derived-features.json`'s extra
   `justifications`/`assumptions`/etc. fields are harmless extras the CLI
   ignores.)
6. Save the CLI's stdout verbatim as
   `<feature-directory>/estimation/<UTC-timestamp>/estimation-record.json`.
7. Report a concise summary to the user in the current turn (FR-015),
   explicitly framed as a prediction (Constitution Principle VI), e.g.:
   "Estimated delivery effort for `<feature-directory>`: X.X DEU
   (confidence Y%, risk: <label>). Full derivation and record saved under
   `<feature-directory>/estimation/<timestamp>/`."

## Postconditions

- A new `estimation/<timestamp>/` directory exists; no existing file under
  `<feature-directory>` was modified or deleted (FR-007).
- Every dimension in `derived-features.json`'s `features` has a matching
  `justifications` entry (data-model.md validation rule).
- If step 3 failed (no Go toolchain), no `estimation/<timestamp>/` directory
  is created at all — a failed engine build must not produce a
  half-written, misleading result.
