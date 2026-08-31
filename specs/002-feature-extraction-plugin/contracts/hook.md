# Contract: `on-tasks-md-change.sh` (PostToolUse hook)

Declared in `plugin/hooks/hooks.json`, two entries (one per matcher, per
research.md — matcher does not combine tool names reliably for our purposes):

```json
{
  "hooks": {
    "PostToolUse": [
      { "matcher": "Write", "hooks": [ { "type": "command", "command": "${CLAUDE_PLUGIN_ROOT}/scripts/on-tasks-md-change.sh" } ] },
      { "matcher": "Edit",  "hooks": [ { "type": "command", "command": "${CLAUDE_PLUGIN_ROOT}/scripts/on-tasks-md-change.sh" } ] }
    ]
  }
}
```

## Input (stdin)

JSON with at least:

```json
{
  "hook_event_name": "PostToolUse",
  "tool_name": "Write",
  "tool_input": { "file_path": "specs/002-feature-extraction-plugin/tasks.md" },
  "cwd": "/path/to/project"
}
```

`on-tasks-md-change.sh` reads `tool_input.file_path` only; every other field
is ignored.

## Behavior

1. Extract `file_path` from stdin (grep/sed, research.md — no `jq` dependency).
2. If `file_path` does not end in `/tasks.md`, exit 0 with no output (not our
   event).
3. Let `dir` = the directory containing `file_path`. If `dir/spec.md` or
   `dir/plan.md` does not exist, exit 0 with no output (Edge Case: stray
   `tasks.md`, FR-014).
4. Otherwise, print to stdout:

   ```json
   {
     "hookSpecificOutput": {
       "hookEventName": "PostToolUse",
       "additionalContext": "A Spec-Kit feature just completed planning: <dir>. Run the estimate-feature skill for this feature directory now."
     }
   }
   ```

5. Exit 0.

## Output contract

- Exit code `0` always (detection failures are not errors — SC-005: this
  hook must never surface as a broken workflow to the user).
- stdout is either empty or exactly one JSON object shaped as above. No
  other output format is emitted.
- The hook performs no writes, no network calls, and does not itself invoke
  `estimatorctl` — it only decides *whether* to prompt Claude to run the
  `estimate-feature` skill (see `skill.md` contract).
