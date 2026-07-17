#!/usr/bin/env bash

set -euo pipefail

root=$(git rev-parse --show-toplevel)
workflow="$root/.github/workflows/label.yml"
classifier="$root/.github/scripts/classify-labels.sh"
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

fail() {
  printf 'FAIL: %s\n' "$*" >&2
  exit 1
}

assert_contains() {
  local haystack=$1
  local needle=$2
  local context=$3
  [[ "$haystack" == *"$needle"* ]] || fail "$context: missing $needle"
}

assert_not_contains() {
  local haystack=$1
  local needle=$2
  local context=$3
  [[ "$haystack" != *"$needle"* ]] || fail "$context: unexpectedly contains $needle"
}

assert_eq() {
  local got=$1
  local want=$2
  local context=$3
  if [[ "$got" != "$want" ]]; then
    printf 'FAIL: %s\n  got:  %q\n  want: %q\n' "$context" "$got" "$want" >&2
    exit 1
  fi
}

job_block() {
  local job=$1
  awk -v job="$job" '
    $0 == "  " job ":" { found = 1 }
    found && $0 ~ /^  [A-Za-z0-9_-]+:$/ && $0 != "  " job ":" { exit }
    found { print }
  ' "$workflow"
}

extract_apply_script() {
  awk '
    /^      - name: Apply labels$/ { step = 1; next }
    step && /^        run: \|$/ { script = 1; next }
    script && /^          / { sub(/^          /, ""); print; next }
    script && /^[[:space:]]*$/ { print; next }
    script { exit }
  ' "$workflow"
}

# The pull_request_target workflow may read and classify untrusted metadata,
# but only the second job may mutate labels. That job must be too small to
# execute repository content, actions, builds, or PR-controlled fields.
classify_job=$(job_block classify)
apply_job=$(job_block apply)
[[ -n "$classify_job" ]] || fail "missing classify job"
[[ -n "$apply_job" ]] || fail "missing apply job"

assert_contains "$classify_job" "contents: read" "classify permissions"
assert_contains "$classify_job" "pull-requests: read" "classify permissions"
assert_not_contains "$classify_job" ": write" "classify permissions"
# shellcheck disable=SC2016 # Assert literal Actions expressions.
assert_contains "$classify_job" 'ref: ${{ github.event.pull_request.base.sha }}' "trusted checkout"
assert_contains "$classify_job" "cache: false" "privileged-trigger cache policy"
assert_contains "$classify_job" "EXPECTED_CHANGED_FILES" "changed-files completeness check"
assert_contains "$classify_job" "gh api --paginate --slurp" "changed-files pagination"
# shellcheck disable=SC2016 # Assert literal Actions expressions.
assert_contains "$classify_job" 'labels: ${{ steps.classify.outputs.labels }}' "classifier output"

assert_contains "$apply_job" "issues: write" "apply permissions"
apply_write_permissions=$(awk '/^[[:space:]]+[a-z-]+: write$/ { count++ } END { print count + 0 }' <<<"$apply_job")
assert_eq "$apply_write_permissions" "1" "apply has exactly one write permission"
assert_not_contains "$apply_job" "contents:" "apply permissions"
assert_not_contains "$apply_job" "pull-requests:" "apply permissions"
assert_not_contains "$apply_job" "uses:" "apply execution boundary"
assert_not_contains "$apply_job" "actions/" "apply execution boundary"
assert_not_contains "$apply_job" "checkout" "apply execution boundary"
assert_not_contains "$apply_job" "go build" "apply execution boundary"
assert_not_contains "$apply_job" "PR_TITLE" "apply data boundary"
assert_not_contains "$apply_job" "PR_BODY" "apply data boundary"
assert_not_contains "$apply_job" "pull_request.head" "apply data boundary"
# shellcheck disable=SC2016 # Assert a literal Actions expression.
assert_contains "$apply_job" '${{ needs.classify.outputs.labels }}' "apply data boundary"

