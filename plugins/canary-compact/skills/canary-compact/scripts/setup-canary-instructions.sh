#!/usr/bin/env sh
set -eu

start_marker="<!-- canary-compact:start -->"
end_marker="<!-- canary-compact:end -->"
default_word="${CANARY_COMPACT_WORD:-CANARY_COMPACT}"

canary_word="${1:-$default_word}"
if [ "$#" -gt 0 ]; then
  shift
fi

if [ "$#" -eq 0 ]; then
  set -- "CLAUDE.md" "AGENTS.md"
fi

tmp_root="${TMPDIR:-/tmp}"

write_block() {
  cat <<EOF
$start_marker
## Canary Compact

When the active assistant session should be compacted, end the assistant reply with the exact canary word:

    $canary_word

Use the canary only for intentional compaction, such as when context is getting tight, a handoff would preserve important state, or the current task is about to cross a safe continuation boundary. Do not include the canary in ordinary explanations, code blocks, examples, or status updates unless compaction is intended.

The Canary Compact plugin watches final assistant replies in Claude Code. When it sees the canary, it asks Claude to prepare a compact-ready handoff and the exact /compact command for the user to submit.
$end_marker
EOF
}

replace_or_append_block() {
  file="$1"
  dir=$(dirname -- "$file")

  if [ "$dir" != "." ] && [ ! -d "$dir" ]; then
    mkdir -p "$dir"
  fi

  if [ -f "$file" ]; then
    start_count=$(grep -xF "$start_marker" "$file" | wc -l | tr -d '[:space:]')
    end_count=$(grep -xF "$end_marker" "$file" | wc -l | tr -d '[:space:]')
    if [ "$start_count" != "$end_count" ]; then
      printf "error: %s has mismatched Canary Compact markers: start=%s end=%s\n" "$file" "$start_count" "$end_count" >&2
      return 1
    fi
    if [ "$start_count" -gt 1 ]; then
      printf "error: %s has multiple Canary Compact managed blocks\n" "$file" >&2
      return 1
    fi
  else
    start_count=0
  fi

  block_file=$(mktemp "$tmp_root/canary-compact-block.XXXXXX")
  tmp_file=$(mktemp "$tmp_root/canary-compact-file.XXXXXX")
  write_block > "$block_file"

  if [ "$start_count" -eq 1 ]; then
    awk -v start="$start_marker" -v end="$end_marker" -v block_file="$block_file" '
      BEGIN {
        while ((getline line < block_file) > 0) {
          block = block line ORS
        }
        in_block = 0
      }
      $0 == start {
        printf "%s", block
        in_block = 1
        next
      }
      $0 == end {
        in_block = 0
        next
      }
      !in_block {
        print
      }
    ' "$file" > "$tmp_file"
  else
    if [ -f "$file" ] && [ -s "$file" ]; then
      cat "$file" > "$tmp_file"
      printf "\n\n" >> "$tmp_file"
    fi
    cat "$block_file" >> "$tmp_file"
  fi

  mv "$tmp_file" "$file"
  rm -f "$block_file"
  printf "updated %s with canary word %s\n" "$file" "$canary_word"
}

for file do
  replace_or_append_block "$file"
done
