package slim

import (
	"fmt"

	"github.com/leodido/go-conventionalcommits"
)

// ColumnPositionTemplate is the template used to communicate the column where errors occur.
var ColumnPositionTemplate = ": col=%02d"

const (
	// ErrType represents an error in the type part of the commit message.
	ErrType = "illegal '%s' character in commit message type"
	// ErrTypeIncomplete represents an error in the type part of the commit message.
	ErrTypeIncomplete = "incomplete commit message type after '%s' character"
	// ErrEmpty represents an error when the input is empty.
	ErrEmpty = "empty input"
)

const start int = 1
const firstFinal int = 14

const enFail int = 21
const enMain int = 1

type machine struct {
	data       []byte
	cs         int
	p, pe, eof int
	pb         int
	err        error
	bestEffort bool
}

func (m *machine) text() []byte {
	return m.data[m.pb:m.p]
}

func (m *machine) emitErrorOnCurrentCharacter(messageTemplate string) error {
	return fmt.Errorf(messageTemplate+ColumnPositionTemplate, string(m.data[m.p]), m.p)
}

func (m *machine) emitErrorOnPreviousCharacter(messageTemplate string) error {
	return fmt.Errorf(messageTemplate+ColumnPositionTemplate, string(m.data[m.p-1]), m.p)
}

// NewMachine creates a new FSM able to parse Conventional Commits.
func NewMachine(options ...conventionalcommits.MachineOption) conventionalcommits.Machine {
	m := &machine{}

	for _, opt := range options {
		opt(m)
	}

	return m
}

// Err returns the last error occurred.
//
// If the result is nil, then the parsing was successfull.
func (m *machine) Err() error {
	return m.err
}

