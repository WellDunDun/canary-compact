# Canary Compact Plugin Scaffold

This repository is a dual-host plugin scaffold for Claude Code and Codex. The plugin source lives in `plugins/canary-compact/`.

## Structure

```text
.
├── .agents/plugins/marketplace.json
├── .claude-plugin/marketplace.json
├── plugins/canary-compact/
│   ├── .claude-plugin/plugin.json
│   ├── .codex-plugin/plugin.json
│   ├── bin/canary-compact-hook
│   ├── hooks/hooks.json
│   ├── LICENSE
│   ├── README.md
│   ├── skills/canary-compact/SKILL.md
│   ├── skills/canary-compact/references/setup.md
│   └── skills/canary-compact/scripts/setup-canary-instructions.sh
├── cmd/canary-compact-hook/
├── scripts/build-plugin-binaries.sh
└── LICENSE
```

## Claude Code

For direct development:

```sh
claude --plugin-dir ./plugins/canary-compact
```

For marketplace installation:

```text
/plugin marketplace add .
/plugin install canary-compact@canary-compact-plugins
```

After installing, add the project-level canary instruction:

```text
/canary-compact:canary-compact setup
```

## Codex

For marketplace installation from the repository root:

```sh
codex plugin marketplace add .
codex plugin add canary-compact@canary-compact-plugins
```

After changing plugin metadata or skills, reinstall the plugin and start a new Codex thread.

## Canary compact behavior

The default canary word is `CANARY_COMPACT`. In Claude Code, the plugin hook checks the final assistant reply at the `Stop` lifecycle event. When the canary appears, the hook keeps the session alive and asks Claude to produce a compact-ready handoff plus a `/compact` command for the user to submit.

Current Claude Code hooks can inspect replies and continue the conversation, but they do not expose an API that directly invokes the built-in `/compact` slash command. Codex plugin manifests in this environment do not accept hooks, so Codex gets the shared skill and compiled helper source but not automatic hook wiring.

## Public distribution notes

The hook runtime path does not require Python, Node, or jq. It calls a small compiled Go helper through `plugins/canary-compact/bin/canary-compact-hook`.

Before publishing a release to a public marketplace, build the platform binaries:

```sh
sh scripts/build-plugin-binaries.sh
```

This writes binaries for macOS and Linux, plus a Windows `.exe`, under `plugins/canary-compact/bin/`. The committed plugin source should include those release binaries because Claude Code installs plugins by copying the plugin directory from the marketplace source.

## Support matrix

| Host | Platform | Status | Runtime path | Notes |
| --- | --- | --- | --- | --- |
| Claude Code | macOS arm64 | Supported | `bin/darwin-arm64/canary-compact-hook` | Automatic `Stop` hook via `sh bin/canary-compact-hook`. |
| Claude Code | macOS amd64 | Supported | `bin/darwin-amd64/canary-compact-hook` | Automatic `Stop` hook via `sh bin/canary-compact-hook`. |
| Claude Code | Linux arm64 | Supported | `bin/linux-arm64/canary-compact-hook` | Automatic `Stop` hook via `sh bin/canary-compact-hook`. |
| Claude Code | Linux amd64 | Supported | `bin/linux-amd64/canary-compact-hook` | Automatic `Stop` hook via `sh bin/canary-compact-hook`. |
| Claude Code | Windows amd64 | Experimental | `bin/windows-amd64/canary-compact-hook.exe` | Binary is built, but the current hook manifest uses POSIX `sh`; publish a Windows-specific hook variant before claiming full support. |
| Codex | All platforms | Skill-only | N/A | Codex plugin manifests in this environment do not support automatic hook wiring. |

## References used for this scaffold

- Claude Code plugin docs: https://code.claude.com/docs/en/plugins
- Claude Code plugin reference: https://code.claude.com/docs/en/plugins-reference
- Claude Code marketplace docs: https://code.claude.com/docs/en/plugin-marketplaces
- Codex local plugin creator reference bundled with this Codex environment.
