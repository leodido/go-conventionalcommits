// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2020- Leonardo Di Donato <leodidonato@gmail.com>
package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindTrailerBlockStart(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "empty",
			input: "",
			want:  -1,
		},
		{
			name:  "trailing-newlines-only",
			input: "\n\n\n",
			want:  -1,
		},
		{
			name:  "description-only",
			input: "fix: x",
			want:  -1,
		},
		{
			name:  "description-then-trailing-newlines",
			input: "fix: x\n\n",
			want:  -1,
		},
		{
			name:  "single-trailer",
			input: "fix: x\n\nFixes: 3",
			want:  len("fix: x\n\n"),
		},
		{
			name:  "multiple-trailers-no-blank-gap",
			input: "fix: x\n\nFixes #3\nSigned-off-by: Leo",
			want:  len("fix: x\n\n"),
		},
		{
			name:  "trailer-block-with-internal-blank-gap",
			input: "fix: x\n\nAcked-by: Leo\n\nBREAKING CHANGE: APIs",
			want:  len("fix: x\n\n"),
		},
		{
			name:  "trailer-block-after-many-blank-lines",
			input: "fix: x\n\n\n\n\nFixes #3\nSigned-off-by: Leo",
			want:  len("fix: x\n\n\n\n\n"),
		},
		{
			name:  "fake-trailer-then-body-paragraph",
			input: "feat: x\n\nFixes #15\n\nLorem ipsum dolor sit amet",
			want:  -1,
		},
		{
			name:  "body-line-ending-trailer-shape-then-real-trailer",
			input: "feat: x\n\nbody1\nFake: trail\n\nReviewed-by: X",
			want:  len("feat: x\n\nbody1\nFake: trail\n\n"),
		},
		{
			name:  "multi-paragraph-body-each-ending-trailer-shape-then-real-trailer",
			input: "feat: x\n\npara1\nFake: a\n\npara2\nFake: b\n\nReviewed-by: X",
			want:  len("feat: x\n\npara1\nFake: a\n\npara2\nFake: b\n\n"),
		},
		{
			name:  "malformed-trailer-line-inside-block",
			input: "fix: x\n\nTested-by: Leo\nX-\nAnother-trailer: x",
			want:  len("fix: x\n\n"),
		},
		{
			name:  "trailer-block-with-trailing-newline",
			input: "fix: x\n\nFixes: 3\n",
			want:  len("fix: x\n\n"),
		},
		{
			name:  "breaking-change-special-case",
			input: "fix: x\n\nBREAKING CHANGE: APIs",
			want:  len("fix: x\n\n"),
		},
		{
			name:  "breaking-change-without-space-after-colon",
			input: "fix: x\n\nBREAKING CHANGE:APIs",
			want:  -1,
		},
		{
			name:  "trailer-with-hash-token",
			input: "fix: x\n\nCloses #42",
			want:  len("fix: x\n\n"),
		},
		{
			name: "trailer-with-tab-after-colon-not-trailer-shaped",
			// `ws` in common.rl is a single ASCII space, so a tab
			// after the colon disqualifies the line; this also
			// mirrors the FSM's behavior on develop.
			input: "fix: x\n\nFixes:\t3",
			want:  -1,
		},
		{
			name: "first-line-trailer-shaped-no-blank-line-no-trailer-block",
			// No blank line means no separator means no trailer block.
			input: "Fixes: 3",
			want:  -1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := findTrailerBlockStart([]byte(tc.input))
			assert.Equal(t, tc.want, got, "input=%q", tc.input)
		})
	}
}

func TestIsTrailerStartLine(t *testing.T) {
	cases := []struct {
		name string
		line string
		want bool
	}{
		{"empty", "", false},
		{"plain-token-then-colon-space", "Fixes: 3", true},
		{"plain-token-then-colon-no-space", "Fixes:3", false},
		{"plain-token-then-colon-tab", "Fixes:\t3", false},
		{"plain-token-then-space-hash", "Closes #3", true},
		{"plain-token-then-space-no-hash", "Closes 3", false},
		{"hyphenated-token", "Signed-off-by: Leo", true},
		{"hyphenated-token-trailing-dash", "X-: foo", false},
		{"hyphenated-token-leading-dash", "-X: foo", false},
		{"breaking-change-with-space", "BREAKING CHANGE: APIs", true},
		{"breaking-change-no-space", "BREAKING CHANGE:APIs", false},
		{"breaking-change-prefix-different-case", "breaking change: APIs", false},
		{"non-alnum-start", " Fixes: 3", false},
		{"alnum-only-no-sep", "Fixes", false},
		{"underscore-not-alnum-by-ragel", "snake_case: x", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := isTrailerStartLine([]byte(tc.line))
			assert.Equal(t, tc.want, got, "line=%q", tc.line)
		})
	}
}

func TestFindTrailerBlockStartLargeInput(t *testing.T) {
	// Sanity check: the implementation should be linear-time. We don't
	// time-bound the test (CI noise would flake it) but we DO want to
	// exercise the perf path with a many-paragraph body so any
	// quadratic regression shows up loud in benchmarks.
	var b strings.Builder
	b.WriteString("feat: x\n\n")
	for i := 0; i < 1000; i++ {
		b.WriteString("paragraph line one\n")
		b.WriteString("paragraph line two: looks-like-trailer\n")
		b.WriteString("\n")
	}
	b.WriteString("Reviewed-by: X")

	idx := findTrailerBlockStart([]byte(b.String()))
	assert.Equal(t, len(b.String())-len("Reviewed-by: X"), idx)
}