// Parse parses the input byte array as a Conventional Commit message with no body neither footer.
//
// When a valid Conventional Commit message is given it outputs its structured representation.
// If the parsing detects an error it returns it with the position where the error occurred.
//
// It can also partially parse input messages returning a partially valid structured representation
// and the error that stopped the parsing.
func (m *machine) Parse(input []byte) (conventionalcommits.Message, error) {
	m.data = input
	m.p = 0
	m.pb = 0
	m.pe = len(input)
	m.eof = len(input)
	m.err = nil
	output := &conventionalCommit{}

	{
		m.cs = start
	}

	{
		if (m.p) == (m.pe) {
			goto _testEof
		}
		switch m.cs {
		case 1:
			goto stCase1
		case 0:
			goto stCase0
		case 2:
			goto stCase2
		case 3:
			goto stCase3
		case 4:
			goto stCase4
		case 5:
			goto stCase5
		case 6:
			goto stCase6
		case 7:
			goto stCase7
		case 14:
			goto stCase14
		case 15:
			goto stCase15
		case 8:
			goto stCase8
		case 9:
			goto stCase9
		case 10:
			goto stCase10
		case 11:
			goto stCase11
		case 12:
			goto stCase12
		case 16:
			goto stCase16
		case 17:
			goto stCase17
		case 18:
			goto stCase18
		case 19:
			goto stCase19
		case 20:
			goto stCase20
		case 13:
			goto stCase13
		case 21:
			goto stCase21
		}
		goto stOut
	stCase1:
		if (m.data)[(m.p)] == 102 {
			goto tr0
		}
		goto st0
	tr2:

		if m.p != m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrType)
		} else {
			m.err = m.emitErrorOnPreviousCharacter(ErrTypeIncomplete)
		}
		(m.p)--

		goto st0
	stCase0:
	st0:
		m.cs = 0
		goto _out
	tr0:

		m.pb = m.p

		goto st2
	st2:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof2
		}
	stCase2:
		switch (m.data)[(m.p)] {
		case 101:
			goto st3
		case 105:
			goto st13
		}
		goto tr2
	st3:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof3
		}
	stCase3:
		if (m.data)[(m.p)] == 97 {
			goto st4
		}
		goto tr2
	st4:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof4
		}
	stCase4:
		if (m.data)[(m.p)] == 116 {
			goto st5
		}
		goto tr2
	st5:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof5
		}
	stCase5:
		switch (m.data)[(m.p)] {
		case 33:
			goto tr7
		case 40:
			goto tr8
		case 58:
			goto tr9
		}
		goto tr2
	tr7:

		output._type = string(m.text())

		m.pb = m.p

		goto st6
	st6:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof6
		}
	stCase6:
		if (m.data)[(m.p)] == 58 {
			goto tr10
		}
		goto st0
	tr9:

		output._type = string(m.text())

		goto st7
	tr10:

		output.exclamation = true

		goto st7
	st7:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof7
		}
	stCase7:
		if (m.data)[(m.p)] == 32 {
			goto st14
		}
		goto st0
	st14:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof14
		}
	stCase14:
		goto tr20
	tr20:

		m.pb = m.p

		goto st15
	st15:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof15
		}
	stCase15:
		goto st15
	tr8:

		output._type = string(m.text())

		goto st8
	st8:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof8
		}
	stCase8:
		if (m.data)[(m.p)] == 41 {
			goto tr13
		}
		goto tr12
	tr12:

		m.pb = m.p

		goto st9
	st9:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof9
		}
	stCase9:
		if (m.data)[(m.p)] == 41 {
			goto tr15
		}
		goto st9
	tr13:

		m.pb = m.p

		output.scope = string(m.text())

		goto st10
	tr15:

		output.scope = string(m.text())

		goto st10
	st10:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof10
		}
	stCase10:
		switch (m.data)[(m.p)] {
		case 33:
			goto tr16
		case 41:
			goto tr15
		case 58:
			goto st12
		}
		goto st9
	tr16:

		m.pb = m.p

		goto st11
	st11:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof11
		}
	stCase11:
		switch (m.data)[(m.p)] {
		case 41:
			goto tr15
		case 58:
			goto tr18
		}
		goto st9
	tr18:

		output.exclamation = true

		goto st12
	st12:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof12
		}
	stCase12:
		switch (m.data)[(m.p)] {
		case 32:
			goto st16
		case 41:
			goto tr15
		}
		goto st9
	st16:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof16
		}
	stCase16:
		if (m.data)[(m.p)] == 41 {
			goto tr23
		}
		goto tr22
	tr22:

		m.pb = m.p

		goto st17
	st17:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof17
		}
	stCase17:
		if (m.data)[(m.p)] == 41 {
			goto tr25
		}
		goto st17
	tr25:

		output.scope = string(m.text())

		goto st18
	tr23:

		output.scope = string(m.text())

		m.pb = m.p

		goto st18
	st18:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof18
		}
	stCase18:
		switch (m.data)[(m.p)] {
		case 33:
			goto tr26
		case 41:
			goto tr25
		case 58:
			goto st20
		}
		goto st17
	tr26:

		m.pb = m.p

		goto st19
	st19:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof19
		}
	stCase19:
		switch (m.data)[(m.p)] {
		case 41:
			goto tr25
		case 58:
			goto tr28
		}
		goto st17
	tr28:

		output.exclamation = true

		goto st20
	st20:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof20
		}
	stCase20:
		switch (m.data)[(m.p)] {
		case 32:
			goto st16
		case 41:
			goto tr25
		}
		goto st17
	st13:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof13
		}
	stCase13:
		if (m.data)[(m.p)] == 120 {
			goto st5
		}
		goto tr2
	st21:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof21
		}
	stCase21:
		switch (m.data)[(m.p)] {
		case 10:
			goto st0
		case 13:
			goto st0
		}
		goto st21
	stOut:
	_testEof2:
		m.cs = 2
		goto _testEof
	_testEof3:
		m.cs = 3
		goto _testEof
	_testEof4:
		m.cs = 4
		goto _testEof
	_testEof5:
		m.cs = 5
		goto _testEof
	_testEof6:
		m.cs = 6
		goto _testEof
	_testEof7:
		m.cs = 7
		goto _testEof
	_testEof14:
		m.cs = 14
		goto _testEof
	_testEof15:
		m.cs = 15
		goto _testEof
	_testEof8:
		m.cs = 8
		goto _testEof
	_testEof9:
		m.cs = 9
		goto _testEof
	_testEof10:
		m.cs = 10
		goto _testEof
	_testEof11:
		m.cs = 11
		goto _testEof
	_testEof12:
		m.cs = 12
		goto _testEof
	_testEof16:
		m.cs = 16
		goto _testEof
	_testEof17:
		m.cs = 17
		goto _testEof
	_testEof18:
		m.cs = 18
		goto _testEof
	_testEof19:
		m.cs = 19
		goto _testEof
	_testEof20:
		m.cs = 20
		goto _testEof
	_testEof13:
		m.cs = 13
		goto _testEof
	_testEof21:
		m.cs = 21
		goto _testEof

	_testEof:
		{
		}
		if (m.p) == (m.eof) {
			switch m.cs {
			case 1:

				m.err = fmt.Errorf(ErrEmpty+ColumnPositionTemplate, m.p)
				(m.p)--

			case 2, 3, 4, 5, 13:

				if m.p != m.pe {
					m.err = m.emitErrorOnCurrentCharacter(ErrType)
				} else {
					m.err = m.emitErrorOnPreviousCharacter(ErrTypeIncomplete)
				}
				(m.p)--

			case 15, 17, 18, 19, 20:

				output.descr = string(m.text())

			case 14, 16:

				m.pb = m.p

				output.descr = string(m.text())

			}
		}

	_out:
		{
		}
	}

	if m.cs < firstFinal || m.cs == enFail {
		if m.bestEffort && output.minimal() {
			// An error occurred but partial parsing is on and partial message is minimally valid
			return output.export(), m.err
		}
		return nil, m.err
	}

	return output.export(), nil
}

// WithBestEffort enables best effort mode.
func (m *machine) WithBestEffort() {
	m.bestEffort = true
}

// HasBestEffort tells whether the receiving machine has best effort mode on or off.
func (m *machine) HasBestEffort() bool {
	return m.bestEffort
}
