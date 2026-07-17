// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2020- Leonardo Di Donato <leodidonato@gmail.com>
package parser

import (
	"errors"
	"strings"
	"testing"

	"github.com/leodido/go-conventionalcommits"
	cctesting "github.com/leodido/go-conventionalcommits/testing"
)

// makeFakeTrailerHeavyBody returns a body section composed of `paragraphs`
// paragraphs, each ending in a trailer-shaped line that is NOT actually a
// trailer (it lives mid-body). Used to exercise the issue #38 pre-scan
// path under the worst realistic shape.
func makeFakeTrailerHeavyBody(paragraphs int) string {
	var b strings.Builder
	for i := 0; i < paragraphs; i++ {
		b.WriteString("paragraph prose line one for context\n")
		b.WriteString("paragraph prose line two with detail\n")
		b.WriteString("Looks-Like-Trailer: but-it-is-body\n")
		if i != paragraphs-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

// Avoid compiler optimizations that could remove the actual call we are benchmarking during benchmarks.
var (
	benchParseResult      conventionalcommits.Message
	errBenchParse         error
	benchInvalidUTF8Index int
)

const utf8BenchmarkInputSize = 4 * 1024

type benchCase struct {
	input []byte
	label string
}

type utf8BenchmarkCase struct {
	input          []byte
	label          string
	expectedOffset int
}

// malformedUTF8BenchmarkCases keeps total input size fixed so ns/op isolates
// the cost of reaching the invalid byte instead of also varying buffer length.
func malformedUTF8BenchmarkCases() []utf8BenchmarkCase {
	tests := []struct {
		label  string
		offset int
	}{
		{label: "4 KiB/invalid at start", offset: 0},
		{label: "4 KiB/invalid at midpoint", offset: utf8BenchmarkInputSize / 2},
		{label: "4 KiB/invalid at end", offset: utf8BenchmarkInputSize - 1},
	}

	cases := make([]utf8BenchmarkCase, 0, len(tests))
	for _, tt := range tests {
		input := []byte(strings.Repeat("x", utf8BenchmarkInputSize))
		input[tt.offset] = 0xff
		cases = append(cases, utf8BenchmarkCase{
			label:          tt.label,
			input:          input,
			expectedOffset: tt.offset,
		})
	}

	return cases
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
	// --- issue #38 pre-scan stress: bodies whose paragraphs end in
	// trailer-shaped lines force the pre-scan to traverse many false
	// candidates before locking onto the real trailer block.
	{
		label: "[ok] 10-fake-trailer-paragraphs then real trailer",
		input: []byte("fix: x\n\n" + makeFakeTrailerHeavyBody(10) + "\n\nReviewed-by: X"),
	},
	{
		label: "[ok] 100-fake-trailer-paragraphs then real trailer",
		input: []byte("fix: x\n\n" + makeFakeTrailerHeavyBody(100) + "\n\nReviewed-by: X"),
	},
	{
		label: "[ok] 10-fake-trailer-paragraphs no real trailer",
		input: []byte("fix: x\n\n" + makeFakeTrailerHeavyBody(10)),
	},
	{
		label: "[ok] 100-fake-trailer-paragraphs no real trailer",
		input: []byte("fix: x\n\n" + makeFakeTrailerHeavyBody(100)),
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

func BenchmarkStrictUTF8Validation(b *testing.B) {
	tests := []benchCase{
		{
			label: "minimal ASCII message",
			input: []byte("fix: x"),
		},
		{
			label: "4 KiB UTF-8 body",
			input: []byte("fix: x\n\n" + strings.Repeat("é", 2048)),
		},
	}

	for _, tt := range tests {
		tt := tt
		b.Run(tt.label, func(b *testing.B) {
			b.Run("default", func(b *testing.B) {
				benchmarkUTF8Validation(b, tt.input, false)
			})
			b.Run("strict", func(b *testing.B) {
				benchmarkUTF8Validation(b, tt.input, true)
			})
		})
	}
}

func benchmarkUTF8Validation(b *testing.B, input []byte, strict bool) {
	machine := NewMachine()
	if strict {
		machine.WithStrictUTF8()
	}
	if _, err := machine.Parse(input); err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.SetBytes(int64(len(input)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchParseResult, errBenchParse = machine.Parse(input)
	}
	if errBenchParse != nil {
		b.Fatal(errBenchParse)
	}
}

func BenchmarkStrictUTF8MalformedInput(b *testing.B) {
	for _, tt := range malformedUTF8BenchmarkCases() {
		tt := tt
		b.Run(tt.label, func(b *testing.B) {
			benchmarkStrictUTF8Rejection(b, tt.input, tt.expectedOffset)
		})
	}
}

func benchmarkStrictUTF8Rejection(b *testing.B, input []byte, expectedOffset int) {
	machine := NewMachine()
	machine.WithStrictUTF8()
	message, err := machine.Parse(input)
	if message != nil {
		b.Fatalf("expected a nil message, got %T", message)
	}
	var invalidUTF8Error *InvalidUTF8Error
	if !errors.As(err, &invalidUTF8Error) {
		b.Fatalf("expected InvalidUTF8Error, got %T: %v", err, err)
	}
	if invalidUTF8Error.ByteOffset != expectedOffset {
		b.Fatalf("expected byte offset %d, got %d", expectedOffset, invalidUTF8Error.ByteOffset)
	}

	// Report the inspected prefix, not the full buffer, so early exits do not
	// claim throughput for bytes the validator never visits.
	bytesScanned := expectedOffset + 1
	b.ReportAllocs()
	b.SetBytes(int64(bytesScanned))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchParseResult, errBenchParse = machine.Parse(input)
	}
	b.StopTimer()
	b.ReportMetric(float64(bytesScanned), "B/scanned")
	if benchParseResult != nil {
		b.Fatalf("expected a nil message, got %T", benchParseResult)
	}
	if !errors.Is(errBenchParse, ErrInvalidUTF8) {
		b.Fatalf("expected ErrInvalidUTF8, got %v", errBenchParse)
	}
}

func BenchmarkFirstInvalidUTF8Index(b *testing.B) {
	tests := []utf8BenchmarkCase{
		{
			label:          "4 KiB/valid ASCII",
			input:          []byte(strings.Repeat("x", utf8BenchmarkInputSize)),
			expectedOffset: -1,
		},
		{
			label:          "4 KiB/valid multibyte UTF-8",
			input:          []byte(strings.Repeat("é", utf8BenchmarkInputSize/len("é"))),
			expectedOffset: -1,
		},
	}
	tests = append(tests, malformedUTF8BenchmarkCases()...)

	for _, tt := range tests {
		tt := tt
		b.Run(tt.label, func(b *testing.B) {
			if offset := firstInvalidUTF8Index(tt.input); offset != tt.expectedOffset {
				b.Fatalf("expected byte offset %d, got %d", tt.expectedOffset, offset)
			}

			// Valid input requires a full scan; malformed input stops immediately
			// after the first invalid byte.
			bytesScanned := len(tt.input)
			if tt.expectedOffset >= 0 {
				bytesScanned = tt.expectedOffset + 1
			}
			b.ReportAllocs()
			b.SetBytes(int64(bytesScanned))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				benchInvalidUTF8Index = firstInvalidUTF8Index(tt.input)
			}
			b.StopTimer()
			b.ReportMetric(float64(bytesScanned), "B/scanned")
			if benchInvalidUTF8Index != tt.expectedOffset {
				b.Fatalf("expected byte offset %d, got %d", tt.expectedOffset, benchInvalidUTF8Index)
			}
		})
	}
}
