#!/bin/bash
# Runs every example .r2 script with a wall-clock timeout, to catch
# parser/runtime regressions without hanging CI on long-running server
# demos. Exits non-zero if any script NOT listed in
# scripts/known_flaky_examples.txt fails or times out — that's treated as a
# real regression. A script listed there failing is expected and does not
# break the build (see that file for why each entry is there).
#
# Usage: scripts/sweep_examples.sh [path-to-r2-binary] [timeout-seconds]
#   Builds ./r2lang-sweep-bin from main.go if no binary path is given.
#
# Deliberately avoids bash 4+ features (associative arrays, etc): macOS
# ships bash 3.2 as /bin/bash and this script needs to run unmodified both
# there and on CI's newer bash.
set -u

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT" || exit 1

BIN="${1:-}"
TIMEOUT_SECS="${2:-20}"
CLEANUP_BIN=0

if [ -z "$BIN" ]; then
  BIN="$ROOT/r2lang-sweep-bin"
  echo "Building $BIN from main.go..."
  if ! go build -o "$BIN" main.go; then
    echo "Build failed, aborting sweep."
    exit 1
  fi
  CLEANUP_BIN=1
fi

KNOWN_FLAKY_FILE="$ROOT/scripts/known_flaky_examples.txt"
# Normalized (comments/blank lines stripped, whitespace-trimmed) copy of the
# known-flaky list, one path per line, used via `grep -Fxq` lookups below.
KNOWN_FLAKY_NORMALIZED=$(mktemp)
if [ -f "$KNOWN_FLAKY_FILE" ]; then
  sed -e 's/#.*//' "$KNOWN_FLAKY_FILE" | sed -e 's/[[:space:]]*$//' -e 's/^[[:space:]]*//' | grep -v '^$' > "$KNOWN_FLAKY_NORMALIZED"
fi

is_known_flaky() {
  grep -Fxq "$1" "$KNOWN_FLAKY_NORMALIZED"
}

pass=0
fail=0
timeout_count=0
regressions_file=$(mktemp)
now_passing_flaky_file=$(mktemp)

while IFS= read -r f; do
  rel="${f#./}"
  out_log=$(mktemp)
  "$BIN" "$f" > "$out_log" 2>&1 &
  pid=$!
  ( sleep "$TIMEOUT_SECS"; kill -9 "$pid" 2>/dev/null ) &
  watcher=$!
  wait "$pid" 2>/dev/null
  exitcode=$?
  kill "$watcher" 2>/dev/null
  wait "$watcher" 2>/dev/null

  if [ "$exitcode" -eq 0 ]; then
    pass=$((pass+1))
    if is_known_flaky "$rel"; then
      echo "$rel" >> "$now_passing_flaky_file"
    fi
  elif [ "$exitcode" -eq 137 ]; then
    timeout_count=$((timeout_count+1))
    echo "TIMEOUT: $rel"
    if ! is_known_flaky "$rel"; then
      echo "$rel :: timeout after ${TIMEOUT_SECS}s" >> "$regressions_file"
    fi
  else
    fail=$((fail+1))
    reason=$(grep -m1 -E "Parser Exception|panic:|Undeclared variable|Error:" "$out_log")
    echo "FAIL($exitcode): $rel :: $reason"
    if ! is_known_flaky "$rel"; then
      echo "$rel :: $reason" >> "$regressions_file"
    fi
  fi
  rm -f "$out_log"
done < <(find . -name "*.r2" \
  -not -path "./examples/dsl/*" \
  -not -path "./examples/proyecto/contable/old_versions/*" \
  -not -path "./examples/proyecto/contable/versiones_anteriores/*")

echo "=== pass=$pass fail=$fail timeout=$timeout_count ==="

if [ "$CLEANUP_BIN" -eq 1 ]; then
  rm -f "$BIN"
fi

if [ -s "$now_passing_flaky_file" ]; then
  echo ""
  echo "NOTE: the following scripts are listed in known_flaky_examples.txt but passed this run — consider removing them from the list:"
  sed 's/^/  - /' "$now_passing_flaky_file"
fi

exit_code=0
if [ -s "$regressions_file" ]; then
  echo ""
  echo "REGRESSIONS (not in known_flaky_examples.txt):"
  sed 's/^/  - /' "$regressions_file"
  exit_code=1
fi

rm -f "$KNOWN_FLAKY_NORMALIZED" "$regressions_file" "$now_passing_flaky_file"
exit "$exit_code"
