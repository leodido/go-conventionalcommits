// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2020- Leonardo Di Donato <leodidonato@gmail.com>
package parser

import (
	"testing"

	"github.com/leodido/go-conventionalcommits"
	cctesting "github.com/leodido/go-conventionalcommits/testing"
)

// Avoid compiler optimizations that could remove the actual call we are benchmarking during benchmarks.
var benchParseResult conventionalcommits.Message

type benchCase struct {
	input []byte
	label string
}

var benchCases = []benchCase{
	{
		label: "[ok] minimal",
		input: []byte("fix: x"),
	},
	{
		label: "[ok] minimal with scope",
		input: []byte("fix(s): x"),
	},
	{
		label: "[ok] minimal breaking with scope",
		input: []byte("fix(s)!: x"),
	},
	{
		label: "[ok] full with 50 characters long description",
		input: []byte("fix(s)!: abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwx"),
	},
	{
		label: "[ok] with body and footer",
		input: []byte(`fix: correct minor typos in code

see the issue for details

on typos fixed.

Reviewed-by: Z
Refs #133`),
	},
	{
		label: "[ok] with body",
		input: []byte(`fix: correct minor typos in code

see the issue for details on typos fixed`),
	},
	{
		label: "[ok] with footer containing one trailer",
		input: []byte(`fix: correct minor typos in code

Acked-by: leodido`),
	},
	{
		label: "[ok] with footer containing many trailers",
		input: []byte(`fix: correct minor typos in code

Acked-by: leodido
Co-authored-by: X
Co-authored-by: Y
Signed-off-by: Leonardo Di Donato <some@email.com>`),
	},
	{
		label: "[no] empty",
		input: []byte(""),
	},
	{
		label: "[no] type but missing colon",
		input: []byte("fix"),
	},
	{
		label: "[no] type but missing description",
		input: []byte("feat: "),
	},
	{
		label: "[no] type and scope but missing description",
		input: []byte("feat(scope): "),
	},
	{
		label: "[no] breaking with type and scope but missing description",
		input: []byte("feat(scope): "),
	},
	{
		// ~~ means it returns the conventionalcommits.Message instance (description cut before the newline) and the error
		label: "[~~] newline in description",
		input: []byte("feat(scope): new\x0Aline"),
	},
	{
		label: "[no] missing whitespace in description",
		input: []byte("feat(scope):a"),
	},
}

func BenchmarkSlimParseMinimalTypes(b *testing.B) {
	for _, tc := range benchCases {
		tc := tc
		m := NewMachine(WithBestEffort())
		b.Run(cctesting.RightPad(tc.label, 50), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				benchParseResult, _ = m.Parse(tc.input)
			}
		})
	}
}

func BenchmarkSlimParseConventionalTypes(b *testing.B) {
	for _, tc := range benchCases {
		tc := tc
		m := NewMachine(WithBestEffort(), WithTypes(conventionalcommits.TypesConventional))
		b.Run(cctesting.RightPad(tc.label, 50), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				benchParseResult, _ = m.Parse(tc.input)
			}
		})
	}
}
