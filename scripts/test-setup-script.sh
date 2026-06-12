#!/usr/bin/env sh
set -eu

root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
setup="$root/plugins/canary-compact/skills/canary-compact/scripts/setup-canary-instructions.sh"
tmp_dir=$(mktemp -d "${TMPDIR:-/tmp}/canary-setup-test.XXXXXX")

cd "$tmp_dir"

"$setup" CUSTOM_CANARY "CLAUDE.md" "my docs/notes.md"

test -f "CLAUDE.md"
test -f "my docs/notes.md"
test ! -e "my"

grep -q "CUSTOM_CANARY" "CLAUDE.md"
grep -q "CUSTOM_CANARY" "my docs/notes.md"

"$setup" OTHER_CANARY "CLAUDE.md" "my docs/notes.md"

start_count=$(grep -xF "<!-- canary-compact:start -->" "CLAUDE.md" | wc -l | tr -d '[:space:]')
end_count=$(grep -xF "<!-- canary-compact:end -->" "CLAUDE.md" | wc -l | tr -d '[:space:]')
test "$start_count" = "1"
test "$end_count" = "1"
grep -q "OTHER_CANARY" "CLAUDE.md"
if grep -q "CUSTOM_CANARY" "CLAUDE.md"; then
  echo "old canary word was not replaced" >&2
  exit 1
fi

printf "%s\n" "<!-- canary-compact:start -->" "tail content" > broken.md
if "$setup" CANARY broken.md 2>broken.err; then
  echo "setup succeeded despite a missing end marker" >&2
  exit 1
fi
grep -q "mismatched Canary Compact markers" broken.err
grep -q "tail content" broken.md

printf "%s\n" "This line mentions <!-- canary-compact:start --> inline." > mention.md
"$setup" INLINE_CANARY mention.md
grep -q "This line mentions" mention.md
start_count=$(grep -xF "<!-- canary-compact:start -->" mention.md | wc -l | tr -d '[:space:]')
test "$start_count" = "1"

echo "setup script tests passed in $tmp_dir"