workflow_text=$(<"$workflow")
assert_not_contains "$workflow_text" "actions/labeler" "external labeler removal"
[[ ! -e "$root/.github/labeler.yml" ]] || fail ".github/labeler.yml is obsolete"
while IFS= read -r use_line; do
  action_ref=${use_line#*@}
  action_ref=${action_ref%% *}
  [[ "$action_ref" =~ ^[0-9a-f]{40}$ ]] || fail "action is not pinned to a full commit: $use_line"
done < <(grep -E '^[[:space:]]+- uses:' "$workflow")

# Build the same trusted validator used by the classifier job.
GOCACHE="$tmp/go-build" go build -o "$tmp/ccvalidate" ./tools/ccvalidate

message="$tmp/message"
files="$tmp/files.json"

printf 'feat(parser)!: isolate label writes\n' >"$message"
printf '%s\n' '[[{"filename":"go.mod"},{"filename":"parser/machine.go.rl"}],[{"filename":"tools/ccvalidate/main.go"}]]' >"$files"
got=$("$classifier" "$tmp/ccvalidate" "$message" "$files" 3)
assert_eq "$got" '["enhancement","breaking-change","dependencies","parser","tools"]' "combined title and path classification"

while read -r commit_type expected_label; do
  printf '%s: classify this type\n' "$commit_type" >"$message"
  printf '%s\n' '[[]]' >"$files"
  got=$("$classifier" "$tmp/ccvalidate" "$message" "$files" 0)
  assert_eq "$got" "[\"$expected_label\"]" "$commit_type category classification"
done <<'EOF'
feat enhancement
fix bug
perf performance
refactor refactor
docs documentation
test test
build build
ci ci
chore chore
style chore
revert chore
EOF

printf 'fix: classify the message\n\nBREAKING CHANGE: reject unsafe output\n' >"$message"
printf '%s\n' '[[]]' >"$files"
got=$("$classifier" "$tmp/ccvalidate" "$message" "$files" 0)
assert_eq "$got" '["bug","breaking-change"]' "breaking footer classification"

printf 'not a conventional title\n' >"$message"
printf '%s\n' '[[{"filename":"parser/utf8_test.go"}]]' >"$files"
got=$("$classifier" "$tmp/ccvalidate" "$message" "$files" 1 2>"$tmp/warning")
assert_eq "$got" '["parser"]' "invalid title retains path labels"
assert_contains "$(<"$tmp/warning")" "Non-Conventional-Commits" "invalid-title warning"

# Filenames remain JSON strings throughout classification. A newline cannot
# turn the suffix of one filename into a second matching path.
printf '%s\n' '[[{"filename":"docs/readme.md\nparser/injected.go"}]]' >"$files"
got=$("$classifier" "$tmp/ccvalidate" "$message" "$files" 1 2>/dev/null)
assert_eq "$got" '[]' "newline-bearing filename"

printf '%s\n' '{"filename":"parser/not-a-page.json"}' >"$files"
if "$classifier" "$tmp/ccvalidate" "$message" "$files" 1 >/dev/null 2>&1; then
  fail "classifier accepted an invalid changed-files envelope"
fi

printf '%s\n' '[]' >"$files"
if "$classifier" "$tmp/ccvalidate" "$message" "$files" 0 >/dev/null 2>&1; then
  fail "classifier accepted a slurp envelope with no response page"
fi

printf '%s\n%s\n' '[[]]' '[[{"filename":"parser/injected.go"}]]' >"$files"
if "$classifier" "$tmp/ccvalidate" "$message" "$files" 1 >/dev/null 2>&1; then
  fail "classifier accepted multiple top-level JSON documents"
fi

printf '%s\n' '[[{"filename":"parser/only-one.go"}]]' >"$files"
if "$classifier" "$tmp/ccvalidate" "$message" "$files" 3001 >/dev/null 2>&1; then
  fail "classifier accepted a truncated changed-files response"
fi

cat >"$tmp/untrusted-validator" <<'EOF'
#!/usr/bin/env bash
printf '%s\n' 'type=feat' 'breaking=false' 'label=owned'
EOF
chmod +x "$tmp/untrusted-validator"
printf '%s\n' '[[]]' >"$files"
if "$classifier" "$tmp/untrusted-validator" "$message" "$files" 0 >/dev/null 2>&1; then
  fail "classifier accepted an unexpected validator field"
fi

# Execute the exact inline script from the write-only job against a fake gh
# binary. This pins label-specific POST/DELETE behavior without granting the
# test credentials or replacing the full label set.
apply_script="$tmp/apply-labels.sh"
extract_apply_script >"$apply_script"
[[ -s "$apply_script" ]] || fail "could not extract apply script"

mkdir "$tmp/bin"
cat >"$tmp/bin/gh" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
printf '%q ' "$@" >>"$GH_LOG"
printf '\n' >>"$GH_LOG"

if [[ " $* " == *" -X POST "* ]]; then
  exit 0
fi
if [[ " $* " == *" -X DELETE "* ]]; then
  if [[ ${FAKE_DELETE_FAILURE:-false} == true ]]; then
    touch "$FAKE_STATE/delete-attempted"
    exit 1
  fi
  exit 0
fi
if [[ ${FAKE_READ_FAIL:-false} == true ]]; then
  exit 1
fi
if [[ " $* " != *" --paginate "* || " $* " != *" --slurp "* ]]; then
  echo "GET must use --paginate --slurp" >&2
  exit 2
fi
labels=${FAKE_CURRENT:-}
if [[ -e "$FAKE_STATE/delete-attempted" ]]; then
  labels=${FAKE_AFTER_DELETE_CURRENT:-$labels}
fi
case ${FAKE_RESPONSE_MODE:-valid} in
  object) printf '%s\n' '{"not":"paginated"}'; exit 0 ;;
  empty) printf '%s\n' '[]'; exit 0 ;;
  multiple) printf '%s\n%s\n' '[]' '[[{"name":"bug"}]]'; exit 0 ;;
