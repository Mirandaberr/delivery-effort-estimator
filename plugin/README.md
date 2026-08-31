# Delivery Effort Estimator — Claude Code Plugin

Auto-derives a Delivery Effort prediction for a Spec-Kit feature right after
it's planned, by reading its `spec.md`/`plan.md`/`tasks.md` and calling the
independent, deterministic `estimatorctl` engine — no hand-authored features
JSON needed. See `specs/002-feature-extraction-plugin/` in the main repo for
the full spec, plan, and design rationale.

## Install

This directory is self-contained — copy it anywhere, no other part of the
`delivery-effort-estimator` repository is required.

```bash
claude --plugin-dir /path/to/plugin
```

**Prerequisite**: a Go toolchain on `PATH`. Nothing else — the engine binary
is built automatically on first use (see below).

## What happens on first use

The first time the `estimate-feature` skill runs (automatically, after a
Spec-Kit feature's `tasks.md` is written, or manually), it builds
`plugin/bin/estimatorctl` from the bundled `plugin/engine/` source. This
takes a few seconds once; subsequent runs reuse the built binary.

## How it works

1. A `PostToolUse` hook (`hooks/hooks.json`) detects when `tasks.md` is
   created/edited inside a real `specs/<id>-*/` feature directory (one with
   sibling `spec.md`/`plan.md`) and tells Claude to run the
   `estimate-feature` skill for it.
2. The skill (`skills/estimate-feature/SKILL.md`) reads that feature's
   `spec.md`/`plan.md`/`tasks.md`, scores 7 normalized dimensions against a
   documented rubric, and writes the derivation with a justification for
   every score.
3. It then runs `bin/estimatorctl estimate` (built from `engine/`) to get
   the actual prediction — the LLM never computes the number itself.
4. Both the derivation and the prediction are saved under that feature's
   own `specs/<id>-*/estimation/<timestamp>/` directory. Nothing is ever
   overwritten; re-running adds a new timestamped result.

You can also run the skill manually on any feature directory instead of
waiting for the automatic trigger.

## Directory layout

```text
plugin/
├── .claude-plugin/plugin.json   manifest
├── hooks/hooks.json             PostToolUse trigger (detection only)
├── skills/estimate-feature/     the actual derivation logic (LLM reasoning)
├── scripts/                     locate-feature.sh, on-tasks-md-change.sh, ensure-engine.sh
├── engine/                      synced copy of the estimatorctl Go source (never hand-edited)
├── bin/                         built estimatorctl binary (git-ignored, built on first use)
└── tests/                       shell tests for the deterministic scripts
```

`engine/` is kept in sync with the main repo's `cmd/`/`internal/` via
`../scripts/sync-plugin-engine.sh` (`make sync-plugin-engine` from the repo
root) — if you're developing the plugin itself, run that after any engine
change, or `make test` will fail the drift check.
