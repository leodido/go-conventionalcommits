// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2026- Leonardo Di Donato <leodidonato@gmail.com>

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const baselineWorkflow = `name: action-pins fixture
on: workflow_dispatch
jobs:
  baseline:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@9c091bb21b7c1c1d1991bb908d89e4e9dddfe3e0 # v7.0.0
      - uses: actions/setup-go@b7ad1dad31e06c5925ef5d2fc7ad053ef454303e # v7.0.0
`

func TestCheckWorkflows(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		workflow string
		wantErr  string
	}{
		{
			name: "accepts reviewed baseline",
		},
		{
			name: "accepts local action",
			workflow: `jobs:
  local:
    runs-on: ubuntu-latest
    steps:
      - uses: ./local-action
`,
		},
		{
			name: "rejects Docker action",
			workflow: `jobs:
  docker:
    runs-on: ubuntu-latest
    steps:
      - uses: docker://example.invalid/action:latest
`,
			wantErr: "Docker actions are not allowed",
		},
		{
			name: "rejects quoted Docker uses key",
			workflow: `jobs:
  docker:
    runs-on: ubuntu-latest
    steps:
      - "uses": docker://example.invalid/action:latest
`,
			wantErr: "Docker actions are not allowed",
		},
		{
			name: "rejects escaped Docker uses key",
			workflow: `jobs:
  docker:
    runs-on: ubuntu-latest
    steps:
      - "\u0075ses": docker://example.invalid/action:latest
`,
			wantErr: "Docker actions are not allowed",
		},
		{
			name: "rejects Docker action alias",
			workflow: `env:
  docker-action: &docker-action docker://example.invalid/action:latest
jobs:
  docker:
    runs-on: ubuntu-latest
    steps:
      - uses: *docker-action
`,
			wantErr: "Docker actions are not allowed",
		},
		{
			name: "rejects Docker action alias key",
			workflow: `env:
  uses-key: &uses-key uses
jobs:
  docker:
    runs-on: ubuntu-latest
    steps:
      - *uses-key: docker://example.invalid/action:latest
`,
			wantErr: "Docker actions are not allowed",
		},
		{
			name: "rejects unsafe checkout key",
			workflow: `jobs:
  unsafe:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@9c091bb21b7c1c1d1991bb908d89e4e9dddfe3e0 # v7.0.0
        with:
          "allow-unsafe-pr-checkout": true
`,
			wantErr: "allow-unsafe-pr-checkout must not be configured",
		},
		{
			name: "rejects escaped unsafe checkout key",
			workflow: `jobs:
  unsafe:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@9c091bb21b7c1c1d1991bb908d89e4e9dddfe3e0 # v7.0.0
        with:
          "allow-unsafe-pr-check\u006fut": true
`,
			wantErr: "allow-unsafe-pr-checkout must not be configured",
		},
		{
			name: "rejects aliased unsafe checkout key",
			workflow: `env:
  unsafe-key: &unsafe-key allow-unsafe-pr-checkout
jobs:
  unsafe:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@9c091bb21b7c1c1d1991bb908d89e4e9dddfe3e0 # v7.0.0
        with:
          *unsafe-key: true
`,
			wantErr: "allow-unsafe-pr-checkout must not be configured",
		},
		{
			name: "rejects unpinned reusable workflow",
			workflow: `jobs:
  reusable:
    uses: owner/repository/.github/workflows/reusable.yml@main
`,
			wantErr: "external action is not pinned to a full commit",
		},
		{
			name: "rejects quoted unpinned uses key",
			workflow: `jobs:
  unpinned:
    runs-on: ubuntu-latest
    steps:
      - "uses": actions/checkout@main
`,
			wantErr: "external action is not pinned to a full commit",
		},
		{
			name: "requires checkout version annotation",
			workflow: `jobs:
  checkout:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@9c091bb21b7c1c1d1991bb908d89e4e9dddfe3e0
`,
			wantErr: "expected # v7.0.0 annotation",
		},
		{
			name: "requires setup go version annotation",
			workflow: `jobs:
  setup-go:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/setup-go@b7ad1dad31e06c5925ef5d2fc7ad053ef454303e
`,
			wantErr: "expected # v7.0.0 annotation",
		},
		{
			name: "requires reviewed setup go commit",
			workflow: `jobs:
  setup-go:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/setup-go@aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa # v7.0.0
`,
			wantErr: "expected actions/setup-go@b7ad1dad31e06c5925ef5d2fc7ad053ef454303e",
		},
		{
			name: "rejects mixed case checkout reference",
			workflow: `jobs:
  checkout:
    runs-on: ubuntu-latest
    steps:
      - uses: Actions/checkout@aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa # v7.0.0
`,
			wantErr: "expected actions/checkout@9c091bb21b7c1c1d1991bb908d89e4e9dddfe3e0",
		},
		{
			name: "ignores env data named uses",
			workflow: `env:
  uses: not-an-action
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - run: true
`,
		},
		{
			name: "allows non-checkout input named allow unsafe checkout",
			workflow: `jobs:
  custom-action:
    runs-on: ubuntu-latest
    steps:
      - uses: owner/repository@aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
        with:
          allow-unsafe-pr-checkout: true
`,
		},
		{
			name: "rejects multiple YAML documents",
			workflow: `jobs:
  first:
    runs-on: ubuntu-latest
    steps:
      - run: true
---
name: second action-pins fixture
on: workflow_dispatch
jobs:
  second:
    runs-on: ubuntu-latest
    steps:
      - uses: owner/repository@main
`,
			wantErr: "workflow must contain exactly one YAML document",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			workflows := t.TempDir()
			writeWorkflow(t, workflows, "baseline.yml", baselineWorkflow)
			if tt.workflow != "" {
				writeWorkflow(t, workflows, "fixture.yml", fixtureWorkflow(tt.workflow))
			}

			err := checkWorkflows(workflows)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("checkWorkflows() error = %v", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("checkWorkflows() succeeded, want error containing %q", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("checkWorkflows() error = %q, want substring %q", err, tt.wantErr)
			}
		})
	}
}

func fixtureWorkflow(body string) string {
	return "name: action-pins fixture\non: workflow_dispatch\n" + body
}

func writeWorkflow(t *testing.T, dir, name, content string) {
	t.Helper()

	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
