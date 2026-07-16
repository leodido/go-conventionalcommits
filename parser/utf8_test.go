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

// utf8Case pins the per-shape acceptance contract for an issue-#50
// input across one type config, in default mode (no options).
//
// When expectOk is true, expectMsg is invoked to check the parsed
// payload. When expectOk is false, expectErr is asserted as the exact
// error string.
type utf8Case struct {
	title     string
	input     []byte
	types     cc.TypeConfig
	expectOk  bool
	expectErr string                           // exact match when expectOk == false
	expectMsg func(t *testing.T, m cc.Message) // payload check when expectOk == true
}

// utf8Cases is the canonical issue-#50 reproducer set. Each entry
// pins:
//   - which input shape (trailer value, scope, free-form type),
//   - which type config it runs under,
//   - the expected payload (GREEN), or the rejection message for the
//     two negative cases that MUST continue to be rejected after the
//     widening (DEL and SOH).
var utf8Cases = []utf8Case{
	// AC1: Latin-1 UTF-8 in trailer value (TypesMinimal).
	{
		title:     "trailer-value-utf8-latin/minimal",
		input:     []byte("feat: x\n\nReviewed-by: léodido"),
		types:     cc.TypesMinimal,
		expectOk:  true,
		expectMsg: expectFooter("reviewed-by", "léodido"),
	},
	// AC1: Latin-1 UTF-8 in trailer value (TypesConventional).
	{
		title:     "trailer-value-utf8-latin/conventional",
		input:     []byte("feat: x\n\nReviewed-by: léodido"),
		types:     cc.TypesConventional,
		expectOk:  true,
		expectMsg: expectFooter("reviewed-by", "léodido"),
	},
	// AC1: Latin-1 UTF-8 in trailer value (TypesFreeForm).
	{
		title:     "trailer-value-utf8-latin/freeform",
		input:     []byte("feat: x\n\nReviewed-by: léodido"),
		types:     cc.TypesFreeForm,
		expectOk:  true,
		expectMsg: expectFooter("reviewed-by", "léodido"),
	},
	// AC2: 4-byte UTF-8 (emoji) in trailer value.
	{
		title:     "trailer-value-utf8-emoji/minimal",
		input:     []byte("feat: x\n\nReviewed-by: 👍"),
		types:     cc.TypesMinimal,
		expectOk:  true,
		expectMsg: expectFooter("reviewed-by", "👍"),
	},
	// AC3: UTF-8 in scope (TypesMinimal).
	{
		title:     "scope-utf8/minimal",
		input:     []byte("feat(scôpe): x"),
		types:     cc.TypesMinimal,
		expectOk:  true,
		expectMsg: expectScope("scôpe"),
	},
	// AC3: UTF-8 in scope (TypesConventional).
	{
		title:     "scope-utf8/conventional",
		input:     []byte("feat(scôpe): x"),
		types:     cc.TypesConventional,
		expectOk:  true,
		expectMsg: expectScope("scôpe"),
	},
	// AC3: UTF-8 in scope (TypesFreeForm).
	{
		title:     "scope-utf8/freeform",
		input:     []byte("feat(scôpe): x"),
		types:     cc.TypesFreeForm,
		expectOk:  true,
		expectMsg: expectScope("scôpe"),
	},
	// AC4: UTF-8 inside a free-form type (`féat`).
	{
		title:     "free-form-type-utf8-mid/freeform",
		input:     []byte("féat: x"),
		types:     cc.TypesFreeForm,
		expectOk:  true,
		expectMsg: expectType("féat"),
	},
	// AC4: free-form type made entirely of UTF-8 bytes (`中文`).
	{
		title:     "free-form-type-utf8-all/freeform",
		input:     []byte("中文: x"),
		types:     cc.TypesFreeForm,
		expectOk:  true,
		expectMsg: expectType("中文"),
	},
	// AC5: mixed ASCII + UTF-8 in trailer value.
	{
		title:     "trailer-value-mixed-ascii-utf8/minimal",
		input:     []byte("feat: x\n\nReviewed-by: leo & léo"),
		types:     cc.TypesMinimal,
		expectOk:  true,
		expectMsg: expectFooter("reviewed-by", "leo & léo"),
	},
	// AC6: multi-line trailer value with CJK UTF-8 across the line
	// boundary.
	{
		title:     "trailer-value-multiline-cjk/minimal",
		input:     []byte("feat: x\n\nBREAKING CHANGE: 中文\n  说明"),
		types:     cc.TypesMinimal,
		expectOk:  true,
		expectMsg: expectFooter("breaking-change", "中文\n  说明"),
	},
	// AC7: trailer value ending mid-multibyte at EOF (truncated UTF-8
	// leader byte). The byte is captured opaquely; the parser does
	// not panic. Validating UTF-8 well-formedness is the caller's
	// responsibility (an opt-in option for that lives in a follow-up
	// PR).
	{
		title:     "trailer-value-mid-multibyte-eof/minimal",
		input:     []byte("feat: x\n\nReviewed-by: \xc3"),
		types:     cc.TypesMinimal,
		expectOk:  true,
		expectMsg: expectFooter("reviewed-by", "\xc3"),
	},
	// AC8: DEL (\x7f) in a trailer value MUST continue to be rejected:
	// the widened class is `(print | 0x80..0xff)`, which still
	// excludes \x7f.
	{
		title:     "trailer-value-del-still-rejected/minimal",
		input:     []byte("feat: x\n\nReviewed-by: a\x7fb"),
		types:     cc.TypesMinimal,
		expectOk:  false,
		expectErr: fmt.Sprintf(ErrEarly+ColumnPositionTemplate, "a", 23),
	},
	// AC9: SOH (\x01) in a trailer value MUST continue to be rejected:
	// control bytes are not in `(print | 0x80..0xff)`.
	{
		title:     "trailer-value-control-still-rejected/minimal",
		input:     []byte("feat: x\n\nReviewed-by: a\x01b"),
		types:     cc.TypesMinimal,
		expectOk:  false,
		expectErr: fmt.Sprintf(ErrEarly+ColumnPositionTemplate, "a", 23),
	},
}

