// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2020- Leonardo Di Donato <leodidonato@gmail.com>
package main

import (
	"bytes"
	"io"
	"testing"
)

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
		// The parser uses Ragel's 'i flag on type tokens, so uppercase
		// types parse successfully. The labeler in
		// .github/workflows/label.yml MUST stay in lockstep with this
		// (it does, via `shopt -s nocasematch`). Codifying both
		// directions of the contract here.
		{name: "uppercase type", header: "FEAT: add foo", wantErr: false},
		{name: "mixed-case type", header: "Feat(API)!: x", wantErr: false},
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
		{name: "malformed UTF-8", header: string(append([]byte("feat: invalid "), 0xff)), wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validate(tc.header, false, io.Discard)
			if tc.wantErr && err == nil {
				t.Fatalf("validate(%q) returned nil error; want non-nil", tc.header)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("validate(%q) returned %v; want nil", tc.header, err)
			}
		})
	}
}

// TestDescribe pins the canonical-form contract that the
// .github/workflows/label.yml labeler relies on. The labeler reads
// these key=value lines and applies labels from them; if this output
// shape ever changes, the labeler breaks.
func TestDescribe(t *testing.T) {
	cases := []struct {
		name   string
		header string
		want   string
	}{
		{
			name:   "plain",
			header: "feat: add foo",
			want:   "type=feat\nbreaking=false\n",
		},
		{
			name:   "scope",
			header: "feat(parser): add foo",
			want:   "type=feat\nscope=parser\nbreaking=false\n",
		},
		{
			name:   "bang implies breaking",
			header: "feat!: x",
			want:   "type=feat\nbreaking=true\n",
		},
		{
			name:   "scope and bang",
			header: "fix(api)!: x",
			want:   "type=fix\nscope=api\nbreaking=true\n",
		},
		{
			// Type and scope come back lowercased: the parser uses
			// Ragel's 'i flag on both tokens, so "FEAT(API)" and
			// "feat(api)" reduce to the same canonical form before we
			// even see them. ccvalidate ALSO lowercases type
			// explicitly to harden this against a future Ragel
			// regeneration that drops 'i. The labeler relies on this
			// canonical form.
			name:   "type and scope canonicalised to lowercase",
			header: "FEAT(API)!: x",
			want:   "type=feat\nscope=api\nbreaking=true\n",
		},
		{
			// Style and revert have no dedicated release-notes
			// category — the labeler buckets them under "chore". The
			// canonical type still surfaces here so consumers can
			// route as they like.
			name:   "style preserved",
			header: "style: tabs",
			want:   "type=style\nbreaking=false\n",
		},
		{
			name:   "revert preserved",
			header: "revert: undo",
			want:   "type=revert\nbreaking=false\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := validate(tc.header, true, &buf); err != nil {
				t.Fatalf("validate(%q, describe=true) returned %v; want nil", tc.header, err)
			}
			if got := buf.String(); got != tc.want {
				t.Fatalf("describe output mismatch for %q:\n got: %q\nwant: %q", tc.header, got, tc.want)
			}
		})
	}
}

// TestDescribeOnInvalidEmitsNothing ensures we don't write a partial
// or stale describe payload when the header fails to parse. Callers
// (the labeler script) MUST be able to rely on "exit non-zero =>
// nothing on stdout to act on".
func TestDescribeOnInvalidEmitsNothing(t *testing.T) {
	var buf bytes.Buffer
	err := validate(`Revert "feat: x"`, true, &buf)
	if err == nil {
		t.Fatal("expected error for non-CC header, got nil")
	}
	if got := buf.String(); got != "" {
		t.Fatalf("describe stdout on invalid input should be empty, got %q", got)
	}
}

// TestParseFlags pins flag handling so a description starting with a
// dash isn't mistakenly treated as a flag.
func TestParseFlags(t *testing.T) {
	cases := []struct {
		name     string
		args     []string
		wantArgs []string
		wantDesc bool
	}{
		{name: "no flags", args: []string{"feat: x"}, wantArgs: []string{"feat: x"}, wantDesc: false},
		{name: "describe flag", args: []string{"--describe", "feat: x"}, wantArgs: []string{"feat: x"}, wantDesc: true},
		{name: "dash-prefixed description not a flag", args: []string{"feat: -x added"}, wantArgs: []string{"feat: -x added"}, wantDesc: false},
		{name: "empty", args: []string{}, wantArgs: []string{}, wantDesc: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotArgs, gotDesc := parseFlags(tc.args)
			if gotDesc != tc.wantDesc {
				t.Fatalf("parseFlags(%v) describe = %v; want %v", tc.args, gotDesc, tc.wantDesc)
			}
			if len(gotArgs) != len(tc.wantArgs) {
				t.Fatalf("parseFlags(%v) args = %v; want %v", tc.args, gotArgs, tc.wantArgs)
			}
			for i := range gotArgs {
				if gotArgs[i] != tc.wantArgs[i] {
					t.Fatalf("parseFlags(%v) args[%d] = %q; want %q", tc.args, i, gotArgs[i], tc.wantArgs[i])
				}
			}
		})
	}
}
