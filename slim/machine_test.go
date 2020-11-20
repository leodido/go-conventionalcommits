package slim

import (
	"fmt"
	"testing"

	"github.com/leodido/go-conventionalcommits"
	cctesting "github.com/leodido/go-conventionalcommits/testing"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	title        string
	input        []byte
	ok           bool
	value        conventionalcommits.Message
	partialValue conventionalcommits.Message
	errorString  string
}

var testCases = []testCase{
	// INVALID / empty
	{
		"empty",
		[]byte(""),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEmpty+ColumnPositionTemplate, 0),
	},
	// INVALID / invalid type (1 char)
	{
		"invalid-type-1-char",
		[]byte("f"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "f", 1),
	},
	// INVALID / invalid type (2 char)
	{
		"invalid-type-2-char",
		[]byte("fx"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "x", 1),
	},
	// INVALID / invalid type (3 char)
	{
		"invalid-type-3-char",
		[]byte("fit"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "t", 2),
	},
	// INVALID / invalid type (4 char)
	{
		"invalid-type-4-char",
		[]byte("feax"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "x", 3),
	},
	// INVALID / missing colon after type fix
	{
		"invalid-after-valid-type-fix",
		[]byte("fix"),
		false,
		nil,
		nil, // no partial result because it is not a minimal valid commit message
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, "x", 2),
	},
	// INVALID / missing colon after type feat
	{
		"invalid-after-valid-type-feat",
		[]byte("feat"),
		false,
		nil,
		nil, // no partial result because it is not a minimal valid commit message
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, "t", 3),
	},
	// INVALID / invalid type (2 char) + colon
	{
		"invalid-type-2-char-colon",
		[]byte("fi:"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, ":", 2),
	},
	// INVALID / invalid type (3 char) + colon
	{
		"invalid-type-3-char-colon",
		[]byte("fea:"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, ":", 3),
	},
	// VALID / minimal commit message
	{
		"valid-minimal-commit-message",
		[]byte("fix: x"),
		true,
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "fix",
				Description: "x",
			},
		},
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "fix",
				Description: "x",
			},
		},
		"",
	},
	// INVALID / missing colon after valid commit message type
	{
		"missing-colon-after-type-3-chars",
		[]byte("fix>"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, ">", 3),
	},
	// INVALID / missing colon after valid commit message type
	{
		"missing-colon-after-type-4-chars",
		[]byte("feat?"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, "?", 4),
	},
	// VALID // type + scope + description
	{
		"valid-with-scope",
		[]byte("fix(aaa): bbb"),
		true,
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "fix",
				Scope:       cctesting.StringAddress("aaa"),
				Description: "bbb",
			},
		},
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "fix",
				Scope:       cctesting.StringAddress("aaa"),
				Description: "bbb",
			},
		},
		"",
	},
	// VALID // type + scope + breaking + description
	{
		"valid-breaking-with-scope",
		[]byte("fix(aaa)!: bbb"),
		true,
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "fix",
				Scope:       cctesting.StringAddress("aaa"),
				Description: "bbb",
				Exclamation: true,
			},
		},
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "fix",
				Scope:       cctesting.StringAddress("aaa"),
				Description: "bbb",
				Exclamation: true,
			},
		},
		"",
	},
}

func TestMachineParse(t *testing.T) {
	fmt.Println("CIAONE")
	runner(t, testCases)
}

func runner(t *testing.T, cases []testCase, machineOpts ...conventionalcommits.MachineOption) {
	t.Helper()

	for _, tc := range cases {
		tc := tc

		t.Run(tc.title, func(t *testing.T) {
			message, messageErr := NewMachine(machineOpts...).Parse(tc.input)
			partial, partialErr := NewMachine(append(machineOpts, WithBestEffort())...).Parse(tc.input)

			if !tc.ok {
				// We expect the test case input to be an invalid commit message
				assert.Nil(t, message)
				assert.Error(t, messageErr)
				assert.EqualError(t, messageErr, tc.errorString)

				// In this case can happen that with best effort mode o
				// the result is not nil rather it contains a minimal valid result
				if partial != nil {
					assert.True(t, partial.Ok())
				}
				assert.Equal(t, tc.partialValue, partial)
				assert.EqualError(t, partialErr, tc.errorString)
			} else {
				// We expect the test case intput to be a valid commit message
				assert.Nil(t, messageErr)
				assert.NotEmpty(t, message)
				assert.True(t, message.Ok())
				assert.Equal(t, message, partial)
				assert.Equal(t, tc.partialValue, partial)
				assert.Equal(t, messageErr, partialErr)
			}

			assert.Equal(t, tc.value, message)
		})
	}
}

func TestMachineBestEffortOption(t *testing.T) {
	p1 := NewMachine().(conventionalcommits.BestEfforter)
	assert.False(t, p1.HasBestEffort())

	p2 := NewMachine(WithBestEffort()).(conventionalcommits.BestEfforter)
	assert.True(t, p2.HasBestEffort())
}

func TestMachineTypeConfigOption(t *testing.T) {
	p := NewMachine(WithTypes(conventionalcommits.TypesFalco))
	mes, err := p.Parse([]byte("new: ciao"))

	res := &ConventionalCommit{
		Minimal: conventionalcommits.Minimal{
			Type:        "new",
			Description: "ciao",
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, res, mes)
}
