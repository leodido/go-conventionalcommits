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

func malformedUTF8(prefix string, invalid []byte, suffix string) []byte {
	input := make([]byte, 0, len(prefix)+len(invalid)+len(suffix))
	input = append(input, prefix...)
	input = append(input, invalid...)
	input = append(input, suffix...)

	return input
}

func strictUTF8ErrorAt(column int) string {
	return fmt.Sprintf("invalid UTF-8"+ColumnPositionTemplate, column)
}

func newStrictUTF8Machine(options ...cc.MachineOption) Machine {
	machine := NewMachine(options...)
	machine.WithStrictUTF8()

	return machine
}

func TestStrictUTF8RejectsFirstMalformedByteAnywhereInInput(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
		bad    []byte
		suffix string
	}{
		{
			name:   "free-form type",
			prefix: "",
			bad:    []byte{0xff},
			suffix: ": description",
		},
		{
			name:   "scope",
			prefix: "feat(",
			bad:    []byte{0x80},
			suffix: "): description",
		},
		{
			name:   "description",
			prefix: "feat: description ",
			bad:    []byte{0xc0, 0x80},
			suffix: " suffix",
		},
		{
			name:   "body",
			prefix: "feat: description\n\nbody ",
			bad:    []byte{0xf5},
			suffix: " suffix",
		},
		{
			name:   "trailer value",
			prefix: "feat: description\n\nReviewed-by: ",
			bad:    []byte{0xc3},
			suffix: "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			input := malformedUTF8(tt.prefix, tt.bad, tt.suffix)

			defaultMessage, defaultErr := NewMachine(WithTypes(cc.TypesFreeForm)).Parse(input)
			require.NoError(t, defaultErr, "strict UTF-8 remains opt-in")
			require.NotNil(t, defaultMessage, "default parsing behavior must remain byte-permissive")

			message, err := newStrictUTF8Machine(WithTypes(cc.TypesFreeForm)).Parse(input)

			assert.Nil(t, message)
			require.EqualError(t, err, strictUTF8ErrorAt(len(tt.prefix)))
		})
	}
}

func TestStrictUTF8AcceptsWellFormedUTF8IncludingEncodedRuneError(t *testing.T) {
	tests := []string{
		"féat(scôpe): café",
		"feat: replacement character \uFFFD",
		"feat: description\n\n中文 👍\n\nReviewed-by: Léa",
	}

	for _, input := range tests {
		message, err := newStrictUTF8Machine(WithTypes(cc.TypesFreeForm)).Parse([]byte(input))

		require.NoError(t, err)
		require.NotNil(t, message)
		assert.True(t, message.Ok())
	}
}

func TestStrictUTF8ReportsFirstMalformedByte(t *testing.T) {
	prefix := "feat: description\n\nbody "
	input := malformedUTF8(prefix, []byte{0xff}, " later \xfe")

	message, err := newStrictUTF8Machine().Parse(input)

	assert.Nil(t, message)
	require.EqualError(t, err, strictUTF8ErrorAt(len(prefix)))
}

func TestStrictUTF8PrecedesGrammarValidation(t *testing.T) {
	prefix := "this is not a conventional commit "
	input := malformedUTF8(prefix, []byte{0xff}, "")

	message, err := newStrictUTF8Machine().Parse(input)

	assert.Nil(t, message)
	require.EqualError(t, err, strictUTF8ErrorAt(len(prefix)))
}

func TestStrictUTF8PrecedesBestEffort(t *testing.T) {
	prefix := "feat: description\n\nbody "
	input := malformedUTF8(prefix, []byte{0xff}, "")

	message, err := newStrictUTF8Machine(WithBestEffort()).Parse(input)

	assert.Nil(t, message, "invalid UTF-8 is an input-precondition failure, not a partial parse")
	require.EqualError(t, err, strictUTF8ErrorAt(len(prefix)))
}

func TestStrictUTF8PersistsAcrossParseCalls(t *testing.T) {
	machine := newStrictUTF8Machine()

	message, err := machine.Parse([]byte("feat: valid"))
	require.NoError(t, err)
	require.NotNil(t, message)

	prefix := "feat: invalid "
	message, err = machine.Parse(malformedUTF8(prefix, []byte{0xff}, ""))
	assert.Nil(t, message)
	require.EqualError(t, err, strictUTF8ErrorAt(len(prefix)))
}

func TestStrictUTF8UsesOriginalInputColumnsForIssue56Scenarios(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
		suffix string
	}{
		{
			name:   "body assembly synthesizes newlines",
			prefix: "feat: description\n\nfirst body line\nsecond body line ",
			suffix: "\n\nReviewed-by: Leo",
		},
		{
			name:   "export trims trailing body blank line",
			prefix: "feat: description\n\nbody before trailing blank line ",
			suffix: "\n\n",
		},
		{
			name:   "body-before-blank-line action nudges parser markers",
			prefix: "feat: description\n\nbody paragraph\nFake: trailer-shaped body ",
			suffix: "\n\nReviewed-by: Leo",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			input := malformedUTF8(tt.prefix, []byte{0xff}, tt.suffix)

			message, err := newStrictUTF8Machine().Parse(input)

			assert.Nil(t, message)
			require.EqualError(t, err, strictUTF8ErrorAt(len(tt.prefix)))
		})
	}
}