// TestUTF8HighByteAcceptance exercises the issue-#50 reproducer set
// in default mode (no options).
//
// Best-effort mode is intentionally NOT covered here because it
// swallows trailer-block errors and returns the description-only
// commit, which is a separate contract from the reject/accept
// decision being asserted.
//
// Strict UTF-8 validation (the Machine.WithStrictUTF8 setting)
// is intentionally NOT covered here either. This test pins which
// bytes the FSM accepts; opt-in behavior lives in its own test file.
func TestUTF8HighByteAcceptance(t *testing.T) {
	for _, tc := range utf8Cases {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			m, err := NewMachine(WithTypes(tc.types)).Parse(tc.input)

			if tc.expectOk {
				assert.NoError(t, err)
				require.NotNil(t, m, "expected non-nil message for valid input %q", tc.input)
				assert.True(t, m.Ok())
				if tc.expectMsg != nil {
					tc.expectMsg(t, m)
				}

				return
			}

			assert.Nil(t, m)
			assert.Error(t, err)
			assert.EqualError(t, err, tc.expectErr)
		})
	}
}

// TestUTF8ExportNormalization distinguishes grammar acceptance from
// the public representation returned by Parse. Type and Scope pass
// through strings.ToLower during export, while trailer values are
// exported without normalization.
func TestUTF8ExportNormalization(t *testing.T) {
	cases := []struct {
		title     string
		input     []byte
		expectMsg func(t *testing.T, m cc.Message)
	}{
		{
			title:     "uppercase-unicode-free-form-type-is-lowercased",
			input:     []byte("FÉAT: x"),
			expectMsg: expectType("féat"),
		},
		{
			title:     "uppercase-unicode-scope-is-lowercased",
			input:     []byte("feat(SCÔPE): x"),
			expectMsg: expectScope("scôpe"),
		},
		{
			title:     "malformed-free-form-type-becomes-replacement-rune",
			input:     []byte{0xff, ':', ' ', 'x'},
			expectMsg: expectType("\uFFFD"),
		},
		{
			title:     "malformed-scope-becomes-replacement-rune",
			input:     []byte{'f', 'e', 'a', 't', '(', 0xff, ')', ':', ' ', 'x'},
			expectMsg: expectScope("\uFFFD"),
		},
		{
			title: "malformed-trailer-value-is-preserved",
			input: []byte{
				'f', 'e', 'a', 't', ':', ' ', 'x', '\n', '\n',
				'R', 'e', 'v', 'i', 'e', 'w', 'e', 'd', '-', 'b', 'y', ':', ' ', 0xff,
			},
			expectMsg: expectFooter("reviewed-by", "\xff"),
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			m, err := NewMachine(WithTypes(cc.TypesFreeForm)).Parse(tc.input)
			require.NoError(t, err)
			require.NotNil(t, m)
			tc.expectMsg(t, m)
		})
	}
}

// TestUTF8MidMultibyteEOFNoPanic asserts that an input ending mid-
// multibyte does not panic the parser regardless of mode or type
// config. Pinned separately because the panic surface is independent
// of the success/failure of the parse.
func TestUTF8MidMultibyteEOFNoPanic(t *testing.T) {
	inputs := [][]byte{
		[]byte("feat: x\n\nReviewed-by: \xc3"),
		[]byte("feat: x\n\nReviewed-by: \xe4\xb8"),
		[]byte("feat: x\n\nReviewed-by: \xf0\x9f\x91"),
	}
	configs := []cc.TypeConfig{
		cc.TypesMinimal,
		cc.TypesConventional,
		cc.TypesFalco,
		cc.TypesFreeForm,
	}
	for _, in := range inputs {
		for _, tc := range configs {
			func() {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("Parse panicked on input %q (config %d): %v", in, tc, r)
					}
				}()
				_, _ = NewMachine(WithTypes(tc)).Parse(in)
				_, _ = NewMachine(WithTypes(tc), WithBestEffort()).Parse(in)
			}()
		}
	}
}

// expectFooter returns a payload check that asserts the parsed
// commit has a footer with the given key whose first value equals
// want.
func expectFooter(key, want string) func(t *testing.T, m cc.Message) {
	return func(t *testing.T, m cc.Message) {
		c := m.(*cc.ConventionalCommit)
		require.Contains(t, c.Footers, key, "missing footer %q", key)
		require.NotEmpty(t, c.Footers[key], "footer %q has no values", key)
		assert.Equal(t, want, c.Footers[key][0])
	}
}

// expectScope returns a payload check that asserts the parsed commit
// has the given scope.
func expectScope(want string) func(t *testing.T, m cc.Message) {
	return func(t *testing.T, m cc.Message) {
		c := m.(*cc.ConventionalCommit)
		require.NotNil(t, c.Scope, "expected scope, got nil")
		assert.Equal(t, want, *c.Scope)
	}
}

// expectType returns a payload check that asserts the parsed commit
// has the given type.
func expectType(want string) func(t *testing.T, m cc.Message) {
	return func(t *testing.T, m cc.Message) {
		c := m.(*cc.ConventionalCommit)
		assert.Equal(t, want, c.Type)
	}
}
