# Canary Compact

Detects a canary word in assistant replies and prompts session compaction.

## Contents

- `.claude-plugin/plugin.json`: Claude Code plugin manifest.
- `.codex-plugin/plugin.json`: Codex plugin manifest.
- `bin/canary-compact-hook`: Runtime launcher for the compiled helper.
- `hooks/hooks.json`: Claude Code hook wiring.
- `LICENSE`: MIT license for the installable plugin package.
- `skills/canary-compact/SKILL.md`: Shared usage skill.
- `skills/canary-compact/references/setup.md`: Setup instructions for project files.
- `skills/canary-compact/scripts/setup-canary-instructions.sh`: Idempotent setup script.

The shared `skills/` directory and compiled helper keep behavior in one place while each host reads its own manifest.

## Behavior

The default canary word is `CANARY_COMPACT`.

In Claude Code, the plugin registers a `Stop` hook. The hook reads `last_assistant_message`, checks for the canary word, and returns `additionalContext` asking Claude to produce a compact-ready handoff and a `/compact` command.

Claude Code hooks can inspect replies and continue the conversation, but they do not expose a hook action that directly invokes the built-in `/compact` slash command. This plugin therefore gets as close as the current hook API allows without recursively spawning Claude.

## Runtime packaging

The hook does not call `python3`, `node`, or `jq`. It calls `sh bin/canary-compact-hook`, which selects a compiled Go helper from `bin/<os>-<arch>/`.

Build release binaries from the repository root:

```sh
sh scripts/build-plugin-binaries.sh
```

The build script writes:

```text
plugins/canary-compact/bin/darwin-amd64/canary-compact-hook
plugins/canary-compact/bin/darwin-arm64/canary-compact-hook
plugins/canary-compact/bin/linux-amd64/canary-compact-hook
plugins/canary-compact/bin/linux-arm64/canary-compact-hook
plugins/canary-compact/bin/windows-amd64/canary-compact-hook.exe
```

## Support matrix

| Host | Platform | Status | Runtime path | Notes |
| --- | --- | --- | --- | --- |
| Claude Code | macOS arm64 | Supported | `bin/darwin-arm64/canary-compact-hook` | Automatic `Stop` hook via `sh bin/canary-compact-hook`. |
| Claude Code | macOS amd64 | Supported | `bin/darwin-amd64/canary-compact-hook` | Automatic `Stop` hook via `sh bin/canary-compact-hook`. |
| Claude Code | Linux arm64 | Supported | `bin/linux-arm64/canary-compact-hook` | Automatic `Stop` hook via `sh bin/canary-compact-hook`. |
| Claude Code | Linux amd64 | Supported | `bin/linux-amd64/canary-compact-hook` | Automatic `Stop` hook via `sh bin/canary-compact-hook`. |
| Claude Code | Windows amd64 | Experimental | `bin/windows-amd64/canary-compact-hook.exe` | Binary is built, but the current hook manifest uses POSIX `sh`; publish a Windows-specific hook variant before claiming full support. |
| Codex | All platforms | Skill-only | N/A | Codex plugin manifests in this environment do not support automatic hook wiring. |

## Claude Code local testing

Load the plugin directly during development:

```sh
claude --plugin-dir ./plugins/canary-compact
```

Then run the skill in Claude Code:

```text
/canary-compact:canary-compact
```

Install the project-level canary instruction:

```text
/canary-compact:canary-compact setup
```

The setup flow updates `CLAUDE.md` and `AGENTS.md` with an idempotent managed block. Pass a custom word and optional target files when needed:

```text
/canary-compact:canary-compact setup CUSTOM_CANARY_WORD CLAUDE.md AGENTS.md
```

You can also install it through this repository's local Claude marketplace:

```text
/plugin marketplace add .
/plugin install canary-compact@canary-compact-plugins
```

## Codex local testing

This repo includes a Codex marketplace entry at `.agents/plugins/marketplace.json`.

Add the local marketplace from the repository root, then install the plugin:

```sh
codex plugin marketplace add .
codex plugin add canary-compact@canary-compact-plugins
```

Start a new Codex thread after installing or reinstalling so new plugin skills are picked up.

## Extending the skill

Add references or scripts under the existing `skills/canary-compact/` directory unless you are introducing a genuinely separate user-facing workflow:

```text
skills/canary-compact/references/<topic>.md
skills/canary-compact/scripts/<helper>.sh
```

Claude Code invokes the main skill as `/canary-compact:canary-compact`. Codex discovers it through the `skills` path in `.codex-plugin/plugin.json`.
