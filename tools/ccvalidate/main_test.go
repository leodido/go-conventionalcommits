// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2020- Leonardo Di Donato <leodidonato@gmail.com>
package main

import "testing"

func TestValidate(t *testing.T) {
	cases := []struct {
		name    string
		header  string
		wantErr bool
	}{
		// ---- valid Conventional Commits ----
		{name: "feat", header: "feat: add foo", wantErr: false},
		{name: "feat with scope", header: "feat(parser): add foo", wantErr: false},
		{name: "feat bang", header: "feat!: breaking", wantErr: false},
		{name: "feat scope bang", header: "feat(parser)!: breaking", wantErr: false},
		{name: "fix", header: "fix: bar", wantErr: false},
		{name: "perf", header: "perf: speed up", wantErr: false},
		{name: "refactor", header: "refactor(api): rename", wantErr: false},
		{name: "docs", header: "docs: typo", wantErr: false},
		{name: "test", header: "test: add cases", wantErr: false},
		{name: "build", header: "build(deps): bump", wantErr: false},
		{name: "ci", header: "ci: tweak workflow", wantErr: false},
		{name: "chore", header: "chore: cleanup", wantErr: false},
		{name: "style", header: "style: format", wantErr: false},
		{name: "revert", header: "revert: undo", wantErr: false},
		{name: "shell metacharacters in description are fine", header: "feat: add $(rm -rf /tmp/x)", wantErr: false},
		{name: "trailing LF accepted", header: "feat: x\n", wantErr: false},
		{name: "trailing CRLF accepted", header: "feat: x\r\n", wantErr: false},

		// ---- invalid ----
		{name: "non-CC", header: "Update README", wantErr: true},
		{name: "GitHub default revert", header: `Revert "feat: x"`, wantErr: true},
		{name: "empty", header: "", wantErr: true},
		{name: "only newline", header: "\n", wantErr: true},
		{name: "type only with colon", header: "feat:", wantErr: true},
		{name: "type with empty description", header: "feat: ", wantErr: true},
		{name: "space before colon", header: "feat : x", wantErr: true},
		{name: "unknown type", header: "wip: draft", wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validate(tc.header)
			if tc.wantErr && err == nil {
				t.Fatalf("validate(%q) returned nil error; want non-nil", tc.header)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("validate(%q) returned %v; want nil", tc.header, err)
			}
		})
	}
}
