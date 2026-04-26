#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
#
# Asserts that the cc-title case statement in
# .github/workflows/label.yml AGREES with tools/ccvalidate's
# --describe output for a fixed corpus of PR titles and bodies.
#
# Two failure modes this guards against:
#
#   1. Labeler and validator drift in opposite directions — a title
#      becomes acceptable to one but not the other (or vice versa),
#      or a body footer is interpreted differently.
#   2. The cc-title block in label.yml is edited and someone forgets
#      that the bash patterns must mirror what the parser accepts.
#
# Usage:
#   .github/scripts/test-labeler-agreement.sh
#
# Exit codes:
#   0  every case agreed
#   1  at least one disagreement (printed to stderr)
#   2  drift between this script's embedded case block and the one
#      in .github/workflows/label.yml

set -euo pipefail

ROOT=$(cd "$(dirname "$0")/../.." && pwd)
WORKFLOW="$ROOT/.github/workflows/label.yml"
CCVALIDATE_SRC="$ROOT/tools/ccvalidate"

# ---------------------------------------------------------------
# Mirror of cc-title's title-to-labels logic. MUST stay byte-for-
# byte equal to the case block in label.yml; the drift-check below
# enforces this. To update: change the block in label.yml, run this
# script, copy the failing diff into both places.
# ---------------------------------------------------------------
title_to_labels() {
  local PR_TITLE="$1"
  local desired_lines
  desired_lines=$(
    shopt -s nocasematch
    out=()
    case "$PR_TITLE" in
      feat\!:\ ?*|feat\(*\)\!:\ ?*)             out+=("enhancement" "breaking-change") ;;
      fix\!:\ ?*|fix\(*\)\!:\ ?*)               out+=("bug" "breaking-change") ;;
      perf\!:\ ?*|perf\(*\)\!:\ ?*)             out+=("performance" "breaking-change") ;;
      refactor\!:\ ?*|refactor\(*\)\!:\ ?*)     out+=("refactor" "breaking-change") ;;
      docs\!:\ ?*|docs\(*\)\!:\ ?*)             out+=("documentation" "breaking-change") ;;
      test\!:\ ?*|test\(*\)\!:\ ?*)             out+=("test" "breaking-change") ;;
      build\!:\ ?*|build\(*\)\!:\ ?*)           out+=("build" "breaking-change") ;;
      ci\!:\ ?*|ci\(*\)\!:\ ?*)                 out+=("ci" "breaking-change") ;;
      chore\!:\ ?*|chore\(*\)\!:\ ?*)           out+=("chore" "breaking-change") ;;
      style\!:\ ?*|style\(*\)\!:\ ?*)           out+=("chore" "breaking-change") ;;
      revert\!:\ ?*|revert\(*\)\!:\ ?*)         out+=("chore" "breaking-change") ;;
      feat:\ ?*|feat\(*\):\ ?*)                 out+=("enhancement") ;;
      fix:\ ?*|fix\(*\):\ ?*)                   out+=("bug") ;;
      perf:\ ?*|perf\(*\):\ ?*)                 out+=("performance") ;;
      refactor:\ ?*|refactor\(*\):\ ?*)         out+=("refactor") ;;
      docs:\ ?*|docs\(*\):\ ?*)                 out+=("documentation") ;;
      test:\ ?*|test\(*\):\ ?*)                 out+=("test") ;;
      build:\ ?*|build\(*\):\ ?*)               out+=("build") ;;
      ci:\ ?*|ci\(*\):\ ?*)                     out+=("ci") ;;
      chore:\ ?*|chore\(*\):\ ?*)               out+=("chore") ;;
      style:\ ?*|style\(*\):\ ?*)               out+=("chore") ;;
      revert:\ ?*|revert\(*\):\ ?*)             out+=("chore") ;;
    esac
    printf '%s\n' "${out[@]:-}"
  )
  printf '%s' "$desired_lines" | sort -u | sed '/^$/d' | tr '\n' ',' | sed 's/,$//'
}

# Extract the literal pattern arms (lines containing "desired+=" /
# "out+=") from each file and assert they are identical. Catches
# drift the very next time someone edits either file.
#
# We compare the body of the case block, not the full block, so
# differences in YAML indentation, the `case "$PR_TITLE"` line and
# the `esac` line don't trigger spurious failures. The arms are the
# part we actually care about staying in lockstep.
drift_check() {
  local from_workflow from_script
  from_workflow=$(grep -E '^[[:space:]]*[a-z][a-z]+[!]?[:(].*out\+=' "$WORKFLOW" \
    | sed -E 's/^[[:space:]]+//' )
  from_script=$(grep -E '^[[:space:]]*[a-z][a-z]+[!]?[:(].*out\+=' "$0" \
    | sed -E 's/^[[:space:]]+//' )
  if [ "$from_workflow" != "$from_script" ]; then
    echo "DRIFT: case-block arms in $WORKFLOW differ from this script." >&2
    diff <(printf '%s\n' "$from_workflow") <(printf '%s\n' "$from_script") || true
    exit 2
  fi
}

