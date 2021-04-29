package parser

import (
	"testing"

	"github.com/leodido/go-conventionalcommits"
	cctesting "github.com/leodido/go-conventionalcommits/testing"
)

// Avoid compiler optimizations that could remove the actual call we are benchmarking during benchmarks
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
