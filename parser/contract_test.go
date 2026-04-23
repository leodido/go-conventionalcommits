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
