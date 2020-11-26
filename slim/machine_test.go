package slim

import (
	"fmt"
	"testing"

	"github.com/leodido/go-conventionalcommits"
	"github.com/stretchr/testify/assert"
)

func TestMachineParse(t *testing.T) {
	runner(t, "minimaltypes", testCases, WithTypes(conventionalcommits.TypesMinimal))
}

func TestMachineParseWithFalcoTypes(t *testing.T) {
	runner(t, "falcotypes", testCasesForFalcoTypes, WithTypes(conventionalcommits.TypesFalco))
}

func TestMachineParseWithConventionalTypes(t *testing.T) {
	runner(t, "conventionaltypes", testCasesForConventionalTypes, WithTypes(conventionalcommits.TypesConventional))
}

func runner(t *testing.T, label string, cases []testCase, machineOpts ...conventionalcommits.MachineOption) {
	t.Helper()

	for _, tc := range cases {
		tc := tc
		title := fmt.Sprintf("%s/%s", label, tc.title)

		t.Run(title, func(t *testing.T) {
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