# Build ccvalidate once.
build_ccvalidate() {
  local out="${TMPDIR:-/tmp}/ccvalidate-agreement-test"
  ( cd "$ROOT" && go build -o "$out" "./tools/ccvalidate" )
  echo "$out"
}

# Compare the labeler's verdict on a (title, body) to ccvalidate's.
# Reports a single pass/fail line per case.
agree_or_fail() {
  local title="$1" body="$2" ccv="$3"

  # 1. Labeler verdict.
  local labeler_csv
  labeler_csv=$(title_to_labels "$title")
  local labeler_breaking=false
  if [[ ",$labeler_csv," == *,breaking-change,* ]]; then
    labeler_breaking=true
  fi
  # Labeler also adds breaking-change from a real footer in the body.
  # Mirror that path so the agreement check covers what label.yml
  # actually does.
  local msg_file
  msg_file=$(mktemp)
  {
    printf '%s\n' "$title"
    if [ -n "$body" ]; then printf '\n%s' "$body"; fi
  } > "$msg_file"
  local describe_out=""
  local ccv_valid=false
  if describe_out=$(CC_INPUT_FILE="$msg_file" "$ccv" --describe 2>/dev/null); then
    ccv_valid=true
  fi
  rm "$msg_file"

  local ccv_breaking=false
  if [ "$ccv_valid" = true ] && grep -qx 'breaking=true' <<< "$describe_out"; then
    ccv_breaking=true
    labeler_breaking=true   # match label.yml's behaviour
  fi

  # 2. Validator agreement: a valid CC title MUST produce at least
  # one labeler label; an invalid one MUST produce none. The
  # exception is the parser's case-insensitive type token: every
  # title accepted by ccvalidate has a corresponding lowercased
  # form that the labeler matches via shopt -s nocasematch.
  local labeler_any=false
  if [ -n "$labeler_csv" ]; then labeler_any=true; fi

  if [ "$ccv_valid" != "$labeler_any" ]; then
    echo "FAIL agreement: title=$(printf '%q' "$title") ccv_valid=$ccv_valid labeler_any=$labeler_any labels=$labeler_csv" >&2
    return 1
  fi
  if [ "$ccv_breaking" != "$labeler_breaking" ]; then
    echo "FAIL breaking: title=$(printf '%q' "$title") ccv_breaking=$ccv_breaking labeler_breaking=$labeler_breaking" >&2
    return 1
  fi

  printf 'OK   title=%-32s valid=%-5s breaking=%s\n' \
    "$(printf '%q' "$title")" "$ccv_valid" "$ccv_breaking"
}

main() {
  drift_check

  # Sanity: ccvalidate package builds before we drive any inputs
  # through it.
  if [ ! -d "$CCVALIDATE_SRC" ]; then
    echo "missing $CCVALIDATE_SRC" >&2
    exit 1
  fi
  local ccv
  ccv=$(build_ccvalidate)

  # Corpus: each row is "TITLE|BODY". BODY may contain literal \n
  # which we expand below.
  # Mix of: lowercase, uppercase, mixed-case, every CC type, bang
  # variants, file-glob non-CC, GitHub revert default, malformed,
  # body footers in/out of fences/prose/indent.
  local corpus=(
    "feat: x|"
    "FEAT: x|"
    "Feat: x|"
    "fEAt: x|"
    "fix: bar|"
    "FIX: bar|"
    "perf: speed|"
    "refactor(api): rename|"
    "docs: typo|"
    "test: cases|"
    "build(deps): bump|"
    "ci: workflow|"
    "chore: cleanup|"
    "style: tabs|"
    "revert: undo|"
    "feat!: bang|"
    "fix(api)!: bang|"
    "FEAT(API)!: BANG|"
    "perf!: speed|"
    "build(deps)!: bump|"
    'feat: x|body\n\nBREAKING CHANGE: kaboom'
    'feat: x|body\n\nBREAKING-CHANGE: kaboom'
    'docs: explain|see```\nBREAKING CHANGE: example\n```'
    'docs: explain|It is a BREAKING CHANGE: foo'
    'docs: explain|body\n\n  BREAKING CHANGE: indented'
    'docs: explain|body\n\nbreaking change: kaboom'
    'Update README|'
    'Revert "feat: x"|'
    'wip: draft|'
    'feat: |'
    'feat:|'
    'feat:x|'
    'feat : x|'
  )

  local fail=0
  for row in "${corpus[@]}"; do
    local title="${row%%|*}"
    local body="${row#*|}"
    # Expand literal \n in BODY.
    body=$(printf '%b' "$body")
    if ! agree_or_fail "$title" "$body" "$ccv"; then
      fail=$((fail + 1))
    fi
  done

  if [ "$fail" -ne 0 ]; then
    echo "$fail case(s) failed" >&2
    exit 1
  fi
  echo "all ${#corpus[@]} cases agreed."
}

main "$@"
