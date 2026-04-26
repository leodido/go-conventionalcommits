// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2020- Leonardo Di Donato <leodidonato@gmail.com>
package parser

import (
	"fmt"
	"testing"

	cc "github.com/leodido/go-conventionalcommits"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// strictCase pins the contract of the WithStrictUTF8 option per
// issue #50 (follow-up to PR #54): when the captured slice for a
// trailer value, scope, or free-form type contains a byte sequence
// that is not well-formed UTF-8, strict mode rejects the parse with
// an error that:
//
//   - includes the failing production name ("trailer value", "scope",
//     "free-form type"),
//   - includes the failing footer key (when applicable),
//   - reports the column of the FIRST invalid byte in the original
//     input (NOT the EOF position; see review §3.1.1),
//   - is deterministic across runs (sorted footer iteration; see
//     review §3.1.2).
//
// expectErr is the EXACT error string asserted via assert.EqualError.
type strictCase struct {
	title     string
	input     []byte
	expectErr string
}

// strictCases is the canonical reproducer set for malformed UTF-8 in
// trailer values, scope, and free-form types under WithStrictUTF8.
var strictCases = []strictCase{
	// Lone leader byte at end of trailer value.
	{
		title:     "trailer-value-lone-leader",
		input:     []byte("feat: x\n\nReviewed-by: \xc3"),
		expectErr: fmt.Sprintf(ErrInvalidUTF8Trailer+ColumnPositionTemplate, "reviewed-by", 22),
	},
	// Lone continuation byte in trailer value.
	{
		title:     "trailer-value-lone-continuation",
		input:     []byte("feat: x\n\nReviewed-by: a\x80b"),
		expectErr: fmt.Sprintf(ErrInvalidUTF8Trailer+ColumnPositionTemplate, "reviewed-by", 23),
	},
	// Overlong-encoding leader (0xc0) in trailer value, padded so the
	// failing column is in the middle, not at EOF.
	{
		title:     "trailer-value-overlong-leader",
		input:     []byte("feat: x\n\nReviewed-by: \xc0\x80 padding text"),
		expectErr: fmt.Sprintf(ErrInvalidUTF8Trailer+ColumnPositionTemplate, "reviewed-by", 22),
	},
	// Malformed UTF-8 in scope.
	{
		title:     "scope-lone-leader",
		input:     []byte("feat(\xc3): x"),
		expectErr: fmt.Sprintf(ErrInvalidUTF8Scope+ColumnPositionTemplate, 5),
	},
	// Malformed UTF-8 in free-form type.
	{
		title:     "free-form-type-lone-leader",
		input:     []byte("\xc3: x"),
		expectErr: fmt.Sprintf(ErrInvalidUTF8Type+ColumnPositionTemplate, 0),
	},
}

// TestStrictUTF8RejectsInvalidUTF8 pins the EXACT error string for
// each malformed input under WithStrictUTF8, including the column
// (review §3.1.1, §5.2) and the footer key for trailer-value cases
// (review §3.4).
func TestStrictUTF8RejectsInvalidUTF8(t *testing.T) {
	for _, tc := range strictCases {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			// Default mode (no strict): bytes accepted.
			m, err := NewMachine(WithTypes(cc.TypesFreeForm)).Parse(tc.input)
			assert.NoError(t, err, "default mode must accept input %q", tc.input)
			assert.NotNil(t, m, "default mode must yield a message for %q", tc.input)

			// Strict mode: rejects with EXACT error string.
			m2, err2 := NewMachine(WithTypes(cc.TypesFreeForm), WithStrictUTF8()).Parse(tc.input)
			assert.Nil(t, m2, "strict mode must return nil message for %q", tc.input)
			require.Error(t, err2, "strict mode must return an error for %q", tc.input)
			assert.EqualError(t, err2, tc.expectErr)
		})
	}
}

// TestStrictUTF8AcceptsValidUTF8 pins that valid UTF-8 (Latin-1, CJK,
// emoji) parses cleanly under WithStrictUTF8.
func TestStrictUTF8AcceptsValidUTF8(t *testing.T) {
	cases := []struct {
		title string
		input []byte
	}{
		{"trailer-latin", []byte("feat: x\n\nReviewed-by: léodido")},
		{"trailer-cjk", []byte("feat: x\n\nReviewed-by: 中文")},
		{"trailer-emoji", []byte("feat: x\n\nReviewed-by: 👍")},
		{"trailer-multiline-cjk", []byte("feat: x\n\nBREAKING CHANGE: 中文\n  说明")},
		{"scope-latin", []byte("feat(scôpe): x")},
		{"free-form-type-latin", []byte("féat: x")},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			m, err := NewMachine(WithTypes(cc.TypesFreeForm), WithStrictUTF8()).Parse(tc.input)
			require.NoError(t, err, "strict mode must accept valid UTF-8 in %q", tc.input)
			require.NotNil(t, m)
			assert.True(t, m.Ok())
		})
	}
}

