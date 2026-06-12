---
name: canary-compact
description: Use when configuring, testing, or explaining canary-triggered session compaction behavior.
argument-hint: "[setup [canary-word] [files...]]"
allowed-tools: Bash
---

# Canary Compact

Use the canary word to request compaction from an assistant reply.

Default canary word:

```text
CANARY_COMPACT
```

When working in Claude Code, this plugin's `Stop` hook checks the final assistant reply. If the reply contains the canary word, the hook keeps the conversation alive and asks Claude to produce a compact-ready handoff plus the exact `/compact` command the user should submit.

Important limitation: current Claude Code hooks can inspect replies and continue the conversation, but they do not expose a hook action that directly invokes `/compact`. In Codex, this plugin exposes the shared skill and helper source, but this environment's Codex plugin manifest does not accept hook wiring.

Configuration:

- `CANARY_COMPACT_WORD`: override the canary word.
- `CANARY_COMPACT_CASE_SENSITIVE=false`: make matching case-insensitive.
- `CANARY_COMPACT_WHOLE_WORD=true`: require token-boundary matching.

To test the detector manually, build the release binaries and pipe a Stop-hook-shaped JSON payload to `bin/canary-compact-hook`.

## Setup

When the user asks to install, configure, or update the project-level canary instruction, read `references/setup.md` and run the bundled setup script from the current project root.

If the skill invocation begins with `setup`, treat `setup` as a subcommand and pass only the remaining arguments to `scripts/setup-canary-instructions.sh`. Do not pass the literal word `setup` as the canary word.

Common invocations:

    /canary-compact:canary-compact setup
    /canary-compact:canary-compact setup CUSTOM_CANARY_WORD
    /canary-compact:canary-compact setup CUSTOM_CANARY_WORD CLAUDE.md AGENTS.md docs/agents.md