esac
jq -cn --arg labels "$labels" '
  [
    $labels
    | split(",")
    | map(select(length > 0) | {name: .})
  ]
'
EOF
chmod +x "$tmp/bin/gh"

export PATH="$tmp/bin:$PATH"
export GH_LOG="$tmp/gh.log"
export FAKE_STATE="$tmp/fake-state"
mkdir "$FAKE_STATE"
export GH_REPO="owner/repo"
export PR_NUMBER="123"
export DESIRED_LABELS_JSON='["enhancement","parser"]'
export FAKE_CURRENT="bug,parser,manual"
: >"$GH_LOG"
bash "$apply_script" >/dev/null
calls=$(<"$GH_LOG")
assert_contains "$calls" "-X POST" "label addition"
assert_contains "$calls" "labels\[\]=enhancement" "label addition"
assert_not_contains "$calls" "labels\[\]=parser" "already-present label"
assert_contains "$calls" "-X DELETE" "stale-label removal"
assert_contains "$calls" "/labels/bug" "stale-label removal"
assert_not_contains "$calls" "/labels/manual" "unmanaged-label preservation"
assert_not_contains "$calls" "-X PUT" "whole-set replacement ban"

injection_marker="$tmp/injection-marker"
# shellcheck disable=SC2016 # The literal command substitution is attack data.
attack_label='$(touch '"$injection_marker"')'
DESIRED_LABELS_JSON=$(jq -cn --arg attack "$attack_label" '["enhancement", $attack]')
export DESIRED_LABELS_JSON
: >"$GH_LOG"
if bash "$apply_script" >/dev/null 2>&1; then
  fail "apply job accepted a non-whitelisted label"
fi
[[ ! -e "$injection_marker" ]] || fail "apply job executed classifier output as shell code"
[[ ! -s "$GH_LOG" ]] || fail "apply job reached GitHub before rejecting its input"

export DESIRED_LABELS_JSON='["enhancement\n"]'
: >"$GH_LOG"
if bash "$apply_script" >/dev/null 2>&1; then
  fail "apply job accepted a label with a trailing newline"
fi
[[ ! -s "$GH_LOG" ]] || fail "apply job reached GitHub before rejecting a non-exact label"

export DESIRED_LABELS_JSON=$'["no-releasenotes"]\n["enhancement"]'
: >"$GH_LOG"
if bash "$apply_script" >/dev/null 2>&1; then
  fail "apply job accepted multiple top-level classifier documents"
fi
[[ ! -s "$GH_LOG" ]] || fail "apply job reached GitHub before rejecting multiple classifier documents"

export DESIRED_LABELS_JSON='[]'
export FAKE_READ_FAIL=true
: >"$GH_LOG"
if bash "$apply_script" >/dev/null 2>&1; then
  fail "apply job continued after failing to read current labels"
fi
assert_not_contains "$(<"$GH_LOG")" "-X POST" "read failure"
assert_not_contains "$(<"$GH_LOG")" "-X DELETE" "read failure"

export FAKE_READ_FAIL=false
export FAKE_DELETE_FAILURE=true
export FAKE_CURRENT="bug,manual"
export FAKE_AFTER_DELETE_CURRENT="manual"
rm -f "$FAKE_STATE/delete-attempted"
: >"$GH_LOG"
bash "$apply_script" >/dev/null
assert_contains "$(<"$GH_LOG")" "/labels/bug" "idempotent deletion"

export FAKE_AFTER_DELETE_CURRENT="bug,manual"
rm -f "$FAKE_STATE/delete-attempted"
: >"$GH_LOG"
if bash "$apply_script" >/dev/null 2>&1; then
  fail "apply job ignored a failed deletion that left the label present"
fi

export FAKE_DELETE_FAILURE=false
export FAKE_RESPONSE_MODE=object
rm -f "$FAKE_STATE/delete-attempted"
: >"$GH_LOG"
if bash "$apply_script" >/dev/null 2>&1; then
  fail "apply job accepted an invalid current-label response"
fi
assert_not_contains "$(<"$GH_LOG")" "-X POST" "invalid current-label response"
assert_not_contains "$(<"$GH_LOG")" "-X DELETE" "invalid current-label response"

for response_mode in empty multiple; do
  export FAKE_RESPONSE_MODE=$response_mode
  : >"$GH_LOG"
  if bash "$apply_script" >/dev/null 2>&1; then
    fail "apply job accepted $response_mode current-label JSON"
  fi
  assert_not_contains "$(<"$GH_LOG")" "-X POST" "$response_mode current-label response"
  assert_not_contains "$(<"$GH_LOG")" "-X DELETE" "$response_mode current-label response"
done

printf 'label workflow tests: PASS\n'
