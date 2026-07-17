#!/usr/bin/env bash

set -euo pipefail

if [[ $# -ne 4 ]]; then
  echo "usage: classify-labels.sh <ccvalidate> <message-file> <changed-files-json> <expected-file-count>" >&2
  exit 2
fi

validator=$1
message_file=$2
changed_files_file=$3
expected_file_count=$4

[[ -x "$validator" ]] || { echo "validator is not executable: $validator" >&2; exit 2; }
[[ -r "$message_file" ]] || { echo "message file is not readable: $message_file" >&2; exit 2; }
[[ -r "$changed_files_file" ]] || { echo "changed-files file is not readable: $changed_files_file" >&2; exit 2; }
[[ "$expected_file_count" =~ ^(0|[1-9][0-9]*)$ ]] || {
  echo "expected file count is not a non-negative integer: $expected_file_count" >&2
  exit 2
}

# `gh api --paginate --slurp` produces an array of response-page arrays.
# Validate that envelope before inspecting paths so a truncated response or a
# future API/CLI shape change fails closed instead of silently dropping labels.
if ! actual_file_count=$(jq -er -s '
  if length == 1 and
    (.[0] |
      type == "array" and
      length > 0 and
      all(.[];
        type == "array" and
        all(.[]; type == "object" and (.filename | type == "string"))
      )
    )
  then
    (.[0] | [.[][]] | length)
  else
    error("invalid changed-files JSON envelope")
  end
' "$changed_files_file"); then
  echo "invalid changed-files JSON envelope" >&2
  exit 1
fi
if [[ "$actual_file_count" != "$expected_file_count" ]]; then
  echo "incomplete changed-files response: got $actual_file_count of $expected_file_count files" >&2
  exit 1
fi

desired=()
commit_type=""
breaking=false

if describe=$(CC_INPUT_FILE="$message_file" "$validator" --describe 2>/dev/null); then
  seen_type=false
  seen_scope=false
  seen_breaking=false

  while IFS='=' read -r key value; do
    case "$key" in
      type)
        [[ "$seen_type" == false ]] || { echo "duplicate type from ccvalidate" >&2; exit 1; }
        seen_type=true
        commit_type=$value
        ;;
      scope)
        [[ "$seen_scope" == false ]] || { echo "duplicate scope from ccvalidate" >&2; exit 1; }
        seen_scope=true
        ;;
      breaking)
        [[ "$seen_breaking" == false ]] || { echo "duplicate breaking flag from ccvalidate" >&2; exit 1; }
        seen_breaking=true
        breaking=$value
        ;;
      *)
        echo "unexpected field from ccvalidate: $key" >&2
        exit 1
        ;;
    esac
  done <<<"$describe"

  if [[ "$seen_type" != true || "$seen_breaking" != true ]]; then
    echo "incomplete output from ccvalidate" >&2
    exit 1
  fi
else
  echo "::warning title=Non-Conventional-Commits PR title::The PR title is not a valid Conventional Commit and will not receive a category label." >&2
fi

case "$commit_type" in
  feat) desired+=(enhancement) ;;
  fix) desired+=(bug) ;;
  perf) desired+=(performance) ;;
  refactor) desired+=(refactor) ;;
  docs) desired+=(documentation) ;;
  test) desired+=(test) ;;
  build) desired+=(build) ;;
  ci) desired+=(ci) ;;
  chore|style|revert) desired+=(chore) ;;
  "") ;;
  *)
    echo "unexpected type from ccvalidate: $commit_type" >&2
    exit 1
    ;;
esac

case "$breaking" in
  true) desired+=(breaking-change) ;;
  false) ;;
  *)
    echo "unexpected breaking value from ccvalidate: $breaking" >&2
    exit 1
    ;;
esac

# Filenames stay structured JSON values. In particular, embedded newlines or
# shell metacharacters can never be reinterpreted as paths or commands.
if jq -e 'any(.[][]; .filename == "go.mod" or .filename == "go.sum")' "$changed_files_file" >/dev/null; then
  desired+=(dependencies)
fi
if jq -e 'any(.[][]; .filename | startswith("parser/"))' "$changed_files_file" >/dev/null; then
  desired+=(parser)
fi
if jq -e 'any(.[][]; .filename | startswith("tools/"))' "$changed_files_file" >/dev/null; then
  desired+=(tools)
fi

allowed='^(enhancement|bug|performance|refactor|documentation|test|build|ci|chore|breaking-change|dependencies|parser|tools)$'
for label in "${desired[@]:-}"; do
  [[ -z "$label" || "$label" =~ $allowed ]] || {
    echo "refusing non-whitelisted label: $label" >&2
    exit 1
  }
done

printf '%s\n' "${desired[@]:-}" | jq -Rsc 'split("\n") | map(select(length > 0))'
