# Setup Reference

Use this reference when the user asks to install or update Canary Compact project instructions.

Run the setup script from the current project root:

    ${CLAUDE_SKILL_DIR}/scripts/setup-canary-instructions.sh [canary-word] [files...]

When called from `/canary-compact:canary-compact setup ...`, strip the literal `setup` subcommand before invoking the script.

If no canary word is passed, the script uses `CANARY_COMPACT_WORD` from the environment, then falls back to `CANARY_COMPACT`.

If no files are passed, the script updates:

    CLAUDE.md
    AGENTS.md

The script is idempotent:

- It creates missing target files.
- It appends a managed Canary Compact block when none exists.
- It replaces the existing managed block when both markers are present.
- It errors without rewriting when markers are mismatched or duplicated.
- It treats marker lines as exact-line matches so incidental inline mentions do not trigger replacement.

Examples:

    ${CLAUDE_SKILL_DIR}/scripts/setup-canary-instructions.sh
    ${CLAUDE_SKILL_DIR}/scripts/setup-canary-instructions.sh CUSTOM_CANARY_WORD
    ${CLAUDE_SKILL_DIR}/scripts/setup-canary-instructions.sh CUSTOM_CANARY_WORD CLAUDE.md AGENTS.md "docs/agent notes.md"

After running the script, report the updated files and the canary word used.
