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

// utf8Case pins the per-shape parse contract for an issue-#50 input
// across one type config, in default mode (no options).
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
//   - the exact default-mode error today (RED), or the expected
//     payload once the FSM is widened (GREEN).
//
// The RED commit registers these with expectOk == false and the
// captured error string. The fix commit flips expectOk to true and
// replaces expectErr with an expectMsg payload check. Bisecting the
// branch shows exactly one commit per assertion flip.
var utf8Cases = []utf8Case{
	// AC1: Latin-1 UTF-8 in trailer value (TypesMinimal).
	{
		title:     "trailer-value-utf8-latin/minimal",
		input:     []byte("feat: x\n\nReviewed-by: léodido"),
		types:     cc.TypesMinimal,
		expectOk:  false,
		expectErr: fmt.Sprintf(ErrEarly+ColumnPositionTemplate, "l", 23),
	},
	// AC1: Latin-1 UTF-8 in trailer value (TypesConventional).
	{
		title:     "trailer-value-utf8-latin/conventional",
		input:     []byte("feat: x\n\nReviewed-by: léodido"),
		types:     cc.TypesConventional,
		expectOk:  false,
		expectErr: fmt.Sprintf(ErrEarly+ColumnPositionTemplate, "l", 23),
	},
	// AC1: Latin-1 UTF-8 in trailer value (TypesFreeForm).
	{
		title:     "trailer-value-utf8-latin/freeform",
		input:     []byte("feat: x\n\nReviewed-by: léodido"),
		types:     cc.TypesFreeForm,
		expectOk:  false,
		expectErr: fmt.Sprintf(ErrEarly+ColumnPositionTemplate, "l", 23),
	},
	// AC2: 4-byte UTF-8 (emoji) in trailer value.
	{
		title:     "trailer-value-utf8-emoji/minimal",
		input:     []byte("feat: x\n\nReviewed-by: 👍"),
		types:     cc.TypesMinimal,
		expectOk:  false,
		expectErr: fmt.Sprintf(ErrEarly+ColumnPositionTemplate, " ", 22),
	},
	// AC3: UTF-8 in scope (TypesMinimal).
	{
		title:     "scope-utf8/minimal",
		input:     []byte("feat(scôpe): x"),
		types:     cc.TypesMinimal,
		expectOk:  false,
		expectErr: fmt.Sprintf(ErrScope+ColumnPositionTemplate, "Ã", 7),
	},
	// AC3: UTF-8 in scope (TypesConventional).
	{
		title:     "scope-utf8/conventional",
		input:     []byte("feat(scôpe): x"),
		types:     cc.TypesConventional,
		expectOk:  false,
		expectErr: fmt.Sprintf(ErrScope+ColumnPositionTemplate, "Ã", 7),
	},
	// AC3: UTF-8 in scope (TypesFreeForm).
	{
		title:     "scope-utf8/freeform",
		input:     []byte("feat(scôpe): x"),
		types:     cc.TypesFreeForm,
		expectOk:  false,
		expectErr: fmt.Sprintf(ErrScope+ColumnPositionTemplate, "Ã", 7),
	},
	// AC4: UTF-8 inside a free-form type (`féat`).
	{
		title:     "free-form-type-utf8-mid/freeform",
		input:     []byte("féat: x"),
		types:     cc.TypesFreeForm,
		expectOk:  false,
		expectErr: fmt.Sprintf(ErrColon+ColumnPositionTemplate, "Ã", 1),
	},
	// AC4: free-form type made entirely of UTF-8 bytes (`中文`).
	{
		title:     "free-form-type-utf8-all/freeform",
		input:     []byte("中文: x"),
		types:     cc.TypesFreeForm,
		expectOk:  false,
		expectErr: fmt.Sprintf(ErrType+ColumnPositionTemplate, "ä", 0),
	},
	// AC5: mixed ASCII + UTF-8 in trailer value.
	{
		title:     "trailer-value-mixed-ascii-utf8/minimal",
		input:     []byte("feat: x\n\nReviewed-by: leo & léo"),
		types:     cc.TypesMinimal,
		expectOk:  false,
		expectErr: fmt.Sprintf(ErrEarly+ColumnPositionTemplate, "l", 29),
	},
	// AC6: multi-line trailer value with CJK UTF-8 across the line
	// boundary.
	{
		title:     "trailer-value-multiline-cjk/minimal",
		input:     []byte("feat: x\n\nBREAKING CHANGE: 中文\n  说明"),
		types:     cc.TypesMinimal,
		expectOk:  false,
		expectErr: fmt.Sprintf(ErrEarly+ColumnPositionTemplate, " ", 26),
	},
	// AC7: trailer value ending mid-multibyte at EOF (truncated UTF-8
	// leader byte). Today this is rejected; after #50 the byte is
	// captured opaquely and the parser does not panic.
	{
		title:     "trailer-value-mid-multibyte-eof/minimal",
		input:     []byte("feat: x\n\nReviewed-by: \xc3"),
		types:     cc.TypesMinimal,
		expectOk:  false,
		expectErr: fmt.Sprintf(ErrEarly+ColumnPositionTemplate, " ", 22),
	},
	// AC8: DEL (\x7f) in a trailer value MUST continue to be rejected
	// after #50: the widened class is `(print | 0x80..0xff)`, which
	// still excludes \x7f.
	{
		title:     "trailer-value-del-still-rejected/minimal",
		input:     []byte("feat: x\n\nReviewed-by: a\x7fb"),
		types:     cc.TypesMinimal,
		expectOk:  false,
		expectErr: fmt.Sprintf(ErrEarly+ColumnPositionTemplate, "a", 23),
	},
	// AC9: SOH (\x01) in a trailer value MUST continue to be rejected
	// after #50: control bytes are not in `(print | 0x80..0xff)`.
	{
		title:     "trailer-value-control-still-rejected/minimal",
		input:     []byte("feat: x\n\nReviewed-by: a\x01b"),
		types:     cc.TypesMinimal,
		expectOk:  false,
		expectErr: fmt.Sprintf(ErrEarly+ColumnPositionTemplate, "a", 23),
	},
}

// TestUTF8ByteTransparency exercises the issue-#50 reproducer set in
// default mode (no options). RED at the test-only commit, GREEN after
// the FSM is widened.
//
// Best-effort mode is intentionally NOT covered here because it
// swallows trailer-block errors and returns the description-only
// commit, which is a separate contract from the reject/accept
// decision being asserted.
//
// Strict-mode validation (the WithStrictUTF8 option, follow-up PR)
// is intentionally NOT covered here either. This test pins what the
// FSM does on its own; option behavior lives in its own test file.
func TestUTF8ByteTransparency(t *testing.T) {
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
