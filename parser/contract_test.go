// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2020- Leonardo Di Donato <leodidonato@gmail.com>
package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leodido/go-conventionalcommits"
)

// TestParseNeverReturnsNilNil pins the Parse contract:
//
//	(message, error) := Parse(input)
//	=> message != nil OR error != nil
//
// A `(nil, nil)` return makes callers nil-deref on the second
// branch of the canonical pattern:
//
//	m, err := p.Parse(in)
//	if err != nil { return err }
//	if !m.Ok() { ... }   // <- NPE here
//
// The trailing-CRLF input below historically reached the strict-
// mode return site with m.cs < first_final and m.err == nil. See
// the second-pass review of PR #47.
func TestParseNeverReturnsNilNil(t *testing.T) {
	cases := []struct {
		name  string
		input []byte
	}{
		{
			"trailing-CRLF-after-only-trailer",
			[]byte("feat: x\n\nReviewed-by: Me\r\n"),
		},
		{
			"trailing-CR-after-only-trailer",
			[]byte("feat: x\n\nReviewed-by: Me\r"),
		},
		{
			"trailing-CRLF-after-body-and-trailer",
			[]byte("feat: x\n\nFixes #15\n\nbody text\n\nReviewed-by: Me\r\n"),
		},
	}
	configs := []struct {
		name string
		opt  conventionalcommits.MachineOption
	}{
		{"minimal", WithTypes(conventionalcommits.TypesMinimal)},
		{"conventional", WithTypes(conventionalcommits.TypesConventional)},
		{"falco", WithTypes(conventionalcommits.TypesFalco)},
		{"freeform", WithTypes(conventionalcommits.TypesFreeForm)},
	}
	for _, cfg := range configs {
		for _, tc := range cases {
			t.Run(cfg.name+"/"+tc.name, func(t *testing.T) {
				m, err := NewMachine(cfg.opt).Parse(tc.input)
				if m == nil && err == nil {
					t.Fatalf("Parse contract violated: returned (nil, nil) for input %q", tc.input)
				}
			})
			t.Run(cfg.name+"/"+tc.name+"/best-effort", func(t *testing.T) {
				m, err := NewMachine(cfg.opt, WithBestEffort()).Parse(tc.input)
				if m == nil && err == nil {
					t.Fatalf("Parse contract violated (best-effort): returned (nil, nil) for input %q", tc.input)
				}
			})
		}
	}
	// Sanity: a happy-path input must still return (non-nil, nil).
	m, err := NewMachine().Parse([]byte("fix: x"))
	assert.NotNil(t, m)
	assert.NoError(t, err)
}

// TestParseAcceptsUTF8 sweeps the UTF-8/high-byte acceptance contract
// added by issue #50 across:
//
//	{TypeConfig: minimal, conventional, falco, freeform}
//	x {input shape: trailer value, scope, free-form type}
//	x {UTF-8 family: Latin-1, CJK, emoji, mid-multibyte truncation}
//
// The matrix exists in addition to the per-shape reproducers in
// utf8_test.go to guarantee that NO type config silently regresses
// grammar acceptance: e.g. if a future change hard-codes `print` in
// one production but not another, this sweep notices.
//
// Falco-types is included because it shares the trailer / scope
// productions with the other configs, so high-byte acceptance in
// those productions must hold there too. Free-form-type-utf8 is only
// run under TypesFreeForm because the other configs use literal-string
// type alternations that intentionally reject anything but their
// allow-list.
func TestParseAcceptsUTF8(t *testing.T) {
	type shape struct {
		name             string
		input            []byte
		freeFormTypeOnly bool
	}
	shapes := []shape{
		{"trailer-value-latin", []byte("feat: x\n\nReviewed-by: léodido"), false},
		{"trailer-value-cjk", []byte("feat: x\n\nReviewed-by: 中文"), false},
		{"trailer-value-emoji", []byte("feat: x\n\nReviewed-by: 👍"), false},
		{"trailer-value-mid-multibyte-eof", []byte("feat: x\n\nReviewed-by: \xc3"), false},
		{"scope-latin", []byte("feat(scôpe): x"), false},
		{"scope-cjk", []byte("feat(中文): x"), false},
		{"free-form-type-latin", []byte("féat: x"), true},
		{"free-form-type-cjk", []byte("中文: x"), true},
	}
	configs := []struct {
		name string
		opt  conventionalcommits.MachineOption
		t    conventionalcommits.TypeConfig
	}{
		{"minimal", WithTypes(conventionalcommits.TypesMinimal), conventionalcommits.TypesMinimal},
		{"conventional", WithTypes(conventionalcommits.TypesConventional), conventionalcommits.TypesConventional},
		{"falco", WithTypes(conventionalcommits.TypesFalco), conventionalcommits.TypesFalco},
		{"freeform", WithTypes(conventionalcommits.TypesFreeForm), conventionalcommits.TypesFreeForm},
	}
	for _, cfg := range configs {
		for _, s := range shapes {
			if s.freeFormTypeOnly && cfg.t != conventionalcommits.TypesFreeForm {
				continue
			}
			cfg, s := cfg, s
			t.Run(cfg.name+"/"+s.name, func(t *testing.T) {
				m, err := NewMachine(cfg.opt).Parse(s.input)
				assert.NoError(t, err, "UTF-8 acceptance: input %q must parse", s.input)
				if assert.NotNil(t, m, "UTF-8 acceptance: input %q must yield a message", s.input) {
					assert.True(t, m.Ok(), "UTF-8 acceptance: input %q must produce an Ok() commit", s.input)
				}
			})
		}
	}
}