// TestStrictUTF8DeterministicFooterError pins that when MULTIPLE
// footers contain malformed UTF-8, the error is deterministic across
// runs (sorted footer iteration; review §3.1.2). With sorted keys the
// FIRST footer alphabetically (here "a-aaa") is the one named in the
// error, regardless of the FSM's insertion order.
func TestStrictUTF8DeterministicFooterError(t *testing.T) {
	in := []byte("feat: x\n\nA-aaa: \xc3\nB-bbb: \xee")
	for i := 0; i < 16; i++ {
		_, err := NewMachine(WithStrictUTF8()).Parse(in)
		require.Error(t, err)
		assert.EqualError(
			t, err,
			fmt.Sprintf(ErrInvalidUTF8Trailer+ColumnPositionTemplate, "a-aaa", 16),
			"iteration %d: deterministic footer error required", i,
		)
	}
}

// TestStrictUTF8BodyOptionRejectsInvalidUTF8 pins the contract of the
// SEPARATE WithStrictUTF8Body option: it is independent of
// WithStrictUTF8, only validates body and description, and reports
// the column of the first invalid byte in the original input.
func TestStrictUTF8BodyOptionRejectsInvalidUTF8(t *testing.T) {
	cases := []struct {
		title     string
		input     []byte
		expectErr string
	}{
		{
			"description-lone-leader",
			[]byte("feat: hello \xc3"),
			fmt.Sprintf(ErrInvalidUTF8Description+ColumnPositionTemplate, 12),
		},
		{
			"body-lone-leader",
			[]byte("feat: x\n\nbody body \xc3 more body"),
			fmt.Sprintf(ErrInvalidUTF8Body+ColumnPositionTemplate, 19),
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			// Default mode (no body strict): bytes accepted.
			m, err := NewMachine().Parse(tc.input)
			assert.NoError(t, err, "default mode must accept input %q", tc.input)
			assert.NotNil(t, m)

			// WithStrictUTF8Body alone: rejects body / description.
			m2, err2 := NewMachine(WithStrictUTF8Body()).Parse(tc.input)
			assert.Nil(t, m2)
			require.Error(t, err2)
			assert.EqualError(t, err2, tc.expectErr)

			// WithStrictUTF8 alone (no body knob): MUST accept body
			// content with high bytes.
			m3, err3 := NewMachine(WithStrictUTF8()).Parse(tc.input)
			assert.NoError(t, err3, "WithStrictUTF8 alone must NOT validate body / description")
			assert.NotNil(t, m3)
		})
	}
}

// TestStrictUTF8BestEffortReturnsCorruptedPartial pins what the
// caller actually gets when (BestEffort + StrictUTF8) is set AND the
// trailer block is malformed UTF-8. Per review §3.3 the field
// contents must be pinned, including U+FFFD-corrupted Type/Scope.
func TestStrictUTF8BestEffortReturnsCorruptedPartial(t *testing.T) {
	// The trailer value contains a lone leader byte. Best effort + strict:
	// validation fails, but we still get the partial commit.
	in := []byte("feat: x\n\nReviewed-by: \xc3")
	m, err := NewMachine(WithBestEffort(), WithStrictUTF8()).Parse(in)
	require.Error(t, err)
	assert.EqualError(t, err, fmt.Sprintf(ErrInvalidUTF8Trailer+ColumnPositionTemplate, "reviewed-by", 22))
	require.NotNil(t, m, "best-effort must yield the partial commit even when strict validation rejects it")
	c := m.(*cc.ConventionalCommit)
	// Type and Description survive clean (they're ASCII).
	assert.Equal(t, "feat", c.Type)
	assert.Equal(t, "x", c.Description)
	// The footer was captured before validation kicked in. The bytes are
	// preserved verbatim in the captured value.
	require.Contains(t, c.Footers, "reviewed-by")
	assert.Equal(t, "\xc3", c.Footers["reviewed-by"][0])
	// Ok() is satisfied because Type + Description are non-empty.
	assert.True(t, m.Ok())
}

// TestStrictUTF8OptionAccessors pins that WithStrictUTF8 and
// WithStrictUTF8Body each set their own flag, and that
// HasStrictUTF8 / HasStrictUTF8Body report state correctly without
// either of them implying the other.
func TestStrictUTF8OptionAccessors(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		m := NewMachine()
		assert.False(t, m.(*machine).HasStrictUTF8())
		assert.False(t, m.(*machine).HasStrictUTF8Body())
	})
	t.Run("strict-only", func(t *testing.T) {
		m := NewMachine(WithStrictUTF8())
		assert.True(t, m.(*machine).HasStrictUTF8())
		assert.False(t, m.(*machine).HasStrictUTF8Body())
	})
	t.Run("body-only", func(t *testing.T) {
		m := NewMachine(WithStrictUTF8Body())
		assert.False(t, m.(*machine).HasStrictUTF8())
		assert.True(t, m.(*machine).HasStrictUTF8Body())
	})
	t.Run("both", func(t *testing.T) {
		m := NewMachine(WithStrictUTF8(), WithStrictUTF8Body())
		assert.True(t, m.(*machine).HasStrictUTF8())
		assert.True(t, m.(*machine).HasStrictUTF8Body())
	})
}
