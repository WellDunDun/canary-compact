# Canary Compact

Canary Compact helps long Claude Code sessions end cleanly before they run out of useful context.

When an assistant reply contains a canary word, the plugin's Claude Code hook keeps the session alive and asks Claude to produce a compact-ready handoff summary plus the exact `/compact` command for you to submit.

Default canary word:

```text
CANARY_COMPACT
```

## What it does

- Watches Claude Code assistant replies at the `Stop` hook event.
- Detects a configurable canary word in the final assistant message.
- Prompts Claude to summarize current task state, decisions, changed files, tests, blockers, and next steps.
- Shows the `/compact` command to run with that handoff as focus instructions.
- Adds a setup skill that can install the canary instruction into `CLAUDE.md` and `AGENTS.md`.

The plugin does not run `/compact` automatically. Claude Code hooks can continue a conversation, but they do not currently expose an API for invoking the built-in `/compact` command directly.

## Set up a project

After installing the plugin, run the setup skill inside Claude Code:

```text
/canary-compact:canary-compact setup
```

That command adds an idempotent managed block to:

- `CLAUDE.md`
- `AGENTS.md`

Use a custom canary word when you want one:

```text
/canary-compact:canary-compact setup SESSION_NEEDS_COMPACT
```

Target specific instruction files when needed:

```text
/canary-compact:canary-compact setup SESSION_NEEDS_COMPACT CLAUDE.md AGENTS.md "docs/agent notes.md"
```

The setup script creates missing files, replaces its own managed block on later runs, and refuses to rewrite files with broken or duplicated markers.

## Use it

Work normally. When the assistant decides the session should be compacted, it should include the canary word in its reply.

When the hook sees the canary word, Claude Code receives extra context telling it to stop unrelated work and produce a handoff. Submit the `/compact` command it gives you, then continue from the compacted session.

## Configuration

Set these environment variables before starting Claude Code when you want to change matching behavior:

| Variable | Default | Meaning |
| --- | --- | --- |
| `CANARY_COMPACT_WORD` | `CANARY_COMPACT` | Canary word to detect. |
| `CANARY_COMPACT_CASE_SENSITIVE` | `true` | Set to `false` for case-insensitive matching. |
| `CANARY_COMPACT_WHOLE_WORD` | `false` | Set to `true` to require token-boundary matching. |

Detections are appended to `canary-detections.jsonl` under Claude Code's plugin data directory when `CLAUDE_PLUGIN_DATA` is available.

## Runtime requirements

No Python, Node, Go, or jq install is required for normal use. The plugin ships compiled helper binaries and uses a small POSIX shell launcher.

## Support matrix

| Host | Platform | Status | Notes |
| --- | --- | --- | --- |
| Claude Code | macOS arm64 | Supported | Automatic `Stop` hook. |
| Claude Code | macOS amd64 | Supported | Automatic `Stop` hook. |
| Claude Code | Linux arm64 | Supported | Automatic `Stop` hook. |
| Claude Code | Linux amd64 | Supported | Automatic `Stop` hook. |
| Claude Code | Windows amd64 | Experimental | Binary is included, but the current hook manifest uses POSIX `sh`. |
| Codex | All platforms | Skill-only | Codex plugin manifests in this environment do not support automatic hook wiring. |

## Troubleshooting

If the hook does not trigger:

1. Make sure Claude Code is running with the plugin installed or loaded with `--plugin-dir`.
2. Confirm the assistant reply contains the exact canary word.
3. Check whether you changed `CANARY_COMPACT_CASE_SENSITIVE` or `CANARY_COMPACT_WHOLE_WORD`.
4. Run `/canary-compact:canary-compact setup` again if the project instruction block is missing or stale.

If you are testing locally after editing the plugin, run `/reload-plugins` or restart Claude Code.

## License

MIT
